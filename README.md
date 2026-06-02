# 🔄 Sync Agent

Un demonio (daemon) ligero y agnóstico escrito en Go, diseñado para sincronizar automáticamente archivos de configuración locales (Dotfiles, ajustes de IDEs) entre múltiples equipos sin intervención manual.

## 🎯 ¿Para qué sirve?

Modernos editores de código y herramientas (como Zed o Cursor) carecen de sincronización nativa en la nube o están limitados por sandboxes (como WebAssembly) que impiden la sincronización en segundo plano. Tradicionalmente, esto se resuelve manteniendo un repositorio Git para los "dotfiles", lo cual requiere hacer `git commit` y `git pull` manualmente en cada cambio de máquina, o depender de enlaces simbólicos (symlinks) acoplados a servicios de nube como Google Drive.

**Sync Agent** resuelve este problema subiendo un nivel: opera directamente a nivel del sistema operativo. 

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