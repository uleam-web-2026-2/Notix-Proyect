# Notix Backend

API Go para el Hito 2. `Tarea` es la entidad principal; `Usuario`, `Tablero`,
`Columna`, `TableroCompartido` y `Comentario` se migran junto con sus
relaciones. Las cuentas de inicio de sesión se mantienen en una tabla separada.

## Requisitos

- Go 1.23 o posterior.
- Docker Desktop iniciado.
- Docker Compose v2.

## Inicio local

El puerto `5432` ya está ocupado en el equipo de desarrollo. Compose publica el
PostgreSQL del proyecto en `5434`, sin interferir con esa instalación:

```powershell
docker compose up -d --wait
go run .
```

La API escucha en `http://localhost:8080`. `.env` contiene la configuración
local de desarrollo y está ignorado por Git. `.env.example` solo enumera los
nombres de las variables, sin valores; si se necesita otra configuración, copie
esa plantilla a `.env` y complete allí los valores localmente. `DB_PORT` cambia
el puerto publicado por Docker y el puerto al que Go se conecta. Si `5434`
también está ocupado, elija otro puerto libre y cambie `DB_PORT` en `.env`.

Al arrancar, la aplicación conecta a PostgreSQL, crea/actualiza las tablas con
GORM y agrega los datos de demostración de forma idempotente. Si la conexión o
la migración fallan, `go run .` termina con un error visible.

## Comprobar PostgreSQL y las tablas

En otra terminal:

```powershell
docker compose ps
docker compose exec postgres pg_isready -U notix -d notix
docker compose exec postgres psql -U notix -d notix -c "\dt"
```

La lista debe incluir `usuarios`, `cuentas`, `tableros`, `columnas`, `tareas`,
`tableros_compartidos` y `comentarios`. Para verificar las dos cuentas sin
consultar hashes:

```powershell
docker compose exec postgres psql -U notix -d notix -c "SELECT correo, rol FROM cuentas ORDER BY id;"
```

Al iniciar, el servidor imprime los IDs del tablero, la columna inicial y la
tarea reservada para la transición prohibida. La instalación nueva usa `1` para
cada uno; si se reutiliza un volumen, use los IDs que imprime el servidor en la
colección Bruno.

## Cuentas de demostración

| Rol | Correo | Clave |
| --- | --- | --- |
| Propietario | `propietaria@notix.test` | `Prop-demo-2026` |
| Colaborador | `colaborador@notix.test` | `Colab-demo-2026` |

Son credenciales inventadas solo para la demostración. El sistema almacena
claves con bcrypt y no implementa registro ni CRUD de cuentas.

## API implementada

- `POST /login`: entrega un JWT de corta duración.
- `GET /tableros/{tableroID}/tareas`: lista tareas si el usuario tiene acceso.
- `GET /tareas/{id}`: consulta una tarea con autorización por tablero.
- `POST /tareas`: crea una tarea en estado `pendiente`.
- `PATCH /tareas/{id}/estado`: aplica la máquina de estados; solo el propietario
  puede cambiar el estado.
- `GET /health`: comprueba que el servidor responde.

Transiciones permitidas: `pendiente → en_proceso`,
`en_proceso → pendiente/completado` y `completado → en_proceso`.
`pendiente → completado` devuelve `422`. Un colaborador que intenta cambiar el
estado recibe `403`.

Las operaciones restantes están especificadas como `(plan)` en
[`docs/openapi.yaml`](docs/openapi.yaml); no se exponen todavía como rutas.

En Bruno, importa la carpeta `docs/hito2_bruno/`, selecciona el entorno
`Local` y ejecuta primero la petición 0; después envía las peticiones 1 a 6 en
orden. La preparación crea la tarea que comparten las peticiones 4 y 5 y genera
tokens recientes para los dos roles.

## Pruebas y contrato

```powershell
go test ./...
npx @redocly/cli lint docs/openapi.yaml --format=stylish
```

La prueba de integración ejercita login, JWT, permisos, creación/listado y
transiciones contra PostgreSQL. Si PostgreSQL no está disponible, esa prueba se
marca como omitida; inicie Compose para ejecutarla realmente.

## Organización

La estructura por capas y las decisiones están explicadas en
[`docs/decisiones.md`](docs/decisiones.md). La hoja y la secuencia de
demostración están en [`docs/hito2_hoja.md`](docs/hito2_hoja.md) y
[`docs/hito2_bruno/`](docs/hito2_bruno/).
