# Hito 2 · Hoja de verificación de Notix

> Esta hoja describe el backend que está en el repositorio. La imagen de referencia usa una API de tickets; Notix implementa tableros y tareas, por lo que aquí se documentan sus rutas y roles reales. Los resultados esperados se distinguen de los resultados realmente ejecutados.

## 1. Usuarios de ejemplo

Las cuentas se crean idempotentemente al iniciar la aplicación, mediante `database.SeedDemoData`. El login compara hashes bcrypt guardados en la tabla `cuentas`.

<table>
  <thead>
    <tr><th scope="col">Usuario (correo)</th><th scope="col">Rol</th><th scope="col">Contraseña de demostración</th><th scope="col">Evidencia</th></tr>
  </thead>
  <tbody>
      <tr><td><code>propietaria@notix.test</code></td><td>Propietario (<code>propietario</code>)</td><td><code>Prop-demo-2026</code></td><td><code>database/database.go</code>, <code>README.md</code></td></tr>
      <tr><td><code>colaborador@notix.test</code></td><td>Colaborador (<code>colaborador</code>)</td><td><code>Colab-demo-2026</code></td><td><code>database/database.go</code>, <code>README.md</code></td></tr>
  </tbody>
</table>

Son credenciales ficticias de desarrollo/demostración incluidas en el código; no son claves personales. La contraseña se almacena como hash bcrypt. No se implementa registro de cuentas por API.

## 2. Las seis peticiones

La colección Bruno está en `docs/hito2_bruno/`. Primero se ejecuta `0 Preparar demo.bru`, que obtiene tokens recientes y prepara variables/IDs para la secuencia. Las rutas de las filas son las implementadas. Los resultados esperados se derivan de las pruebas declaradas en la colección y el código. **No ejecuté la colección en este entorno**, por lo que el resultado real queda pendiente.

<table>
  <thead>
    <tr><th scope="col">#</th><th scope="col">Caso y petición</th><th scope="col">Credencial necesaria</th><th scope="col">Resultado esperado</th><th scope="col">Resultado real</th></tr>
  </thead>
  <tbody>
      <tr><td>1</td><td><code>GET /tableros/{boardId}/tareas</code> sin token (<code>1 Sin token.bru</code>)</td><td>Ninguna</td><td><code>401</code>, error <code>unauthorized</code></td><td>Pendiente de ejecutar</td></tr>
      <tr><td>2</td><td><code>POST /login</code> con el usuario colaborador (<code>2 Login colaborador.bru</code>)</td><td>Correo y contraseña del colaborador, en el JSON</td><td><code>200</code>; devuelve token Bearer y rol <code>colaborador</code></td><td>Pendiente de ejecutar</td></tr>
      <tr><td>3</td><td><code>GET /tableros/{boardId}/tareas</code> con token de colaborador (<code>3 Listar tareas.bru</code>)</td><td>Bearer de colaborador con invitación activa al tablero demo</td><td><code>200</code>; devuelve tareas del tablero</td><td>Pendiente de ejecutar</td></tr>
      <tr><td>4</td><td><code>PATCH /tareas/{taskId}/estado</code> solicitando <code>en_proceso</code> (<code>4 Cambio colaborador 403.bru</code>)</td><td>Bearer de colaborador</td><td><code>403</code>; solo propietario puede cambiar estado</td><td>Pendiente de ejecutar</td></tr>
      <tr><td>5</td><td><code>PATCH /tareas/{taskId}/estado</code> de <code>pendiente</code> a <code>en_proceso</code> (<code>5 Cambio propietario 200.bru</code>)</td><td>Bearer de propietario</td><td><code>200</code>; devuelve la tarea con estado actualizado</td><td>Pendiente de ejecutar</td></tr>
      <tr><td>6</td><td><code>PATCH /tareas/{forbiddenTaskId}/estado</code> de <code>pendiente</code> a <code>completado</code> (<code>6 Transicion prohibida 422.bru</code>)</td><td>Bearer de propietario</td><td><code>422</code>, error <code>invalid_transition</code>; la tarea reservada permanece en <code>pendiente</code></td><td>Pendiente de ejecutar</td></tr>
  </tbody>
</table>

Los valores `boardId`, `taskId` y `forbiddenTaskId` se resuelven con los datos de demostración. No se debe asumir que siempre sean `1`: al reutilizar una base de datos los IDs pueden cambiar y el servidor los imprime al iniciar.

## 3. Matriz ruta × rol

`Propietario` y `colaborador` se refieren a cuentas autenticadas. En rutas de tareas, los resultados de acceso autorizado presuponen que la cuenta pertenece al tablero. El middleware devuelve `401` si falta un token válido. Para usuarios autenticados sin acceso al tablero, las consultas y la creación devuelven `403`.

<table>
  <thead>
    <tr><th scope="col">Ruta realmente registrada</th><th scope="col">Sin sesión</th><th scope="col">Propietario con acceso</th><th scope="col">Colaborador con acceso</th></tr>
  </thead>
  <tbody>
      <tr><td><code>GET /health</code></td><td><code>200</code></td><td><code>200</code></td><td><code>200</code></td></tr>
      <tr><td><code>POST /login</code></td><td><code>200</code> con credenciales válidas; <code>401</code> si no coinciden</td><td>Igual; requiere credenciales en el cuerpo</td><td>Igual; requiere credenciales en el cuerpo</td></tr>
      <tr><td><code>GET /tableros/{tableroID}/tareas</code></td><td><code>401</code></td><td><code>200</code>; <code>403</code> si no tiene acceso</td><td><code>200</code>; <code>403</code> si no tiene acceso</td></tr>
      <tr><td><code>GET /tareas/{id}</code></td><td><code>401</code></td><td><code>200</code>; <code>403</code> si no tiene acceso</td><td><code>200</code>; <code>403</code> si no tiene acceso</td></tr>
      <tr><td><code>POST /tareas</code></td><td><code>401</code></td><td><code>201</code>; <code>403</code> si no tiene acceso a la columna/tablero</td><td><code>201</code>; <code>403</code> si no tiene acceso a la columna/tablero</td></tr>
      <tr><td><code>PATCH /tareas/{id}/estado</code></td><td><code>401</code></td><td><code>200</code> si la transición es válida</td><td><code>403</code></td></tr>
  </tbody>
</table>

IDs inválidos producen `400`; recursos inexistentes producen `404`; validaciones de tarea o transiciones no permitidas producen `422`. No hay rutas `PUT` ni `DELETE` registradas para estos recursos. La fuente de esta matriz es `internal/httpapi/router.go`; las operaciones marcadas como plan en `docs/openapi.yaml` no forman parte de la matriz.

## 4. Recuadro de cumplimiento

Criterios marcados como cumplidos solo cuando se verificaron en archivos del repositorio. La existencia de pruebas o aserciones no acredita una ejecución exitosa.

- [x] Los structs de las entidades principales están definidos en `internal/models/models.go`; sus nombres de tabla están definidos y GORM ejecuta `AutoMigrate` en `database/database.go`.
- [x] Hay validación en el servicio: título requerido con máximo 150 caracteres, contenido de hasta 5000, identificador positivo de columna y transiciones de estado permitidas. El decodificador rechaza JSON mal formado, claves desconocidas y cuerpos con más de un objeto.
- [x] Existe una colección Bruno con los seis casos de esta hoja en `docs/hito2_bruno/`; los archivos declaran sus códigos esperados.
- [ ] Resultado real de las seis peticiones Bruno: pendiente de ejecutar y registrar.
- [x] Existe una prueba de integración de la API en `internal/httpapi/router_integration_test.go`, que cubre autenticación, permisos, creación/consulta y transiciones con PostgreSQL.
- [x] Se ejecutó `go test ./...`: la captura compartida muestra salida `ok` para `notix/internal/httpapi` y `notix/internal/services`; los otros paquetes indican `[no test files]`.
- [ ] Ejecución confirmada de la prueba de integración contra PostgreSQL: la captura no usa `-v`, por lo que no permite saber si la prueba se ejecutó o se omitió por falta de base de datos.
- [x] El código separa modelos, handlers, servicios, repositorios y acceso HTTP; esta organización se explica en `docs/decisiones.md`.
- [ ] Capturas de respaldo de cada paso: no localizadas entre los archivos revisados; pendientes de agregar/verificar.

## Fuentes revisadas

`internal/httpapi/router.go`, `internal/handlers/`, `internal/middleware/auth.go`, `internal/services/`, `internal/repository/`, `internal/models/models.go`, `database/database.go`, `internal/httpapi/router_integration_test.go`, `README.md` y `docs/hito2_bruno/`.
