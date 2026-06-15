package engine

import (
	"context"
	"maps"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/geomark27/sync-agent/internal/config"
)

// fakeProvider simula un backend de nube en memoria, sin tocar la red.
type fakeProvider struct {
	mu       sync.Mutex
	store    map[string]string
	pushedCh chan map[string]string
}

func newFakeProvider(initial map[string]string) *fakeProvider {
	return &fakeProvider{
		store:    initial,
		pushedCh: make(chan map[string]string, 8),
	}
}

func (f *fakeProvider) Pull(ctx context.Context) (map[string]string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make(map[string]string, len(f.store))
	maps.Copy(out, f.store)
	return out, nil
}

func (f *fakeProvider) Push(ctx context.Context, files map[string]string) error {
	f.mu.Lock()
	maps.Copy(f.store, files)
	f.mu.Unlock()
	f.pushedCh <- files
	return nil
}

// TestEngineDetectaYEmpujaCambios verifica el flujo completo: el motor vigila un
// archivo real, y al modificarlo en disco lo empuja a la nube con su contenido.
func TestEngineDetectaYEmpujaCambios(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".zshrc")
	if err := os.WriteFile(path, []byte("inicial"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.AppConfig{Paths: []string{path}}
	provider := newFakeProvider(map[string]string{})
	eng := New(cfg, provider)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() {
		_ = eng.Run(ctx)
		close(done)
	}()

	// Dar tiempo a que el pull inicial y el registro del watcher se completen.
	time.Sleep(300 * time.Millisecond)

	// Modificar el archivo: debe disparar un push tras la ventana de debounce.
	if err := os.WriteFile(path, []byte("nuevo contenido"), 0o644); err != nil {
		t.Fatal(err)
	}

	select {
	case pushed := <-provider.pushedCh:
		if got := pushed[".zshrc"]; got != "nuevo contenido" {
			t.Fatalf("se empujó %q, se esperaba %q", got, "nuevo contenido")
		}
	case <-time.After(8 * time.Second):
		t.Fatal("el motor no empujó el cambio dentro del plazo")
	}

	cancel()
	<-done
}

// TestEngineAplicaPullInicial verifica que el contenido remoto se escribe en el
// archivo local cuando difieren al arrancar.
func TestEngineAplicaPullInicial(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")
	if err := os.WriteFile(path, []byte("local-viejo"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.AppConfig{Paths: []string{path}}
	provider := newFakeProvider(map[string]string{"settings.json": "remoto-nuevo"})
	eng := New(cfg, provider)

	if err := eng.initialPull(context.Background()); err != nil {
		t.Fatalf("initialPull devolvió error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "remoto-nuevo" {
		t.Fatalf("el archivo local quedó como %q, se esperaba %q", got, "remoto-nuevo")
	}
}
