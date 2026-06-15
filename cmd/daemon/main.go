package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/geomark27/sync-agent/internal/build"
	"github.com/geomark27/sync-agent/internal/cloud"
	"github.com/geomark27/sync-agent/internal/config"
	"github.com/geomark27/sync-agent/internal/engine"
)

func main() {
	// Subcomandos sencillos: sync-agent <init|version|help>
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			runInit()
			return
		case "version", "--version", "-v":
			fmt.Printf("sync-agent %s\n", build.Version)
			return
		case "help", "--help", "-h":
			printUsage()
			return
		}
	}

	runDaemon()
}

// runDaemon arranca el agente de sincronización.
func runDaemon() {
	cfgPath := flag.String("config", defaultConfigPath(), "ruta al archivo de configuración JSON")
	flag.Parse()

	log.Printf("🚀 Iniciando Sync Agent %s...", build.Version)

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		log.Fatalf("❌ no se pudo cargar la configuración (%s): %v\n   Sugerencia: ejecuta 'sync-agent init' para crear una.", *cfgPath, err)
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

// runInit crea un archivo de configuración de ejemplo en la ruta por defecto.
func runInit() {
	path := defaultConfigPath()
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("ℹ️  Ya existe una configuración en:\n   %s\n", path)
		return
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		log.Fatalf("❌ no se pudo crear el directorio: %v", err)
	}

	const plantilla = `{
  "machine_id": "mi-equipo",
  "paths": [],
  "gist_token": "",
  "gist_id": ""
}
`
	if err := os.WriteFile(path, []byte(plantilla), 0o600); err != nil {
		log.Fatalf("❌ no se pudo escribir la configuración: %v", err)
	}

	fmt.Printf("✅ Configuración creada en:\n   %s\n\n", path)
	fmt.Println("Siguientes pasos:")
	fmt.Println("  1. Edita ese archivo y completa: gist_token, gist_id y paths.")
	fmt.Println("  2. Ejecuta: sync-agent")
}

func printUsage() {
	fmt.Println(`Sync Agent — sincroniza archivos de configuración entre equipos.

Uso:
  sync-agent                   Inicia el agente (configuración por defecto)
  sync-agent --config <ruta>   Inicia con una configuración específica
  sync-agent init              Crea un archivo de configuración de ejemplo
  sync-agent version           Muestra la versión
  sync-agent help              Muestra esta ayuda`)
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
