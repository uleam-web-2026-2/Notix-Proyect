# Alcance y archivos del taller 10

## Alcance documentado

Esta entrega registra la comprobación de campos controlados por el servidor en `POST /tareas`, la repetición de la colección Bruno, la consulta de una tarea después de una transición inválida y una ejecución actual de las pruebas Go. Los resultados reales están en [peticiones.md](peticiones.md) y [go_test.txt](go_test.txt); la evidencia visual de Go está en [go_test_2026-11-05.png](../hito2_capturas/go_test_2026-11-05.png).

No se amplió el alcance funcional de la API en esta continuación. No se modificaron controladores, servicios, modelos, rutas, configuración, base de datos ni pruebas automatizadas.

## Archivos tocados en esta continuación

| Archivo | Cambio |
| --- | --- |
| `docs/taller10/peticiones.md` | Se agregaron las respuestas observadas para `POST /tareas`, la consulta posterior al `422` y el resultado de la colección Bruno CLI. |
| `docs/taller10/go_test.txt` | Se reemplazó la salida anterior por la salida de `RUN_DB_INTEGRATION=true go test -count=1 ./...`, ejecutado el 5 de noviembre de 2026. |
| `docs/hito2_capturas/README.md` | Se agregó la referencia a la captura actual de Go. |
| `docs/hito2_capturas/go_test_2026-11-05.png` | Se añadió una captura de pantalla de la salida real de `go test`. |
| `docs/taller10/alcance.md` | Se creó este registro del alcance y de los archivos. |

La ejecución de Bruno CLI modificó temporalmente el entorno local de la colección; `docs/hito2_bruno/environments/Local.bru` se restauró y no queda modificado.

## Archivos que quedaron como estaban en esta continuación

- `docs/hito2_hoja.md`, a petición expresa.
- `docs/hito2_bruno/environments/Local.bru` y los archivos de la colección Bruno.
- Todo el código fuente, incluidas rutas HTTP, handlers, servicios, modelos y repositorios.
- La configuración del proyecto y la base de datos.
- Los archivos de pruebas automatizadas.
- Los documentos ajenos a esta comprobación.

## Cambios previos del taller presentes en la rama

La rama ya contiene cambios de fases anteriores del taller respecto del punto de partida. Esta lista los distingue de los cinco archivos tocados en esta continuación:

```text
db/consultas.sql
db/schema.sql
docs/decisiones.md
docs/openapi.yaml
docs/taller10/redocly.txt
internal/handlers/tarea_dtos.go
internal/handlers/tarea_handler.go
internal/httpapi/router.go
internal/httpapi/router_integration_test.go
internal/models/models.go
internal/repository/filas/consultas.sql.go
internal/repository/filas/db.go
internal/repository/filas/models.go
internal/repository/tarea_sqlc.go
internal/services/tarea_repository.go
internal/services/tarea_service.go
main.go
sqlc.yaml
```

Estos cambios previos no se editaron en esta continuación.

## Verificación y límites

- La colección Bruno CLI `4.3.0` se ejecutó localmente: 7 peticiones y 7 pruebas aprobadas.
- `RUN_DB_INTEGRATION=true go test -count=1 ./...` terminó correctamente: los paquetes `internal/httpapi` e `internal/services` aprobaron; los otros paquetes reportaron `[no test files]`.
- Una ejecución inicial dentro del sandbox no pudo conectar con PostgreSQL local; se repitió con acceso local y finalizó correctamente. El fallo de sandbox no se presenta como fallo de la aplicación.
- No se consultó la tarea después de la transición prohibida hasta esta comprobación; ahora se consultó y su respuesta observada fue `200` con estado `pendiente`.
- No se ejecutó Redocly ni se generó una captura de Bruno CLI en esta continuación.
