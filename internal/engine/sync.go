package engine

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/geomark27/sync-agent/internal/cloud"
	"github.com/geomark27/sync-agent/internal/config"
	"github.com/geomark27/sync-agent/internal/watcher"
	"github.com/geomark27/sync-agent/pkg/hasher"
)

// debounceDelay es la ventana de quietud antes de agrupar y procesar cambios.
const debounceDelay = 2 * time.Second

// Engine orquesta la sincronización: une configuración, watcher, hasher y nube.
type Engine struct {
	cfg         *config.AppConfig
	cloud       cloud.Provider
	localByName map[string]string // nombre remoto -> ruta local absoluta
	hashes      map[string]string // ruta local -> último hash conocido
}

// New construye el motor a partir de la configuración y un proveedor de nube.
// El nombre remoto de cada archivo es su nombre base; las colisiones se avisan
// y se ignoran.
func New(cfg *config.AppConfig, provider cloud.Provider) *Engine {
	localByName := make(map[string]string, len(cfg.Paths))
	for _, p := range cfg.Paths {
		abs, err := filepath.Abs(p)
		if err != nil {
			abs = p
		}
		name := filepath.Base(abs)
		if existing, ok := localByName[name]; ok {
			log.Printf("⚠️  colisión de nombre '%s' entre %q y %q; se ignora el segundo", name, existing, abs)
			continue
		}
		localByName[name] = abs
	}

	return &Engine{
		cfg:         cfg,
		cloud:       provider,
		localByName: localByName,
		hashes:      make(map[string]string),
	}
}

// Run ejecuta el ciclo de vida completo: pull inicial, registro de hashes,
// vigilancia y push de cambios. Retorna cuando el contexto se cancela.
func (e *Engine) Run(ctx context.Context) error {
	// 1. Traer el estado remoto y aplicarlo localmente si difiere.
	if err := e.initialPull(ctx); err != nil {
		log.Printf("⚠️  pull inicial falló: %v", err)
	}

	// 2. Registrar los hashes locales actuales (tras el pull).
	for _, path := range e.localByName {
		if _, ok := e.hashes[path]; ok {
			continue // ya calculado durante el pull
		}
		if h, err := hasher.HashFile(path); err == nil {
			e.hashes[path] = h
		}
	}

	// 3. Vigilar los directorios padre de cada archivo configurado.
	w, err := watcher.New()
	if err != nil {
		return err
	}
	dirs := make(map[string]struct{})
	for _, path := range e.localByName {
		dirs[filepath.Dir(path)] = struct{}{}
	}
	for dir := range dirs {
		if err := w.Add(dir); err != nil {
			log.Printf("⚠️  no se pudo vigilar %q: %v", dir, err)
		}
	}
	go w.Run(ctx)

	// 4. Procesar lotes de cambios ya agrupados por el debouncer.
	batches := watcher.Debounce(ctx, w.Events(), debounceDelay)
	log.Printf("👀 Vigilando %d archivo(s)...", len(e.localByName))

	for {
		select {
		case <-ctx.Done():
			return nil
		case batch, ok := <-batches:
			if !ok {
				return nil
			}
			e.handleBatch(ctx, batch)
		}
	}
}

// initialPull descarga el Gist y escribe localmente los archivos cuyo contenido
// remoto difiera del local. En el arranque, la nube tiene prioridad.
func (e *Engine) initialPull(ctx context.Context) error {
	remote, err := e.cloud.Pull(ctx)
	if err != nil {
		return err
	}

	for name, content := range remote {
		localPath, ok := e.localByName[name]
		if !ok {
			continue // archivo remoto no mapeado en esta máquina
		}

		remoteHash := hasher.HashBytes([]byte(content))
		if localHash, err := hasher.HashFile(localPath); err == nil && localHash == remoteHash {
			e.hashes[localPath] = localHash
			continue // ya están sincronizados
		}

		if err := os.MkdirAll(filepath.Dir(localPath), 0o755); err != nil {
			log.Printf("⚠️  no se pudo crear el directorio de %q: %v", localPath, err)
			continue
		}
		if err := os.WriteFile(localPath, []byte(content), 0o644); err != nil {
			log.Printf("⚠️  no se pudo escribir %q: %v", localPath, err)
			continue
		}
		e.hashes[localPath] = remoteHash
		log.Printf("📥 actualizado desde la nube: %s", localPath)
	}
	return nil
}

// handleBatch evalúa un lote de rutas modificadas y empuja a la nube solo las
// que realmente cambiaron de contenido (comparando hashes).
func (e *Engine) handleBatch(ctx context.Context, batch []string) {
	changed := make(map[string]string) // nombre remoto -> contenido

	for _, evPath := range batch {
		abs, err := filepath.Abs(evPath)
		if err != nil {
			abs = evPath
		}
		name := filepath.Base(abs)

		// Filtrar: solo nos interesan las rutas exactas configuradas.
		if want, ok := e.localByName[name]; !ok || want != abs {
			continue
		}

		h, err := hasher.HashFile(abs)
		if err != nil {
			continue // archivo borrado o ilegible
		}
		if h == e.hashes[abs] {
			continue // sin cambios reales
		}

		data, err := os.ReadFile(abs)
		if err != nil {
			continue
		}
		e.hashes[abs] = h
		changed[name] = string(data)
		log.Printf("📤 cambio detectado: %s", abs)
	}

	if len(changed) == 0 {
		return
	}

	if err := e.cloud.Push(ctx, changed); err != nil {
		log.Printf("⚠️  push falló: %v", err)
		// Invalidar los hashes para reintentar en el próximo cambio.
		for name := range changed {
			delete(e.hashes, e.localByName[name])
		}
		return
	}
	log.Printf("✅ sincronizados %d archivo(s) con la nube", len(changed))
}
