# go-api-template

![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=flat)
![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen?style=flat)

Boilerplate de producción para construir APIs REST en Go. Incluye autenticación JWT, ORM con GORM, rate limiting, documentación OpenAPI interactiva, respuestas estandarizadas, hot-reload y Docker listo para desplegar.

---

## Características

- 🔐 **Autenticación JWT** — Access + refresh tokens con rotación automática
- 🗄️ **ORM con GORM + SQLite** — Fácil de cambiar a PostgreSQL o MySQL
- 🌐 **Gin Framework** — Con middlewares de CORS, Logger, RequestID y Recovery
- 📦 **Respuestas JSON estandarizadas** — Estructura `{success, data, error, request_id}` consistente
- 🚦 **Rate limiting por IP y global** — Token-bucket con protección anti brute-force en endpoints de auth; headers `X-Ratelimit-*` y `Retry-After` en cada respuesta
- 📚 **Swagger UI interactivo** — Documentación OpenAPI auto-generada desde anotaciones en el código; disponible en `/docs/index.html`
- 📄 **Paginación genérica** — Helpers de GORM listos para usar en cualquier recurso
- ⚡ **Hot-reload con Air** — Recarga automática en desarrollo al guardar cambios
- 🐳 **Docker + Docker Compose** — Multi-stage build optimizado para producción
- 🛡️ **Graceful shutdown** — Cierre limpio del servidor ante señales del SO
- 🔧 **Makefile** — Comandos útiles para desarrollo, build, tests y más

---

## Requisitos

| Herramienta | Versión mínima | Notas |
|---|---|---|
| Go | 1.23+ | [golang.org](https://golang.org/dl/) |
| gcc | cualquiera | Requerido por GORM/SQLite (CGO) |
| libsqlite3-dev | cualquiera | Debian/Ubuntu: `apt install libsqlite3-dev` |
| Docker | 20.10+ | Opcional, para despliegue con contenedores |
| air | latest | Opcional, para hot-reload en desarrollo |

---

## Inicio rápido

```bash
# 1. Clonar el repositorio
git clone https://github.com/tu-usuario/go-api-template.git
cd go-api-template

# 2. Copiar las variables de entorno
cp .env.example .env

# 3. Editar .env con tus valores (opcional para desarrollo local)

# 4. Instalar dependencias y ejecutar
make run
```

La API estará disponible en `http://localhost:8080`.  
La documentación interactiva estará en `http://localhost:8080/docs/index.html`.

---

## Variables de entorno

Copia `.env.example` a `.env` y ajusta los valores según tu entorno.

| Variable | Descripción | Default |
|---|---|---|
| `APP_PORT` | Puerto en el que escucha la API | `8080` |
| `APP_ENV` | Entorno de ejecución (`development` \| `production`) | `development` |
| `DB_PATH` | Ruta al archivo de base de datos SQLite | `./data.db` |
| `JWT_SECRET` | Clave secreta para firmar los tokens JWT | `change_me` |
| `JWT_ACCESS_EXPIRY` | Duración del access token | `15m` |
| `JWT_REFRESH_EXPIRY` | Duración del refresh token | `168h` |
| `CORS_ORIGINS` | Orígenes permitidos (separados por coma o `*`) | `*` |
| `LOG_LEVEL` | Nivel de log (`debug` \| `info` \| `warn` \| `error`) | `info` |
| `RATE_LIMIT_ENABLED` | Activa el rate limiting global y por IP | `true` |
| `SWAGGER_ENABLED` | Sirve la UI de Swagger en `/docs/index.html` | `true` |

> ⚠️ **Importante:** Cambia `JWT_SECRET` por un valor seguro en producción. Considera poner `SWAGGER_ENABLED=false` en producción para no exponer la documentación.

---

## Estructura del proyecto

```
go-api-template/
├── cmd/
│   └── main.go                  # Punto de entrada con graceful shutdown
├── docs/                        # Spec OpenAPI generada (make swagger)
│   ├── docs.go                  # Spec embebida en el binario
│   ├── swagger.json             # Spec en JSON
│   └── swagger.yaml             # Spec en YAML
├── internal/
│   ├── config/
│   │   └── config.go            # Carga y valida variables de entorno
│   ├── database/
│   │   └── database.go          # Conexión a SQLite y auto-migración de modelos
│   ├── handlers/
│   │   ├── auth_handler.go      # Handlers de registro, login, refresh y logout
│   │   ├── health_handler.go    # Health check del servicio
│   │   ├── swagger_types.go     # Tipos nombrados para schemas de la documentación
│   │   └── user_handler.go      # Handler del perfil del usuario autenticado
│   ├── middleware/
│   │   ├── auth.go              # Valida el JWT en el header Authorization
│   │   ├── cors.go              # CORS configurable vía CORS_ORIGINS
│   │   ├── logger.go            # Log estructurado JSON de cada request (slog)
│   │   ├── ratelimit.go         # Rate limiting global y por IP (token-bucket)
│   │   └── request_id.go        # Inyecta un UUID por request (X-Request-ID)
│   ├── models/
│   │   └── user.go              # Modelos User y RefreshToken para GORM
│   ├── router/
│   │   └── router.go            # Registro de rutas, grupos y middlewares
│   ├── services/
│   │   └── auth_service.go      # Lógica de negocio de autenticación
│   └── utils/
│       ├── pagination.go        # Helper genérico de paginación con GORM
│       └── response.go          # Helpers para respuestas JSON estandarizadas
├── bin/                         # Binario compilado (generado, en .gitignore)
├── .air.toml                    # Configuración de hot-reload para Air
├── .env                         # Variables de entorno locales (en .gitignore)
├── .env.example                 # Plantilla de variables de entorno (commiteado)
├── .gitignore
├── Dockerfile                   # Multi-stage build para producción
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## Endpoints

| Método | Ruta | Auth | Descripción |
|---|---|---|---|
| `GET` | `/health` | ❌ Público | Estado del servicio |
| `GET` | `/docs/*` | ❌ Público | Swagger UI (desactivar en producción) |
| `POST` | `/api/v1/auth/register` | ❌ Público | Registra un nuevo usuario |
| `POST` | `/api/v1/auth/login` | ❌ Público | Inicia sesión y devuelve tokens JWT |
| `POST` | `/api/v1/auth/refresh` | ❌ Público | Renueva el access token con el refresh token |
| `POST` | `/api/v1/auth/logout` | ❌ Público | Invalida el refresh token |
| `GET` | `/api/v1/users/me` | ✅ Bearer token | Perfil del usuario autenticado |

---

## Rate limiting

Cada respuesta incluye headers informativos sobre el estado del límite:

| Header | Descripción |
|---|---|
| `X-Ratelimit-Limit` | Tasa configurada en req/s para ese endpoint |
| `X-Ratelimit-Remaining` | Tokens disponibles para esa IP en este momento |
| `Retry-After` | Segundos a esperar antes de reintentar (solo en 429) |

Las tres capas de límite configuradas en `router.go`:

| Capa | Aplica a | Límite | Burst |
|---|---|---|---|
| Global | Todo el servidor | 500 req/s | 1000 |
| Auth por IP | `/api/v1/auth/*` | 10 req/min | 5 |
| API por IP | `/api/v1/*` (protegido) | 60 req/s | 120 |

Desactiva temporalmente con `RATE_LIMIT_ENABLED=false` (útil en tests).

---

## Formato de respuesta

Todas las respuestas siguen la misma estructura JSON.

### Éxito

```json
{
  "success": true,
  "data": { "id": 1, "name": "Juan Pérez", "email": "juan@ejemplo.com" },
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Error

```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "user not found"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Paginación

```json
{
  "success": true,
  "data": [
    { "id": 1, "title": "Primer post" },
    { "id": 2, "title": "Segundo post" }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

Los códigos de error posibles: `BAD_REQUEST`, `UNAUTHORIZED`, `FORBIDDEN`, `NOT_FOUND`, `CONFLICT`, `TOO_MANY_REQUESTS`, `INTERNAL_ERROR`.

---

## Documentación Swagger

La documentación OpenAPI se genera automáticamente a partir de las anotaciones en los handlers.

```bash
# Instalar el CLI de swag (solo la primera vez)
go install github.com/swaggo/swag/cmd/swag@v1.8.12

# Regenerar docs después de modificar handlers
make swagger
```

La UI interactiva queda disponible en `http://localhost:8080/docs/index.html`. Para probar endpoints protegidos, haz click en **Authorize** e ingresa tu access token con el formato `Bearer <token>`.

> Recomendación: usa `SWAGGER_ENABLED=false` en producción para no exponer la documentación de tu API.

---

## Agregar un nuevo recurso

A continuación se muestra cómo agregar un recurso `Post` de forma consistente con la arquitectura del proyecto.

### 1. Crear el modelo

```go
// internal/models/post.go
package models

import "gorm.io/gorm"

type Post struct {
    gorm.Model
    Title   string `gorm:"not null"`
    Content string
    UserID  uint   `gorm:"not null;index"`
}
```

### 2. Registrar la migración

En `internal/database/database.go`, agrega el modelo al bloque de auto-migrate:

```go
func migrate() {
    DB.AutoMigrate(
        &models.User{},
        &models.RefreshToken{},
        &models.Post{},   // <-- añadir aquí
    )
}
```

### 3. Crear el servicio

```go
// internal/services/post_service.go
package services

import (
    "go-api-template/internal/database"
    "go-api-template/internal/models"
    "go-api-template/internal/utils"
)

type CreatePostInput struct {
    Title   string `json:"title"   binding:"required,min=3"`
    Content string `json:"content"`
}

func CreatePost(userID uint, input CreatePostInput) (*models.Post, error) {
    post := &models.Post{
        Title:   input.Title,
        Content: input.Content,
        UserID:  userID,
    }
    return post, database.DB.Create(post).Error
}

func GetPosts(params utils.PageParams) ([]models.Post, int64, error) {
    var posts []models.Post
    var total int64
    database.DB.Model(&models.Post{}).Count(&total)
    database.DB.Scopes(utils.Paginate(params)).Find(&posts)
    return posts, total, nil
}
```

### 4. Crear el handler

```go
// internal/handlers/post_handler.go
package handlers

import (
    "go-api-template/internal/services"
    "go-api-template/internal/utils"

    "github.com/gin-gonic/gin"
)

// GetPosts godoc
//
//	@Summary		Listar posts
//	@Tags			posts
//	@Produce		json
//	@Param			page      query  int  false  "Página"     default(1)
//	@Param			per_page  query  int  false  "Por página" default(20)
//	@Success		200  {object}  utils.PaginatedResponse
//	@Security		BearerAuth
//	@Router			/api/v1/posts [get]
func GetPosts(c *gin.Context) {
    params := utils.ParsePageParams(c)
    posts, total, err := services.GetPosts(params)
    if err != nil {
        utils.InternalError(c, "could not retrieve posts")
        return
    }
    utils.Paginated(c, posts, utils.NewMeta(params, total))
}

// CreatePost godoc
//
//	@Summary		Crear post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			request  body  services.CreatePostInput  true  "Datos del post"
//	@Success		201  {object}  utils.Response
//	@Security		BearerAuth
//	@Router			/api/v1/posts [post]
func CreatePost(c *gin.Context) {
    var input services.CreatePostInput
    if err := c.ShouldBindJSON(&input); err != nil {
        utils.BadRequest(c, err.Error())
        return
    }
    userID := c.MustGet("user_id").(uint)
    post, err := services.CreatePost(userID, input)
    if err != nil {
        utils.InternalError(c, "could not create post")
        return
    }
    utils.Created(c, post)
}
```

### 5. Registrar las rutas

En `internal/router/router.go`, dentro del grupo `protected`:

```go
protected.GET("/posts",  handlers.GetPosts)
protected.POST("/posts", handlers.CreatePost)
```

### 6. Regenerar la documentación

```bash
make swagger
```

---

## Docker

```bash
# Levantar el servicio con Docker Compose
make docker-up

# Detener y eliminar los contenedores
make docker-down
```

La base de datos SQLite persiste en el volumen `./data/data.db` en tu máquina host.

---

## Makefile

| Comando | Descripción |
|---|---|
| `make help` | Lista todos los comandos disponibles |
| `make run` | Ejecuta la app directamente con `go run` |
| `make build` | Compila el binario en `./bin/api` |
| `make swagger` | Genera/actualiza la documentación OpenAPI en `docs/` |
| `make test` | Ejecuta todos los tests con cobertura |
| `make test-race` | Ejecuta tests con detección de race conditions |
| `make tidy` | Ejecuta `go mod tidy` |
| `make lint` | Ejecuta `golangci-lint` (requiere instalación previa) |
| `make clean` | Elimina el binario compilado y `data.db` |
| `make air` | Inicia hot-reload con Air (requiere instalación previa) |
| `make docker-up` | Levanta el servicio con Docker Compose |
| `make docker-down` | Detiene y elimina los contenedores |

### Instalar herramientas de desarrollo opcionales

```bash
# Air (hot-reload)
go install github.com/air-verse/air@latest

# swag CLI (generación de docs Swagger)
go install github.com/swaggo/swag/cmd/swag@v1.8.12

# golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

---

## Licencia

Este proyecto está bajo la licencia **MIT**. Consulta el archivo [LICENSE](LICENSE) para más detalles.
