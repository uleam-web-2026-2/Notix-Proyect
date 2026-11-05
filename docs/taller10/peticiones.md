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

La petición de prueba con `id` y `estado` en `POST /tareas` queda pendiente de ejecutar. No se autorizó el envío de esa petición a la API local. La inspección del código anterior a los DTO muestra que `decodeJSON` rechazaba campos desconocidos, por lo que `400` era el resultado previsto por el código; no se registra como resultado observado.
