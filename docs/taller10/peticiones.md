# Línea de base del Hito 2

Peticiones ejecutadas contra la API local antes de las fases del taller. Se ejecutó primero la petición `0 Preparar demo` para obtener tokens recientes y preparar los IDs; después se enviaron las peticiones 1 a 6 de Bruno en orden.

| Petición | Antes | Después |
| --- | --- | --- |
| 1 sin token | 401 | 401 |
| 2 login | 200 | 200 |
| 3 con token | 200 | 200 |
| 4 rol sin permiso | 403 | 403 |
| 5 cambio permitido | 200 | 200 |
| 6 transición prohibida | 422 | 422 |

## Resultado de la línea de base

Los códigos observados fueron 401, 200, 200, 403, 200 y 422, respectivamente. Los IDs de demostración los imprimió el servidor al iniciar; no se fijaron IDs en esta hoja.

## Resultado posterior a las fases 1–3

El 5 de noviembre de 2026 se repitieron las seis peticiones contra la API local con el adaptador sqlc seleccionado en `main.go`. Se usaron las cuentas de demostración configuradas por `database.SeedDemoData`: `colaborador@notix.test` y `propietaria@notix.test`; sus contraseñas no se incluyen aquí. Para las peticiones 4 y 5 se creó la tarea de demostración 27; la petición 6 usó la tarea prohibida 1, que permanecía en `pendiente`.

| # | Petición ejecutada | Credencial | Resultado observado |
| --- | --- | --- | --- |
| 1 | `GET /tableros/1/tareas` | Ninguna | `401` |
| 2 | `POST /login` | Correo y clave del colaborador | `200`, respuesta incluyó `token` Bearer |
| 3 | `GET /tableros/1/tareas` | Token del colaborador | `200` |
| 4 | `PATCH /tareas/27/estado` con `{"estado":"en_proceso"}` | Token del colaborador | `403` |
| 5 | `PATCH /tareas/27/estado` con `{"estado":"en_proceso"}` | Token de propietaria | `200` |
| 6 | `PATCH /tareas/1/estado` con `{"estado":"completado"}` | Token de propietaria | `422` |

## Prueba del paso 2.1: campos controlados por el servidor

El 5 de noviembre de 2026 se envió a la API local `POST /tareas` con `id` y `estado` incluidos en el cuerpo. La petición fue rechazada con `400` y el error observado fue `invalid_json` (`El cuerpo JSON no es válido.`). No se documenta aquí un resultado esperado como si fuera real.

Se utilizó el token del colaborador y este cuerpo:

```json
{
  "id": 999,
  "estado": "completado",
  "titulo": "Intento campos controlados por servidor",
  "contenido": "Debe rechazarse antes de crear la tarea.",
  "columna_id": 1
}
```

```json
{
  "error": {
    "code": "invalid_json",
    "message": "El cuerpo JSON no es válido."
  }
}
```

## Verificación después de la transición prohibida

Después de observar el `422` de la petición 6, se consultó `GET /tareas/1` con el token de propietaria. La respuesta observada fue `200` y el cuerpo conservó `"estado": "pendiente"`.

## Ejecución de la colección Bruno CLI

El 5 de noviembre de 2026 se ejecutó la colección local con Bruno CLI `4.3.0` mediante `npx --yes @usebruno/cli@4.3.0 run --env Local`, desde `docs/hito2_bruno`. El resultado fue **7 de 7 peticiones y 7 de 7 pruebas aprobadas**:

| # | Petición de la colección | Código observado | Prueba |
| --- | --- | --- | --- |
| 0 | Preparar demo | `201` | Aprobada |
| 1 | Sin token | `401` | Aprobada |
| 2 | Login colaborador | `200` | Aprobada |
| 3 | Listar tareas | `200` | Aprobada |
| 4 | Cambio colaborador | `403` | Aprobada |
| 5 | Cambio propietario | `200` | Aprobada |
| 6 | Transición prohibida | `422` | Aprobada |

Bruno actualizó temporalmente valores del entorno local durante la ejecución; `docs/hito2_bruno/environments/Local.bru` se restauró a su contenido original y no forma parte de los cambios.
