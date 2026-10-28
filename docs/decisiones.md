# Decisiones de arquitectura · Notix

## Arquitectura observada

El backend usa **capas clásicas** con separación por responsabilidad, organizadas bajo el paquete interno de Go. No es una implementación de arquitectura limpia completa ni de puertos y adaptadores: el router construye concretamente repositorios y servicios, y los servicios dependen en general de tipos concretos de repositorio. Hay una interfaz pequeña (`AccountStore`) para la consulta de cuentas, pero eso por sí solo no convierte el conjunto en arquitectura hexagonal.

## Estructura relevante

| Ruta | Responsabilidad observada |
| --- | --- |
| `main.go` | Punto de entrada: carga configuración, conecta y migra/siembra la base de datos, registra los IDs de demostración y levanta el servidor HTTP. |
| `config/` | Carga variables de entorno, aplica valores por defecto y valida puerto, TTL del token y secreto JWT. |
| `database/` | Conexión PostgreSQL con GORM, migración de modelos y creación idempotente de cuentas/datos demo. |
| `internal/models/` | Structs persistentes, relaciones GORM, nombres de tabla y constantes de roles, invitaciones y estados. |
| `internal/repository/` | Consultas y escrituras de persistencia: tareas, cuentas y autorización/acceso a tablero. Traduce ausencia/conflicto en errores de dominio del repositorio. |
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

1. `main.go` carga configuración y prepara PostgreSQL. GORM migra `Usuario`, `Tablero`, `Columna`, `Tarea`, `TableroCompartido`, `Comentario` y `Cuenta`.
2. `internal/httpapi/router.go` crea repositorios, servicios y handlers, y registra las rutas.
3. En rutas protegidas, `middleware.RequireAuth` lee `Authorization: Bearer ...`, valida JWT y adjunta el identificador/rol al contexto. Sin token válido responde `401`.
4. El handler valida/decodifica datos HTTP y llama al servicio correspondiente.
5. El servicio aplica reglas de acceso y negocio y utiliza repositorio para leer o escribir con GORM en PostgreSQL.
6. El handler convierte el resultado o error conocido en JSON y código HTTP; el router devuelve la respuesta al cliente.

En autenticación, `POST /login` sigue el camino handler → servicio de autenticación → repositorio de cuentas; el servicio compara bcrypt y firma un JWT HS256. En las demás rutas protegidas, el middleware autentica primero y el servicio autoriza el acceso al tablero/tarea.

## Decisiones respaldadas por el código

- **Autenticación:** JWT firmado HS256; claims de usuario, rol, emisor, sujeto y vencimiento. El TTL por defecto es 8 horas y puede configurarse. Las contraseñas demo se generan y comparan con bcrypt. En desarrollo se permite un secreto por defecto; fuera de desarrollo se requiere `JWT_SECRET` de al menos 32 caracteres.
- **Autorización:** la pertenencia al tablero se comprueba en servicios/repositorios para consultar o crear tareas. Solo el rol `propietario` puede cambiar el estado; un colaborador recibe `403`.
- **Validación:** los handlers exigen JSON válido, rechazan claves desconocidas y limitan el cuerpo a 64 KiB. Los servicios validan título, longitud del contenido, columna y transiciones. El modelo de tarea incluye una restricción de estados permitidos en base de datos.
- **Persistencia:** PostgreSQL a través de GORM, migración automática y relaciones/índices declarados en structs. La actualización de estado condiciona la escritura al estado anterior para detectar conflictos concurrentes.
- **Errores:** handlers centralizan conversión de errores conocidos a JSON y estados HTTP (400, 401, 403, 404, 409 y 422). Errores no clasificados se registran en servidor y se exponen como 500 genérico.

## Motivos y efectos para mantenimiento

La separación permite cambiar el formato HTTP en handlers sin mezclar reglas de negocio, concentrar reglas de tareas/autorización en servicios y mantener las consultas de GORM en repositorios. Los structs compartidos describen el modelo persistente y las relaciones. La composición explícita del router deja visible el cableado de dependencias. La prueba de integración construye el router contra una transacción de PostgreSQL, mientras que los servicios se pueden revisar con menos dependencia del transporte HTTP.

La contrapartida observada es el acoplamiento de varios servicios a repositorios concretos y la dependencia directa de GORM/modelos persistentes; por ello no se afirma que el dominio esté aislado como en arquitectura limpia. La presencia de una prueba de integración no implica que se haya ejecutado con éxito en este trabajo.
