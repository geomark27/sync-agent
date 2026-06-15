# 📖 Guía de Usuario — Sync Agent

Esta guía está escrita en lenguaje sencillo para que cualquier persona pueda
entender **qué hace** Sync Agent, **para qué sirve** y **cómo usarlo paso a
paso**, aunque no seas experto en programación.

> ✅ **No necesitas instalar Go ni saber programar.** El programa se instala
> como una aplicación lista para usar, con un solo comando.

---

## 1. ¿Qué es y para qué sirve? (en palabras simples)

Imagina que usas dos computadoras: una en la oficina y otra en casa. En ambas
tienes configurados tu editor de código, tu terminal y tus "atajos"
personales. El problema es que cada vez que cambias algo en una máquina, tienes
que copiarlo manualmente a la otra para que todo se vea igual.

**Sync Agent es un asistente invisible que hace ese trabajo por ti.**

Se queda corriendo en segundo plano, vigila los archivos de configuración que tú
elijas y, cuando guardas un cambio, lo sube automáticamente a la nube. Al
encender tu otra computadora, descarga esos cambios y deja todo idéntico. Así
tu entorno de trabajo es el mismo en todas tus máquinas, **sin que tengas que
hacer nada manualmente**.

> 💡 **En una frase:** mantiene tus archivos de configuración sincronizados
> entre varias computadoras, de forma automática.

### ¿Qué tipo de archivos puede sincronizar?
Cualquier archivo de texto, por ejemplo:
- La configuración de tu terminal (`.zshrc`, `.bashrc`).
- Los ajustes de tu editor (por ejemplo `settings.json` de Zed, Cursor o VS Code).
- Cualquier archivo de configuración personal que quieras tener igual en todos lados.

---

## 2. ¿Cómo funciona por dentro? (lo justo para entenderlo)

```
   Tu computadora                    La nube (GitHub Gist)
   ┌───────────────┐                 ┌───────────────────┐
   │  Editas un    │   ── sube ──►   │  Guarda la última │
   │  archivo      │                 │  versión          │
   └───────────────┘                 └───────────────────┘
                                              │
   Otra computadora                           │
   ┌───────────────┐                          │
   │  Al encender, │   ◄── descarga ──────────┘
   │  se actualiza │
   └───────────────┘
```

- **No necesita una base de datos ni un servidor propio.** Usa un "Gist" de
  GitHub (un espacio gratuito para guardar archivos) como almacén central.
- **Es eficiente:** solo sube algo cuando el contenido **realmente cambió**
  (lo comprueba con una "huella digital" del archivo), así no gasta red ni
  batería sin necesidad.
- **Agrupa los cambios:** muchos editores guardan varias veces en pocos
  segundos; el agente espera a que termines y sube todo de una sola vez.

---

## 3. Lo que necesitas antes de empezar (requisitos)

1. **Una cuenta de GitHub** (gratuita): https://github.com
2. **Un "Gist" secreto**, que será el lugar donde se guardan tus archivos.
3. **Un "token" de GitHub**, que es como una contraseña especial que le da
   permiso al agente para usar tu Gist.

No te preocupes si no sabes qué son: en el siguiente paso te explico cómo
conseguirlos.

---

## 4. Preparación en GitHub (paso a paso)

### Paso 4.1 — Crear el Gist (el almacén en la nube)
1. Entra a 👉 https://gist.github.com
2. En el cuadro de nombre de archivo escribe, por ejemplo: `placeholder.txt`
3. En el contenido escribe cualquier cosa, por ejemplo: `inicial`
4. Abajo a la derecha, haz clic en **"Create secret gist"** (Gist secreto).
5. Mira la dirección (URL) de tu navegador. Se verá algo así:
   `https://gist.github.com/tu-usuario/`**`a1b2c3d4e5f6...`**
   Esa parte final larga es el **ID del Gist**. Cópialo y guárdalo, lo
   necesitarás después.

### Paso 4.2 — Crear el token (el permiso de acceso)
1. Entra a 👉 https://github.com/settings/tokens
2. Haz clic en **"Generate new token"** → **"Generate new token (classic)"**.
3. Ponle un nombre que reconozcas, por ejemplo: `sync-agent`.
4. En la lista de permisos (*scopes*), marca **solo** la casilla **`gist`**.
5. Baja y haz clic en **"Generate token"**.
6. GitHub te mostrará el token **una sola vez** (empieza con `ghp_...`).
   **Cópialo y guárdalo en un lugar seguro ahora mismo**, porque no podrás
   volver a verlo.

> 🔒 **Importante sobre seguridad:** ese token es como una llave de tu cuenta.
> Nunca lo compartas, no lo envíes por chat ni correo, y no lo subas a ningún
> repositorio público.

---

## 5. Instalar el programa (un solo comando)

No necesitas descargar nada manualmente ni instalar Go. Abre una terminal y
pega el comando según tu sistema operativo:

### 🪟 Windows (PowerShell)
```powershell
iwr -useb https://raw.githubusercontent.com/geomark27/sync-agent/main/scripts/install.ps1 | iex
```

### 🐧 Linux / 🍎 macOS (Terminal)
```bash
curl -fsSL https://raw.githubusercontent.com/geomark27/sync-agent/main/scripts/install.sh | bash
```

Esto descarga el programa ya listo, lo instala y lo deja disponible como el
comando `sync-agent`. **Cierra y vuelve a abrir la terminal** al terminar (para
que reconozca el nuevo comando).

> ✅ Para comprobar que quedó instalado, escribe: `sync-agent version`

---

## 6. Configurar el agente

### Paso 6.1 — Crear el archivo de configuración
Escribe en la terminal:

```bash
sync-agent init
```

Esto crea un archivo de configuración vacío y te dice exactamente dónde quedó
(por ejemplo, en Linux: `~/.config/sync-agent/config.json`).

### Paso 6.2 — Rellenar tus datos
Abre ese archivo con cualquier editor de texto y complétalo así:

```json
{
  "machine_id": "mi-laptop-oficina",
  "paths": [
    "/home/usuario/.zshrc",
    "/home/usuario/.config/zed/settings.json"
  ],
  "gist_token": "ghp_aqui_va_tu_token",
  "gist_id": "aqui_va_el_id_de_tu_gist"
}
```

**Qué significa cada campo:**

| Campo | Qué debes poner |
|-------|-----------------|
| `machine_id` | Un nombre para identificar esta computadora (lo eliges tú, ej. `laptop-casa`). |
| `paths` | La lista de archivos que quieres sincronizar, con su **ruta completa**. |
| `gist_token` | El token que creaste en el Paso 4.2. |
| `gist_id` | El ID del Gist que copiaste en el Paso 4.1. |

> 💡 **¿Cómo saber la ruta completa de un archivo?** En la terminal, ve a la
> carpeta del archivo y escribe `pwd` (te muestra la ruta de la carpeta); la
> ruta completa es esa carpeta + `/` + el nombre del archivo.

---

## 7. Poner en marcha el agente

Escribe simplemente:

```bash
sync-agent
```

Si todo está bien, verás mensajes como estos:

```
🚀 Iniciando Sync Agent v1.0.0...
📥 actualizado desde la nube: /home/usuario/.zshrc      (si había algo más nuevo en la nube)
👀 Vigilando 2 archivo(s)...
```

A partir de aquí, **déjalo corriendo**. Cada vez que guardes un cambio en
alguno de tus archivos vigilados, verás algo como:

```
📤 cambio detectado: /home/usuario/.zshrc
✅ sincronizados 1 archivo(s) con la nube
```

Para **detenerlo de forma segura**, pulsa `Ctrl + C`.

> 💡 **¿Quieres que arranque solo y corra siempre en segundo plano?** Mira la
> sección 11 ("Dejarlo corriendo automáticamente").

---

## 8. Usarlo en tu segunda computadora

1. Instala el programa en la otra máquina con el mismo comando del Paso 5.
2. Ejecuta `sync-agent init` y edita su `config.json` usando **el mismo
   `gist_token` y `gist_id`** (para que apunte al mismo almacén en la nube).
3. En `paths`, pon las rutas tal como existen **en esa computadora** (pueden
   ser distintas si tu usuario o carpetas cambian).
4. Ejecuta `sync-agent`.

Al arrancar, esta segunda máquina **descargará lo último de la nube** y dejará
tus archivos al día. ¡Listo!

> ℹ️ Los archivos se reconocen entre máquinas por su **nombre** (por ejemplo
> `.zshrc`). Por eso funciona aunque la carpeta sea diferente en cada equipo.

---

## 9. ¿Qué debes hacer tú? (resumen de tu rol)

✅ **Sí debes:**
- Crear una vez tu Gist y tu token (Paso 4).
- Instalar con un comando (Paso 5).
- Llenar el `config.json` en cada máquina (Paso 6).
- Dejar el agente corriendo mientras trabajas (Paso 7).

❌ **No tienes que:**
- Instalar Go ni compilar nada.
- Copiar archivos manualmente entre computadoras.
- Hacer `commit`/`pull` a mano.
- Preocuparte por subir cambios: el agente lo hace solo.

---

## 10. Preguntas frecuentes y solución de problemas

**`sync-agent: command not found` (o "no se reconoce el comando")**
→ Cierra y vuelve a abrir la terminal después de instalar. En Linux/Mac también
puedes ejecutar `source ~/.zshrc` (o `source ~/.bashrc`).

**"❌ no se pudo cargar la configuración"**
→ Aún no creaste el archivo. Ejecuta `sync-agent init`.

**"❌ configuración incompleta: 'gist_token' y 'gist_id' son obligatorios"**
→ Faltó llenar el token o el ID del Gist en tu `config.json`.

**"push falló (401 ...)" o "(404 ...)"**
→ El token es incorrecto/expiró (401) o el ID del Gist está mal (404). Vuelve
a revisar los Pasos 4.1 y 4.2.

**Cambié un archivo y no se subió.**
→ El agente espera ~2 segundos tras tu último guardado antes de subir (para
agrupar cambios). Espera un momento. Además, asegúrate de que el archivo esté
listado en `paths`.

**¿Qué pasa si edito el mismo archivo en las dos máquinas a la vez?**
→ El agente es sencillo en este aspecto: gana el último cambio que detecta
(no combina cambios). Para evitar perder trabajo, edita en una máquina a la
vez cuando sea posible.

**¿Consume muchos recursos?**
→ No. Está hecho en Go y diseñado para ser muy ligero en CPU y memoria.

---

## 11. Dejarlo corriendo automáticamente (opcional, avanzado)

Por defecto, `sync-agent` corre mientras la terminal esté abierta. Si quieres
que arranque solo al encender el equipo y siga en segundo plano:

- **Linux (systemd):** crea un servicio de usuario en
  `~/.config/systemd/user/sync-agent.service` que ejecute la ruta del binario
  (`~/.local/bin/sync-agent`) y actívalo con
  `systemctl --user enable --now sync-agent`.
- **macOS (launchd):** crea un archivo `.plist` en `~/Library/LaunchAgents/`
  que apunte al binario y cárgalo con `launchctl load`.
- **Windows:** usa el **Programador de tareas** y crea una tarea que ejecute
  `sync-agent.exe` "Al iniciar sesión".

Si no estás familiarizado con esto, no pasa nada: basta con abrir una terminal
y ejecutar `sync-agent` cuando empieces a trabajar.

---

## 12. Recordatorio de seguridad

- Tu `gist_token` es información sensible: trátalo como una contraseña.
- No lo compartas ni lo subas a repositorios públicos.
- Si crees que se filtró, ve a https://github.com/settings/tokens, bórralo y
  genera uno nuevo.

---

## 13. Para desarrolladores (compilar desde el código)

Esta sección **NO es necesaria para usar el programa**; solo aplica si quieres
modificar el código o generar tú mismo los ejecutables. Requiere tener
[Go](https://go.dev/dl/) instalado.

```bash
make build      # Compila el binario en ./bin/sync-agent
make run        # Ejecuta en modo desarrollo
make test       # Corre las pruebas
make build-all  # Compila para Linux, Windows y Mac
make release    # Publica una nueva versión en GitHub Releases
```

Una vez publicada una versión con `make release`, los comandos de instalación
del Paso 5 (los *one-liners*) descargarán automáticamente ese ejecutable.

---

¿Dudas o algo no funciona como esperas? Revisa primero la sección 10; la mayoría
de los problemas se resuelven verificando el `config.json`, el token y el ID del
Gist.
