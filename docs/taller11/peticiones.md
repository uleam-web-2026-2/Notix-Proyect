# Semana 11 · las seis peticiones del hito 2

| Petición | Con go run | En el contenedor | Con Compose |
| --- | --- | --- | --- |
| 1 sin token | 401 | | |
| 2 login | 200 | | |
| 3 con token | 200 | | |
| 4 rol sin permiso | 403 | | |
| 5 cambio permitido | 200 | | |
| 6 transición prohibida | 422 | | |

Petición 3 con go run (respuesta completa):

```json
[{"id":1,"titulo":"DEMO: transición prohibida","contenido":"Debe permanecer en pendiente durante toda la demostración.","estado":"pendiente","columna_id":1,"columna":{"id":1,"nombre":"Por hacer","orden":1,"tablero_id":1},"creador_id":1,"creador":{"id":1,"nombre":"Propietaria Demo","correo":"propietaria@notix.test","creado":"2026-10-24T11:54:47.169774-05:00"},"creado":"2026-10-24T11:54:47.280102-05:00"},{"id":16,"titulo":"Tarea nueva para la demostración","contenido":"Este registro se usará una sola vez en las peticiones 4 y 5.","estado":"en_proceso","columna_id":1,"columna":{"id":1,"nombre":"Por hacer","orden":1,"tablero_id":1},"creador_id":2,"creador":{"id":2,"nombre":"Colaborador Demo","correo":"colaborador@notix.test","creado":"2026-10-24T11:54:47.251984-05:00"},"creado":"2026-10-24T13:43:53.878281-05:00"},{"id":23,"titulo":"Tarea nueva para la demostración","contenido":"Este registro se usará una sola vez en las peticiones 4 y 5.","estado":"en_proceso","columna_id":1,"columna":{"id":1,"nombre":"Por hacer","orden":1,"tablero_id":1},"creador_id":2,"creador":{"id":2,"nombre":"Colaborador Demo","correo":"colaborador@notix.test","creado":"2026-10-24T11:54:47.251984-05:00"},"creado":"2026-10-24T14:32:32.82536-05:00"},{"id":24,"titulo":"Tarea nueva para la demostración","contenido":"Este registro se usará una sola vez en las peticiones 4 y 5.","estado":"en_proceso","columna_id":1,"columna":{"id":1,"nombre":"Por hacer","orden":1,"tablero_id":1},"creador_id":2,"creador":{"id":2,"nombre":"Colaborador Demo","correo":"colaborador@notix.test","creado":"2026-10-24T11:54:47.251984-05:00"},"creado":"2026-10-28T15:39:51.777858-05:00"},{"id":25,"titulo":"Tarea nueva para la demostración","contenido":"Este registro se usará una sola vez en las peticiones 4 y 5.","estado":"en_proceso","columna_id":1,"columna":{"id":1,"nombre":"Por hacer","orden":1,"tablero_id":1},"creador_id":2,"creador":{"id":2,"nombre":"Colaborador Demo","correo":"colaborador@notix.test","creado":"2026-10-24T11:54:47.251984-05:00"},"creado":"2026-10-28T15:51:51.215489-05:00"},{"id":26,"titulo":"Tarea nueva para la demostración","contenido":"Este registro se usará una sola vez en las peticiones 4 y 5.","estado":"pendiente","columna_id":1,"columna":{"id":1,"nombre":"Por hacer","orden":1,"tablero_id":1},"creador_id":2,"creador":{"id":2,"nombre":"Colaborador Demo","correo":"colaborador@notix.test","creado":"2026-10-24T11:54:47.251984-05:00"},"creado":"2026-10-29T09:22:31.148008-05:00"},{"id":27,"titulo":"Tarea nueva para la demostración","contenido":"Este registro se usará una sola vez en las peticiones 4 y 5.","estado":"en_proceso","columna_id":1,"columna":{"id":1,"nombre":"Por hacer","orden":1,"tablero_id":1},"creador_id":2,"creador":{"id":2,"nombre":"Colaborador Demo","correo":"colaborador@notix.test","creado":"2026-10-24T11:54:47.251984-05:00"},"creado":"2026-10-29T09:22:50.312281-05:00"},{"id":34,"titulo":"Tarea nueva para la demostración","contenido":"Este registro se usará una sola vez en las peticiones 4 y 5.","estado":"pendiente","columna_id":1,"columna":{"id":1,"nombre":"Por hacer","orden":1,"tablero_id":1},"creador_id":2,"creador":{"id":2,"nombre":"Colaborador Demo","correo":"colaborador@notix.test","creado":"2026-10-24T11:54:47.251984-05:00"},"creado":"2026-11-10T11:04:23.088146-05:00"}]
```

La misma respuesta se registrará en las columnas «En el contenedor» y «Con
Compose» al completar la fase 2.
