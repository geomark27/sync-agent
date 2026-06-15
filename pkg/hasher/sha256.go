package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// HashBytes devuelve el hash SHA-256 (hexadecimal) de los datos dados.
func HashBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// HashFile calcula el hash SHA-256 de un archivo leyéndolo en streaming,
// evitando cargarlo completo en memoria.
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
