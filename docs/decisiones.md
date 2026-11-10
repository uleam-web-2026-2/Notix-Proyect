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

## Semana 11 · el contenedor

### D1 · Qué entra en la imagen

La imagen final `notix:latest` ocupa 33,4 MB en disco (9,86 MB de contenido) e incluye el binario estático de 17,8 MB y los archivos mínimos de Distroless. Las herramientas de Go, las dependencias descargadas y el código fuente quedan en la etapa `compilacion`, que ocupa 1,62 GB en disco (386 MB de contenido). Sin `.dockerignore`, `.env` se enviaría en el contexto de construcción y quedaría en `/src` dentro de la etapa de compilación; no se copiaría a la imagen final, porque la etapa final solo copia `/servidor`. Excluirlo evita exponer el secreto en el contexto y en las capas/cache de compilación.

### D2 · Dónde está la base, vista desde el contenedor

Dentro de un contenedor, `localhost` apunta al propio contenedor, donde no corre PostgreSQL. La base del proyecto está publicada en el puerto `5434` de la máquina anfitriona; por eso el contenedor usa `host.docker.internal` como `DB_HOST` y conserva `DB_PORT=5434`. En un servidor, se configuraría el nombre DNS o dirección de la base accesible desde ese entorno, junto con sus credenciales y reglas de red; no se dependería de `host.docker.internal`.

### D3 · Los puertos

Se ejecutó `-p 8080:8080`: el primer `8080` es el puerto de la máquina al que se conecta Bruno y el segundo es el puerto dentro del contenedor. El segundo debe coincidir con `PORT=8080`, donde escucha la API. Si no coincide, Docker dirige el tráfico a un puerto sin servidor y la API no será accesible por ese mapeo.

### D4 · La hora

Las fechas de la respuesta 3 conservan el mismo instante: por ejemplo, `2026-10-24T11:54:47.280102-05:00` con `go run` y `2026-10-24T16:54:47.280102Z` en el contenedor. Distroless entrega UTC y la API serializa el instante con su zona/offset, por lo que cambió la representación, no la hora real. Se mantiene la fecha en UTC en el contenedor y los clientes deben interpretar el offset para mostrarla en su zona local.

### Registro de peticiones HTTP

El middleware de logging escribe en la salida estándar el método, la ruta y el código HTTP de cada petición, de modo que `docker logs` permite verificar las respuestas atendidas por el contenedor. No registra parámetros de consulta ni tokens.
