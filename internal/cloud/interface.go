package cloud

import "context"

// Provider abstrae un backend de almacenamiento en la nube. Permite
// intercambiar GitHub Gists por otros backends (S3, Dropbox, etc.) sin
// tocar el motor de sincronización.
//
// Las claves de los mapas son nombres lógicos de archivo (normalmente el
// nombre base de la ruta local); los valores son el contenido del archivo.
type Provider interface {
	// Pull descarga todos los archivos remotos y los devuelve indexados por nombre.
	Pull(ctx context.Context) (map[string]string, error)

	// Push crea o actualiza en el backend los archivos indicados.
	Push(ctx context.Context, files map[string]string) error
}
