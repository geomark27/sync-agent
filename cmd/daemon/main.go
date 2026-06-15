package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/geomark27/sync-agent/internal/cloud"
	"github.com/geomark27/sync-agent/internal/config"
	"github.com/geomark27/sync-agent/internal/engine"
)

func main() {
	log.Println("🚀 Iniciando Sync Agent...")

	cfgPath := flag.String("config", defaultConfigPath(), "ruta al archivo de configuración JSON")
	flag.Parse()

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		log.Fatalf("❌ no se pudo cargar la configuración (%s): %v", *cfgPath, err)
	}
	if cfg.GistToken == "" || cfg.GistID == "" {
		log.Fatal("❌ configuración incompleta: 'gist_token' y 'gist_id' son obligatorios")
	}
	if len(cfg.Paths) == 0 {
		log.Fatal("❌ configuración incompleta: define al menos una ruta en 'paths'")
	}

	// Contexto cancelable por señales del sistema operativo.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	provider := cloud.NewGistProvider(cfg.GistToken, cfg.GistID)
	eng := engine.New(cfg, provider)

	errCh := make(chan error, 1)
	go func() {
		errCh <- eng.Run(ctx)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			log.Fatalf("❌ el agente terminó con error: %v", err)
		}
	case <-ctx.Done():
		log.Println("🛑 Señal recibida. Apagando Sync Agent de forma segura...")
		<-errCh // esperar a que el motor cierre limpiamente
	}
}

// defaultConfigPath devuelve la ruta de configuración por defecto dentro del
// directorio de configuración del usuario (p. ej. ~/.config/sync-agent/).
func defaultConfigPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "config.json"
	}
	return filepath.Join(dir, "sync-agent", "config.json")
}
