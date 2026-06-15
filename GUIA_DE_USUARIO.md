# 📖 Guía de Usuario — Sync Agent

Esta guía está escrita en lenguaje sencillo para que cualquier persona pueda
entender **qué hace** Sync Agent, **para qué sirve** y **cómo usarlo paso a
paso**, aunque no seas experto en programación.

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

## 4. Preparación paso a paso

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
> repositorio público. El proyecto ya está configurado para no subirlo por
> accidente (`.gitignore`).

---

## 5. Configurar el agente

1. En la carpeta del proyecto verás un archivo llamado `config.example.json`.
   Haz una copia y renómbrala a `config.json`.
2. Abre `config.json` con cualquier editor de texto y rellénalo con tus datos:

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

## 6. Poner en marcha el agente

Abre una terminal en la carpeta del proyecto y ejecuta:

```bash
# 1) Compilar el programa (solo la primera vez o tras actualizar)
go build -o bin/sync-agent ./cmd/daemon

# 2) Iniciar el agente
./bin/sync-agent --config ./config.json
```

Si todo está bien, verás mensajes como estos:

```
🚀 Iniciando Sync Agent...
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

---

## 7. Usarlo en tu segunda computadora

1. Copia el proyecto (o vuélvelo a descargar) en la otra máquina.
2. Crea ahí también su propio `config.json`, usando **el mismo `gist_token` y
   `gist_id`** (para que apunte al mismo almacén en la nube).
3. En `paths`, pon las rutas tal como existen **en esa computadora** (pueden
   ser distintas si tu usuario o carpetas cambian).
4. Compila y ejecuta igual que en el Paso 6.

Al arrancar, esta segunda máquina **descargará lo último de la nube** y dejará
tus archivos al día. ¡Listo!

> ℹ️ Los archivos se reconocen entre máquinas por su **nombre** (por ejemplo
> `.zshrc`). Por eso funciona aunque la carpeta sea diferente en cada equipo.

---

## 8. ¿Qué debes hacer tú? (resumen de tu rol)

✅ **Sí debes:**
- Crear una vez tu Gist y tu token.
- Llenar el `config.json` en cada máquina.
- Dejar el agente corriendo mientras trabajas.

❌ **No tienes que:**
- Copiar archivos manualmente entre computadoras.
- Hacer `commit`/`pull` a mano.
- Preocuparte por subir cambios: el agente lo hace solo.

---

## 9. Preguntas frecuentes y solución de problemas

**"❌ no se pudo cargar la configuración"**
→ Revisa que la ruta del `config.json` sea correcta y que el archivo exista.

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

## 10. Recordatorio de seguridad

- Tu `gist_token` es información sensible: trátalo como una contraseña.
- No lo compartas ni lo subas a repositorios públicos.
- Si crees que se filtró, ve a https://github.com/settings/tokens, bórralo y
  genera uno nuevo.

---

¿Dudas o algo no funciona como esperas? Revisa primero la sección 9; la mayoría
de los problemas se resuelven verificando el `config.json`, el token y el ID del
Gist.
