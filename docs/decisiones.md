# Decisiones de arquitectura · Notix

## Arquitectura observada

El backend usa **capas clásicas** con separación por responsabilidad, organizadas bajo el paquete interno de Go. No es una implementación completa de arquitectura limpia ni de puertos y adaptadores. En la fase 1 del taller, `TareaService` pasó a depender de las interfaces `TareaRepository` y `BoardAccessRepository`, declaradas en `internal/services`; el router recibe la primera como dependencia. Esta frontera permite sustituir el adaptador de tareas, aunque las demás dependencias y el modelo siguen acoplados a GORM y a los modelos persistentes.

## Estructura relevante

| Ruta | Responsabilidad observada |
| --- | --- |
| `main.go` | Punto de entrada: carga configuración, conecta y migra/siembra la base de datos, registra los IDs de demostración y levanta el servidor HTTP. |
| `config/` | Carga variables de entorno, aplica valores por defecto y valida puerto, TTL del token y secreto JWT. |
| `database/` | Conexión PostgreSQL con GORM, migración de modelos y creación idempotente de cuentas/datos demo. |
| `internal/models/` | Structs persistentes, relaciones GORM, nombres de tabla y constantes de roles, invitaciones y estados. |
| `internal/repository/` | Consultas y escrituras de persistencia: tareas, cuentas y autorización/acceso a tablero. Traduce ausencia/conflicto en errores de dominio del repositorio. |
| `internal/repository/filas/` | Código generado por sqlc v1.31.1 para las consultas declaradas en `db/consultas.sql`. |
| `db/` y `sqlc.yaml` | Esquema acotado de `tareas`, consultas sqlc y configuración de generación de código. |
| `internal/services/` | Reglas de negocio: autenticación y JWT, autorización de pertenencia al tablero, validación de tareas y transiciones. |
| `internal/handlers/` | Adaptadores HTTP: leen path/body, llaman al servicio y serializan respuestas JSON y errores con códigos HTTP. |
| `internal/middleware/` | Exige el esquema Bearer, valida el JWT y pasa el principal autenticado en el contexto de la petición. |
| `internal/httpapi/` | Compone repositorios, servicios y handlers en el `ServeMux`; contiene además la prueba de integración del router. |
| `docs/` | OpenAPI, colección Bruno de demostración y documentación del proyecto/hitos. `docs/openapi.yaml` incluye explícitamente operaciones planeadas además de las disponibles. |
| `go.mod` | Módulo y dependencias Go: JWT, godotenv, bcrypt vía `x/crypto`, GORM y driver PostgreSQL. |
| `docker-compose.yml` | Servicio local de PostgreSQL para desarrollo, según el README. |
| `.env.example` / `.env` | Plantilla de nombres de configuración y valores locales (el archivo `.env` se indica como ignorado por Git). |

La estructura de carpetas se describe a partir de los archivos consultados; no se atribuyen a la API handlers para comentarios, tableros o columnas porque no están conectados al router.

## Flujo de una petición

1. `main.go` carga configuración y prepara PostgreSQL. GORM migra `Usuario`, `Tablero`, `Columna`, `Tarea`, `TableroCompartido`, `Comentario` y `Cuenta`; luego obtiene el `*sql.DB` del pool y construye el repositorio compuesto de tareas.
2. `internal/httpapi/router.go` recibe el repositorio de tareas por interfaz, crea los repositorios de acceso/cuentas, servicios y handlers, y registra las rutas.
3. En rutas protegidas, `middleware.RequireAuth` lee `Authorization: Bearer ...`, valida JWT y adjunta el identificador/rol al contexto. Sin token válido responde `401`.
4. El handler valida/decodifica datos HTTP y llama al servicio correspondiente.
5. El servicio aplica reglas de acceso y negocio mediante la interfaz del repositorio. El adaptador compuesto `TareaSqlcRepository` usa sqlc para `FindByID` y `UpdateState`; delega `ListByBoard`, `Create` y `FindColumn` al repositorio GORM embebido. La consulta por ID complementa la tarea con la columna usando ese repositorio, necesaria para comprobar el tablero.
6. El handler convierte el resultado o error conocido en JSON y código HTTP; el router devuelve la respuesta al cliente.

En autenticación, `POST /login` sigue el camino handler → servicio de autenticación → repositorio de cuentas; el servicio compara bcrypt y firma un JWT HS256. En las demás rutas protegidas, el middleware autentica primero y el servicio autoriza el acceso al tablero/tarea.

## Decisiones respaldadas por el código

- **Autenticación:** JWT firmado HS256; claims de usuario, rol, emisor, sujeto y vencimiento. El TTL por defecto es 8 horas y puede configurarse. Las contraseñas demo se generan y comparan con bcrypt. En desarrollo se permite un secreto por defecto; fuera de desarrollo se requiere `JWT_SECRET` de al menos 32 caracteres.
- **Autorización:** la pertenencia al tablero se comprueba en servicios/repositorios para consultar o crear tareas. Solo el rol `propietario` puede cambiar el estado; un colaborador recibe `403`.
- **Validación:** los handlers exigen JSON válido, rechazan claves desconocidas y limitan el cuerpo a 64 KiB. Los servicios validan título, longitud del contenido, columna y transiciones. El modelo de tarea incluye una restricción de estados permitidos en base de datos.
- **Persistencia:** PostgreSQL con migración automática de GORM y relaciones/índices declarados en structs. En la configuración ejecutada por `main.go`, sqlc v1.31.1 genera la búsqueda por ID y la actualización condicional de estado; el adaptador GORM conserva las demás operaciones de tarea. La actualización condiciona la escritura al estado anterior para detectar conflictos concurrentes.
- **Errores:** handlers centralizan conversión de errores conocidos a JSON y estados HTTP (400, 401, 403, 404, 409 y 422). Errores no clasificados se registran en servidor y se exponen como 500 genérico.

## Motivos y efectos para mantenimiento

La separación permite cambiar el formato HTTP en handlers sin mezclar reglas de negocio, concentrar reglas de tareas/autorización en servicios y mantener consultas en repositorios. Los DTOs de tarea mantienen el contrato JSON separado de los structs persistentes. La interfaz del repositorio de tareas permite alternar su implementación sin cambiar las reglas del servicio. La composición explícita del router deja visible el cableado de dependencias. La prueba de integración construye el router contra una transacción de PostgreSQL.

La contrapartida observada es que los servicios aún usan modelos persistentes y no toda dependencia está abstraída; por ello no se afirma que el dominio esté aislado como en arquitectura limpia. El adaptador sqlc también mezcla una consulta generada con la carga GORM de la columna requerida para autorización.

## Semana 10 · transacciones

### D1 · Caso que requiere atomicidad

Las operaciones actuales de `TareaService` no escriben en dos tablas dentro de un mismo caso de uso: `Create` inserta una fila de `tareas` y `ChangeState` actualiza esa fila. Por tanto, el código real no necesita hoy una transacción que coordine dos escrituras. El fragmento de referencia `ConfirmarYCobrar` sí muestra el problema hipotético: primero guarda el cambio de estado de una reserva y luego crea un pago; si falla el segundo paso, queda una reserva confirmada sin pago. En el dominio Notix, un ejemplo futuro sería crear una tarea y también su primer comentario, pero el API actual no implementa esa operación conjunta.

### D2 · Alternativas consideradas

1. Ejecutar ambas llamadas al repositorio por separado. Es simple, pero si la segunda escritura falla deja persistido solo el primer cambio; incumple la atomicidad del caso hipotético.
2. Iniciar una transacción en el handler o en `TareaService`, pasando directamente `*gorm.DB` o `*sql.Tx`. Puede agrupar las escrituras, pero acopla la capa de negocio a GORM o `database/sql` y dificulta sustituir adaptadores y probar el caso con un repositorio falso.
3. Definir una unidad de trabajo/interfaz transaccional en la capa de negocio e implementarla en cada adaptador. La operación completa se ejecutaría mediante repositorios vinculados a una misma transacción; el servicio expresaría el caso de uso sin conocer el motor. Es la opción conceptual preferible si Notix incorpora una operación con varias escrituras.

### D3 · Diseño previsto si aparece ese caso

Si se implementara el caso futuro de crear tarea y comentario de forma atómica, el servicio recibiría una abstracción de transacción de negocio (por ejemplo, `WithinTransaction(func(TareaRepository, ComentarioRepository) error) error`). El adaptador GORM la implementaría con `db.Transaction`; el adaptador sqlc abriría `BeginTx`, construiría `filas.New(tx)` mediante `Queries.WithTx(tx)`, confirmaría al terminar sin error y revertiría ante error. Ambos repositorios usados dentro de la operación tendrían que compartir la misma transacción. El servicio mantendría las reglas y el orden del caso; handler, servicio y entidades no importarían GORM ni `database/sql`.

Este diseño es una decisión razonada para un requisito futuro, no código implementado en el taller: las operaciones reales no requieren coordinar dos escrituras. La generación sqlc actual incluye únicamente búsqueda por ID y actualización de estado.

### D4 · Alternativa descartada

Se descarta abrir/confirmar/revertir la transacción directamente desde el handler o hacer que el servicio importe `gorm.DB` o `sql.Tx`. Eso trasladaría persistencia al transporte o filtraría el motor a la capa de negocio, acoplaría el caso de uso a un adaptador y haría más difícil probarlo sin PostgreSQL.

### D5 · Condición para sqlc

En una futura operación transaccional, todas las consultas sqlc del caso tendrían que ejecutarse con el mismo `*sql.Tx`, mediante un único conjunto de `Queries` ligado con `WithTx`. No se debe combinar una escritura sqlc fuera de la transacción con otra GORM dentro de ella (ni crear una transacción independiente por adaptador), porque ante un fallo volverían a existir resultados parciales.

### D6 · Alcance de esta entrega

No se añadieron tablas, entidades ni código transaccional: el caso actual no contiene una operación de dos escrituras que lo justifique. Esta sección documenta el análisis y la opción para un requisito futuro; no afirma que `ConfirmarYCobrar`, pagos, comentarios automáticos ni una unidad de trabajo existan en Notix.
