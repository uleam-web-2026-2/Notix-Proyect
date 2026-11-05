# Línea de base del Hito 2

Peticiones ejecutadas contra la API local antes de las fases del taller. Se ejecutó primero la petición `0 Preparar demo` para obtener tokens recientes y preparar los IDs; después se enviaron las peticiones 1 a 6 de Bruno en orden.

| Petición | Antes | Después |
| --- | --- | --- |
| 1 sin token | 401 | |
| 2 login | 200 | |
| 3 con token | 200 | |
| 4 rol sin permiso | 403 | |
| 5 cambio permitido | 200 | |
| 6 transición prohibida | 422 | |

## Resultado de la línea de base

Los códigos observados fueron 401, 200, 200, 403, 200 y 422, respectivamente. Los IDs de demostración los imprimió el servidor al iniciar; no se fijaron IDs en esta hoja.
