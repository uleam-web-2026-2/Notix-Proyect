# Semana 12 · auditoría de inyección

Búsqueda: `git grep -n -E -f docs/taller12/patrones.txt -- "*.go" ":(exclude)*_test.go"`.
Salida completa: `auditoria.txt` (22 coincidencias).

| Línea | ¿Arma SQL? | ¿Llega texto del cliente? | ¿Va aparte? | Qué hicimos |
| --- | --- | --- | --- | --- |
| `database/database.go:30` `dsn.Query()` | No. Arma parámetros de la URL de conexión, no SQL. | No. Sale de configuración local. | No aplica. | Nada. |
| `database/database.go:83` `Where("creador_id = ? AND titulo = ?", ...)` | Sí. | No. Son datos de siembra definidos por el servidor. | Sí, ambos valores viajan en `?`. | Nada: ya usa parámetros. |
| `database/database.go:91` `Where("tablero_id = ? AND orden = ?", ...)` | Sí. | No. Son datos de siembra definidos por el servidor. | Sí, ambos valores viajan en `?`. | Nada: ya usa parámetros. |
| `database/database.go:104` `Where("tablero_id = ? AND usuario_id = ? AND estado_invitacion = ?", ...)` | Sí. | No. Son datos de siembra definidos por el servidor. | Sí, los tres valores viajan en `?`. | Nada: ya usa parámetros. |
| `database/database.go:117` `Where("columna_id = ? AND titulo = ?", ...)` | Sí. | No. Son datos de siembra definidos por el servidor. | Sí, ambos valores viajan en `?`. | Nada: ya usa parámetros. |
| `database/database.go:132` `Where("correo = ?", email)` | Sí. | No. `email` es una cuenta de demostración definida por el servidor. | Sí, con `?`. | Nada: ya usa parámetros. |
| `database/database.go:142` `Where("correo = ?", email)` | Sí. | No. `email` es una cuenta de demostración definida por el servidor. | Sí, con `?`. | Nada: ya usa parámetros. |
| `internal/handlers/tarea_handler.go:48` `h.tasks.Find(...)` | No. Es una llamada al servicio. | El ID llegó por la ruta, pero `parseID` ya lo convirtió a `uint`. | No aplica. | Nada. |
| `internal/handlers/tarea_handler.go:105` `name + ...` | No. Solo arma el mensaje de error. | No. `name` es una etiqueta fija del handler. | No aplica. | Nada. |
| `internal/repository/access_repository.go:22` `First(&board, boardID)` | Sí. | El ID de tablero puede venir de la ruta y ya fue convertido a `uint`. | Sí. GORM recibe el ID como valor. | Nada: el ID se valida antes. |
| `internal/repository/access_repository.go:32` `Where("tablero_id = ? AND usuario_id = ? AND estado_invitacion = ?", ...)` | Sí. | Indirectamente, los IDs validados y el usuario del token. | Sí, los tres valores usan `?`. | Nada: ya usa parámetros. |
| `internal/repository/access_repository.go:33` `First(&share)` | Sí. Ejecuta la consulta construida antes. | No agrega texto del cliente. | Los valores anteriores ya van como parámetros. | Nada. |
| `internal/repository/account_repository.go:22` `Where("correo = ?", email).First(&account)` | Sí. | Sí, el correo llega desde el cuerpo de `POST /login`. | Sí, con `?`. | Nada: ya usa parámetros. |
| `internal/repository/tarea_repository.go:24` `Where("columnas.tablero_id = ?", boardID)` | Sí. | El ID de tablero puede venir de la ruta y ya fue convertido a `uint`. | Sí, con `?`. | Nada: ya usa parámetros. |
| `internal/repository/tarea_repository.go:27` `Order("columnas.orden ASC, tareas.creado ASC, tareas.id ASC")` | Sí. | No. El orden es texto fijo del servidor. | No aplica. | Nada: no recibe un campo de orden del cliente. |
| `internal/repository/tarea_repository.go:28` `Find(&tasks)` | Sí. Ejecuta la consulta construida antes. | No agrega texto del cliente. | Los valores anteriores ya van como parámetros. | Nada. |
| `internal/repository/tarea_repository.go:40` `First(&task, id)` | Sí. | El ID de tarea puede venir de la ruta y ya fue convertido a `uint`. | Sí. GORM recibe el ID como valor. | Nada: el ID se valida antes. |
| `internal/repository/tarea_repository.go:59` `Where("id = ? AND estado = ?", id, from)` | Sí. | El ID ya fue convertido a `uint`; el estado fue validado contra transiciones permitidas. | Sí, ambos valores usan `?`. | Nada: ya usa parámetros. |
| `internal/repository/tarea_repository.go:69` `Select("count(*) > 0").Where("id = ?", id)` | Sí. | El ID ya fue convertido a `uint`. | Sí, el ID usa `?`; el `Select` es texto fijo del servidor. | Nada. |
| `internal/repository/tarea_repository.go:80` `First(&column, id)` | Sí. | El ID de columna llega del JSON como `uint` y se valida como positivo. | Sí. GORM recibe el ID como valor. | Nada: el ID se valida antes. |
| `internal/services/tarea_service.go:39` `Find(...)` | No. Es una llamada al repositorio. | El ID ya fue convertido a `uint` en el handler. | No aplica. | Nada. |
| `main.go:41` `":" + cfg.Port` | No. Arma la dirección HTTP del servidor. | No. Sale de configuración local. | No aplica. | Nada. |

**Resultado:** se revisaron las 22 coincidencias. Las 16 operaciones que construyen o ejecutan SQL no concatenan texto del cliente: los valores variables usan parámetros de GORM (`?`) o llegan como identificadores numéricos ya convertidos. No hubo una consulta vulnerable que cerrar.
