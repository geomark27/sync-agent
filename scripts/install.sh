#!/usr/bin/env bash
# ================================================================
#  Sync Agent - Instalador para Linux / Mac
#  Uso: curl -fsSL https://raw.githubusercontent.com/geomark27/sync-agent/main/scripts/install.sh | bash
# ================================================================

set -e

REPO="geomark27/sync-agent"
INSTALL_DIR="$HOME/.local/bin"

# Detectar sistema operativo
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux)  ASSET="sync-agent-linux-amd64"  ;;
  darwin) ASSET="sync-agent-darwin-amd64" ;;
  *)
    echo "Sistema operativo no soportado: $OS"
    exit 1
    ;;
esac

echo ""
echo "  Sync Agent - Instalador"
echo "  ────────────────────────────────────"
echo ""

# 1. Obtener última versión
echo "[1/3] Buscando ultima version..."
VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
  | grep '"tag_name"' \
  | cut -d'"' -f4)
echo "      version: $VERSION"

# 2. Crear directorio y descargar binario
mkdir -p "$INSTALL_DIR"
echo "[2/3] Descargando $ASSET..."
URL="https://github.com/$REPO/releases/download/$VERSION/$ASSET"
curl -fsSL "$URL" -o "$INSTALL_DIR/sync-agent"
chmod +x "$INSTALL_DIR/sync-agent"
echo "      descargado en $INSTALL_DIR/sync-agent"

# 3. Verificar PATH
echo "[3/3] Verificando PATH..."
if echo "$PATH" | grep -q "$INSTALL_DIR"; then
  echo "      ya estaba en PATH"
else
  echo "      agregando $INSTALL_DIR al PATH..."
  SHELL_RC=""
  if [ -f "$HOME/.zshrc" ]; then
    SHELL_RC="$HOME/.zshrc"
  elif [ -f "$HOME/.bashrc" ]; then
    SHELL_RC="$HOME/.bashrc"
  fi

  if [ -n "$SHELL_RC" ]; then
    echo "" >> "$SHELL_RC"
    echo "# sync-agent" >> "$SHELL_RC"
    echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$SHELL_RC"
    echo "      agregado a $SHELL_RC"
  else
    echo "      agrega manualmente a tu shell: export PATH=\"\$PATH:$INSTALL_DIR\""
  fi
fi

# Resultado final
echo ""
echo "  Sync Agent $VERSION instalado correctamente"
echo ""
echo "  Proximos pasos:"
echo "    1. Ejecuta: source ~/.zshrc  (o abre una nueva terminal)"
echo "    2. Ejecuta: sync-agent init"
echo "    3. Edita el config.json (gist_token, gist_id, paths)"
echo "    4. Ejecuta: sync-agent"
echo ""
