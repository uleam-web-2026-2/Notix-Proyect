# Semana 13 · las seis peticiones del hito 2

| Petición | Antes (taller 12) | Con el WebSocket | Con el reto |
| --- | --- | --- | --- |
| 1 sin token | `401` |  |  |
| 2 login | `200` (token devuelto) |  |  |
| 3 con token | `200` |  |  |
| 4 rol sin permiso | `403` |  |  |
| 5 cambio permitido | `200` |  |  |
| 6 transición prohibida | `422` |  |  |

La columna «Con el reto» se llena solo si hicieron el reto.
Antes de esta vuelta, el servidor se arrancó otra vez: el límite del taller 12
cuenta los inicios de sesión durante un minuto.

Los códigos de «Antes (taller 12)» se observaron ejecutando localmente la rama
`taller12` (`9a1620f`) y enviando peticiones HTTP directas. La colección de
Bruno no se ejecutó en este entorno; por eso estos resultados no se presentan
como una ejecución de Bruno.
