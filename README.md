# 🔄 Sync Agent

Un demonio (daemon) ligero y agnóstico escrito en Go, diseñado para sincronizar automáticamente archivos de configuración locales (Dotfiles, ajustes de IDEs) entre múltiples equipos sin intervención manual.

## 🎯 ¿Para qué sirve?

Modernos editores de código y herramientas (como Zed o Cursor) carecen de sincronización nativa en la nube o están limitados por sandboxes (como WebAssembly) que impiden la sincronización en segundo plano. Tradicionalmente, esto se resuelve manteniendo un repositorio Git para los "dotfiles", lo cual requiere hacer `git commit` y `git pull` manualmente en cada cambio de máquina, o depender de enlaces simbólicos (symlinks) acoplados a servicios de nube como Google Drive.

**Sync Agent** resuelve este problema subiendo un nivel: opera directamente a nivel del sistema operativo. 

> 📖 ¿Eres usuario y buscas una explicación clara paso a paso? Lee la **[Guía de Usuario](GUIA_DE_USUARIO.md)**.

Vigila de forma silenciosa las rutas que tú definas, detecta cuando guardas un cambio en tus configuraciones e impacta esos cambios automáticamente en un almacenamiento centralizado (como GitHub Gists). Cuando enciendes tu otro equipo, el agente local descarga los cambios y actualiza tus archivos, permitiendo que tu entorno de desarrollo sea idéntico en cualquier máquina con cero fricción.

## ✨ Características Principales

* **Agnóstico al IDE:** No depende de las APIs de Zed, Cursor, VS Code o la terminal (Zsh/Powerlevel10k). Funciona observando el sistema de archivos, por lo que puedes sincronizar lo que quieras.
* **Extremadamente Ligero:** Construido en Go puro. Aprovecha goroutines para el manejo concurrente, manteniendo el consumo de CPU y RAM al mínimo mientras corre en segundo plano.
* **I/O Optimizado:** Implementa mecanismos de *Debouncing* (agrupación de eventos de escritura) y verificación de estado local mediante *Hashes SHA-256*, asegurando que solo se realicen llamadas de red o escrituras en disco cuando el contenido realmente ha cambiado.
* **Cloud-Backed (GitHub Gists):** Utiliza Gists secretos como backend de almacenamiento por defecto, lo que proporciona historial de versiones y control de acceso seguro sin necesidad de levantar una base de datos propia.

## 🏗️ Arquitectura

El agente está compuesto por módulos independientes que se comunican de forma asíncrona mediante canales:

1. **Config Manager:** Carga las rutas a observar y las credenciales desde un archivo local seguro.
2. **FS Watcher:** Utiliza `fsnotify` para escuchar eventos nativos del sistema operativo en tiempo real.
3. **Debouncer / Batcher:** Filtra y agrupa las múltiples señales de guardado que emiten los IDEs.
4. **State Tracker:** Calcula y compara los hashes de los archivos.
5. **Cloud Transport:** Gestiona la comunicación bidireccional (Push/Pull) con la API de la nube.

## 🚀 Estructura del Proyecto

El código sigue el estándar de diseño de proyectos de Go:

```text
/sync-agent
├── cmd/daemon/         # Punto de entrada de la aplicación
├── internal/
│   ├── config/         # Carga de configuraciones
│   ├── watcher/        # Lógica de observación de archivos e I/O
│   ├── cloud/          # Clientes de APIs externas (GitHub)
│   └── engine/         # Orquestador y resolución de sincronización
└── pkg/hasher/         # Utilidades criptográficas compartidas
```

## ⚡ Instalación rápida (usuario final, sin Go)

El programa se distribuye como un ejecutable listo para usar. **No requiere
instalar Go.** Instálalo con un solo comando:

**Windows (PowerShell):**
```powershell
iwr -useb https://raw.githubusercontent.com/geomark27/sync-agent/main/scripts/install.ps1 | iex
```

**Linux / macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/geomark27/sync-agent/main/scripts/install.sh | bash
```

Luego: `sync-agent init` → edita el `config.json` → `sync-agent`.

> 📖 Explicación detallada paso a paso en la **[Guía de Usuario](GUIA_DE_USUARIO.md)**.
> Los *one-liners* descargan el binario del último [release](https://github.com/geomark27/sync-agent/releases); requieren que exista al menos una versión publicada (`make release`).

## ⚙️ Configuración

Crea un archivo `config.json` (puedes partir de `config.example.json`). **Nunca lo subas al repositorio**: contiene tu token personal (ya está incluido en `.gitignore`).

```json
{
  "machine_id": "laptop-personal",
  "paths": [
    "/home/usuario/.zshrc",
    "/home/usuario/.config/zed/settings.json"
  ],
  "gist_token": "ghp_tu_token_personal",
  "gist_id": "id_del_gist_secreto"
}
```

| Campo | Descripción |
|-------|-------------|
| `machine_id` | Identificador legible de la máquina (informativo). |
| `paths` | Lista de **rutas absolutas** a los archivos que quieres sincronizar. |
| `gist_token` | Token personal de GitHub con permiso (scope) `gist`. |
| `gist_id` | ID de un Gist (preferiblemente secreto) que actúa como almacenamiento. |

> **Requisitos previos:** crea un [Gist](https://gist.github.com) (puede ser secreto) y genera un [token personal](https://github.com/settings/tokens) con el scope `gist`.

## 🚀 Uso

```bash
sync-agent init               # Crea un config.json de ejemplo en la ruta por defecto
sync-agent                    # Inicia el agente (usa ~/.config/sync-agent/config.json)
sync-agent --config ./config.json   # Inicia con una configuración específica
sync-agent version            # Muestra la versión instalada
```

Al arrancar, el agente:
1. **Descarga** el estado remoto del Gist y actualiza los archivos locales que difieran (en el arranque, la nube tiene prioridad).
2. **Vigila** los directorios de los archivos configurados.
3. Ante un cambio, agrupa los eventos (*debounce* de 2 s), verifica el hash SHA-256 y **sube** a la nube solo lo que realmente cambió.

Detén el agente de forma segura con `Ctrl+C` (`SIGINT`/`SIGTERM`).

## 🛠️ Para desarrolladores (compilar desde el código)

Requiere [Go](https://go.dev/dl/). El `Makefile` automatiza las tareas comunes:

```bash
make build      # Compila el binario en ./bin/sync-agent
make run        # Ejecuta en modo desarrollo (make run ARGS='--config ./config.json')
make test       # Ejecuta las pruebas
make lint       # go fmt + go vet
make build-all  # Compila para Linux, Windows y macOS (amd64)
make release    # lint + test + compila + tag + publica binarios en GitHub Releases
```

> `make release` sube los binarios al *release* de GitHub; los scripts de
> instalación rápida descargan justamente esos artefactos.

## ⚠️ Limitaciones actuales

* **Identificación por nombre base:** los archivos se mapean entre máquinas por su nombre de archivo (`.zshrc`, `settings.json`). Dos rutas con el mismo nombre base colisionan (el agente lo avisa y omite la segunda).
* **Resolución de conflictos sencilla:** en el arranque gana la versión de la nube; durante la ejecución gana el último cambio local detectado (*last-write-wins*). No hay fusión (*merge*) de cambios concurrentes.
* **Sin sondeo remoto continuo:** los cambios remotos se aplican al arrancar, no en caliente mientras el agente corre.