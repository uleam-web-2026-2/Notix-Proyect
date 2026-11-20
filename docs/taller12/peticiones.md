# Semana 12 · las seis peticiones del hito 2

| Petición | Antes (taller 11) | Con las defensas | Con el usuario de la API (reto) |
| --- | --- | --- | --- |
| 1 sin token | 401 | 401 | |
| 2 login | 200 | 200 | |
| 3 con token | 200 | 200 | |
| 4 rol sin permiso | 403 | 403 | |
| 5 cambio permitido | 200 | 200 | |
| 6 transición prohibida | 422 | 422 | |

Los seis códigos de la columna «Con las defensas» se observaron en la API de
esta rama mediante solicitudes HTTP locales el 19 de noviembre de 2026. La
preparación creó la tarea `44`; la petición 4 recibió `403` y no la cambió, la
petición 5 recibió `200` y la transición prohibida de la tarea `1` recibió
`422`. Los eventos reales del servidor están en `logs_seis.txt`.
