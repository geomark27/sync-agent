package watcher

import (
	"context"
	"sort"
	"testing"
	"time"
)

func TestDebounceAgrupaYDeduplica(t *testing.T) {
	in := make(chan string)
	out := Debounce(t.Context(), in, 30*time.Millisecond)

	// Ráfaga rápida con un duplicado: debe colapsar en un único lote.
	in <- "a"
	in <- "b"
	in <- "a"

	select {
	case batch := <-out:
		sort.Strings(batch)
		if len(batch) != 2 || batch[0] != "a" || batch[1] != "b" {
			t.Fatalf("lote inesperado: %v", batch)
		}
	case <-time.After(time.Second):
		t.Fatal("el debouncer no emitió el lote a tiempo")
	}
}

func TestDebounceCierraAlCancelar(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan string)
	out := Debounce(ctx, in, time.Hour)

	cancel()

	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("se esperaba que el canal estuviera cerrado")
		}
	case <-time.After(time.Second):
		t.Fatal("el debouncer no cerró el canal tras cancelar el contexto")
	}
}
