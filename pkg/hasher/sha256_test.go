package hasher

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHashBytes(t *testing.T) {
	// SHA-256 conocido de la cadena vacía.
	const emptySum = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if got := HashBytes([]byte{}); got != emptySum {
		t.Fatalf("HashBytes(\"\") = %s, se esperaba %s", got, emptySum)
	}

	// Misma entrada -> mismo hash; entradas distintas -> hashes distintos.
	hola1 := HashBytes([]byte("hola"))
	hola2 := HashBytes([]byte("hola"))
	adios := HashBytes([]byte("adios"))
	if hola1 != hola2 {
		t.Fatal("el hash no es determinista")
	}
	if hola1 == adios {
		t.Fatal("entradas distintas produjeron el mismo hash")
	}
}

func TestHashFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(path, []byte("contenido"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := HashFile(path)
	if err != nil {
		t.Fatalf("HashFile devolvió error: %v", err)
	}
	if want := HashBytes([]byte("contenido")); got != want {
		t.Fatalf("HashFile = %s, se esperaba %s", got, want)
	}

	if _, err := HashFile(filepath.Join(dir, "no-existe")); err == nil {
		t.Fatal("se esperaba error al hashear un archivo inexistente")
	}
}
