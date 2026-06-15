package watcher

import (
	"context"
	"log"

	"github.com/fsnotify/fsnotify"
)

// Watcher envuelve fsnotify y emite las rutas modificadas a través de un canal.
type Watcher struct {
	fsw    *fsnotify.Watcher
	events chan string
}

// New crea un Watcher listo para registrar rutas.
func New() (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &Watcher{
		fsw:    fsw,
		events: make(chan string, 64),
	}, nil
}

// Add registra una ruta (archivo o directorio) para vigilancia. Para detectar
// guardados atómicos de los editores conviene vigilar el directorio padre.
func (w *Watcher) Add(path string) error {
	return w.fsw.Add(path)
}

// Events expone el canal de rutas modificadas.
func (w *Watcher) Events() <-chan string {
	return w.events
}

// Run bombea los eventos de fsnotify hacia el canal de salida hasta que el
// contexto se cancela. Solo propaga escrituras y creaciones.
func (w *Watcher) Run(ctx context.Context) {
	defer w.fsw.Close()
	defer close(w.events)

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-w.fsw.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				select {
				case w.events <- event.Name:
				case <-ctx.Done():
					return
				}
			}
		case err, ok := <-w.fsw.Errors:
			if !ok {
				return
			}
			log.Printf("⚠️  error del watcher: %v", err)
		}
	}
}
