package build

// Version es la versión del binario. Su valor real se inyecta en tiempo de
// compilación mediante -ldflags "-X .../internal/build.Version=vX.Y.Z".
// En compilaciones locales queda como "dev".
var Version = "dev"
