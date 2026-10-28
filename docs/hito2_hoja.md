# Hito 2 · Hoja de verificación de Notix

> Esta hoja describe el backend que está en el repositorio. La imagen de referencia usa una API de tickets; Notix implementa tableros y tareas, por lo que aquí se documentan sus rutas y roles reales. Los resultados esperados se distinguen de los resultados realmente ejecutados.

## 1. Usuarios de ejemplo

Las cuentas se crean idempotentemente al iniciar la aplicación, mediante `database.SeedDemoData`. El login compara hashes bcrypt guardados en la tabla `cuentas`.

| Usuario (correo) | Rol | Contraseña de demostración | Evidencia |
| --- | --- | --- | --- |
| `propietaria@notix.test` | Propietario (`propietario`) | `Prop-demo-2026` | `database/database.go`, `README.md` |
| `colaborador@notix.test` | Colaborador (`colaborador`) | `Colab-demo-2026` | `database/database.go`, `README.md` |

Son credenciales ficticias de desarrollo/demostración incluidas en el código; no son claves personales. La contraseña se almacena como hash bcrypt. No se implementa registro de cuentas por API.

## 2. Las seis peticiones

La colección Bruno está en `docs/hito2_bruno/`. Primero se ejecuta `0 Preparar demo.bru`, que obtiene tokens recientes y prepara variables/IDs para la secuencia. Las rutas de las filas son las implementadas. Los resultados esperados se derivan de las pruebas declaradas en la colección y el código. **No ejecuté la colección en este entorno**, por lo que el resultado real queda pendiente.

| # | Caso y petición | Credencial necesaria | Resultado esperado | Resultado real |
| --- | --- | --- | --- | --- |
| 1 | `GET /tableros/{boardId}/tareas` sin token (`1 Sin token.bru`) | Ninguna | `401`, error `unauthorized` | Pendiente de ejecutar |
| 2 | `POST /login` con el usuario colaborador (`2 Login colaborador.bru`) | Correo y contraseña del colaborador, en el JSON | `200`; devuelve token Bearer y rol `colaborador` | Pendiente de ejecutar |
| 3 | `GET /tableros/{boardId}/tareas` con token de colaborador (`3 Listar tareas.bru`) | Bearer de colaborador con invitación activa al tablero demo | `200`; devuelve tareas del tablero | Pendiente de ejecutar |
| 4 | `PATCH /tareas/{taskId}/estado` solicitando `en_proceso` (`4 Cambio colaborador 403.bru`) | Bearer de colaborador | `403`; solo propietario puede cambiar estado | Pendiente de ejecutar |
| 5 | `PATCH /tareas/{taskId}/estado` de `pendiente` a `en_proceso` (`5 Cambio propietario 200.bru`) | Bearer de propietario | `200`; devuelve la tarea con estado actualizado | Pendiente de ejecutar |
| 6 | `PATCH /tareas/{forbiddenTaskId}/estado` de `pendiente` a `completado` (`6 Transicion prohibida 422.bru`) | Bearer de propietario | `422`, error `invalid_transition`; la tarea reservada permanece en `pendiente` | Pendiente de ejecutar |

Los valores `boardId`, `taskId` y `forbiddenTaskId` se resuelven con los datos de demostración. No se debe asumir que siempre sean `1`: al reutilizar una base de datos los IDs pueden cambiar y el servidor los imprime al iniciar.

## 3. Matriz ruta × rol

`Propietario` y `colaborador` se refieren a cuentas autenticadas. En rutas de tareas, los resultados de acceso autorizado presuponen que la cuenta pertenece al tablero. El middleware devuelve `401` si falta un token válido. Para usuarios autenticados sin acceso al tablero, las consultas y la creación devuelven `403`.

| Ruta realmente registrada | Sin sesión | Propietario con acceso | Colaborador con acceso |
| --- | --- | --- | --- |
| `GET /health` | `200` | `200` | `200` |
| `POST /login` | `200` con credenciales válidas; `401` si no coinciden | Igual; requiere credenciales en el cuerpo | Igual; requiere credenciales en el cuerpo |
| `GET /tableros/{tableroID}/tareas` | `401` | `200`; `403` si no tiene acceso | `200`; `403` si no tiene acceso |
| `GET /tareas/{id}` | `401` | `200`; `403` si no tiene acceso | `200`; `403` si no tiene acceso |
| `POST /tareas` | `401` | `201`; `403` si no tiene acceso a la columna/tablero | `201`; `403` si no tiene acceso a la columna/tablero |
| `PATCH /tareas/{id}/estado` | `401` | `200` si la transición es válida | `403` |

IDs inválidos producen `400`; recursos inexistentes producen `404`; validaciones de tarea o transiciones no permitidas producen `422`. No hay rutas `PUT` ni `DELETE` registradas para estos recursos. La fuente de esta matriz es `internal/httpapi/router.go`; las operaciones marcadas como plan en `docs/openapi.yaml` no forman parte de la matriz.

## 4. Recuadro de cumplimiento

Criterios marcados como cumplidos solo cuando se verificaron en archivos del repositorio. La existencia de pruebas o aserciones no acredita una ejecución exitosa.

- [x] Los structs de las entidades principales están definidos en `internal/models/models.go`; sus nombres de tabla están definidos y GORM ejecuta `AutoMigrate` en `database/database.go`.
- [x] Hay validación en el servicio: título requerido con máximo 150 caracteres, contenido de hasta 5000, identificador positivo de columna y transiciones de estado permitidas. El decodificador rechaza JSON mal formado, claves desconocidas y cuerpos con más de un objeto.
- [x] Existe una colección Bruno con los seis casos de esta hoja en `docs/hito2_bruno/`; los archivos declaran sus códigos esperados.
- [ ] Resultado real de las seis peticiones Bruno: pendiente de ejecutar y registrar.
- [x] Existe una prueba de integración de la API en `internal/httpapi/router_integration_test.go`, que cubre autenticación, permisos, creación/consulta y transiciones con PostgreSQL.
- [ ] Ejecución exitosa de `go test ./...`: no comprobada en esta revisión. La prueba de integración puede omitirse si PostgreSQL no está disponible.
- [x] El código separa modelos, handlers, servicios, repositorios y acceso HTTP; esta organización se explica en `docs/decisiones.md`.
- [ ] Capturas de respaldo de cada paso: no localizadas entre los archivos revisados; pendientes de agregar/verificar.

## Fuentes revisadas

`internal/httpapi/router.go`, `internal/handlers/`, `internal/middleware/auth.go`, `internal/services/`, `internal/repository/`, `internal/models/models.go`, `database/database.go`, `internal/httpapi/router_integration_test.go`, `README.md` y `docs/hito2_bruno/`.
