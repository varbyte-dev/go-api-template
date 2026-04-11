# go-api-template

![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=flat)
![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen?style=flat)

Boilerplate de producción para construir APIs REST en Go. Incluye autenticación JWT, ORM con GORM, middlewares esenciales, respuestas estandarizadas, hot-reload y Docker listo para desplegar.

---

## Características

- 🔐 **Autenticación JWT** — Access + refresh tokens con rotación automática
- 🗄️ **ORM con GORM + SQLite** — Fácil de cambiar a PostgreSQL o MySQL
- 🌐 **Gin Framework** — Con middlewares de CORS, Logger, RequestID y Recovery
- 📦 **Respuestas JSON estandarizadas** — Estructura `{success, data, error, request_id}` consistente
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
# APP_PORT, JWT_SECRET, etc.

# 4. Instalar dependencias y ejecutar
make run
```

La API estará disponible en `http://localhost:8080`.

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

> ⚠️ **Importante:** Cambia `JWT_SECRET` por un valor seguro y aleatorio antes de ir a producción.

---

## Estructura del proyecto

```
go-api-template/
├── cmd/
│   └── main.go                  # Punto de entrada: inicializa config, DB, router y servidor
├── internal/
│   ├── config/
│   │   └── config.go            # Carga y valida variables de entorno
│   ├── database/
│   │   └── database.go          # Conexión a SQLite y auto-migración de modelos
│   ├── handlers/
│   │   ├── auth_handler.go      # Handlers de registro, login, refresh y logout
│   │   ├── user_handler.go      # Handler de perfil del usuario autenticado
│   │   └── health_handler.go    # Health check del servicio
│   ├── middleware/
│   │   ├── auth.go              # Middleware de validación de JWT
│   │   ├── cors.go              # Middleware de CORS configurable
│   │   ├── logger.go            # Middleware de logging de requests
│   │   └── request_id.go        # Middleware que inyecta un UUID por request
│   ├── models/
│   │   └── user.go              # Modelo de usuario para GORM
│   ├── router/
│   │   └── router.go            # Registro de rutas y grupos de la API
│   ├── services/
│   │   └── auth_service.go      # Lógica de negocio de autenticación
│   └── utils/
│       ├── response.go          # Helpers para respuestas JSON estandarizadas
│       └── pagination.go        # Helper genérico de paginación con GORM
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
| `GET` | `/health` | ❌ Público | Verifica que el servicio está en línea |
| `POST` | `/api/v1/auth/register` | ❌ Público | Registra un nuevo usuario |
| `POST` | `/api/v1/auth/login` | ❌ Público | Inicia sesión y devuelve tokens JWT |
| `POST` | `/api/v1/auth/refresh` | ❌ Público | Renueva el access token usando el refresh token |
| `POST` | `/api/v1/auth/logout` | ❌ Público | Invalida el refresh token |
| `GET` | `/api/v1/users/me` | ✅ Bearer token | Devuelve el perfil del usuario autenticado |

---

## Formato de respuesta

Todas las respuestas siguen la misma estructura JSON para facilitar el manejo en el cliente.

### Éxito

```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "usuario@ejemplo.com",
    "name": "Juan Pérez"
  },
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

---

## Agregar un nuevo recurso

A continuación se muestra cómo agregar un recurso `Post` de forma consistente con la arquitectura del proyecto.

### 1. Crear el modelo

Crea `internal/models/post.go`:

```go
// internal/models/post.go
package models

import "gorm.io/gorm"

type Post struct {
    gorm.Model
    Title   string `json:"title"   gorm:"not null"`
    Content string `json:"content"`
    UserID  uint   `json:"user_id" gorm:"not null"`
}
```

### 2. Registrar la migración

En `internal/database/database.go`, agrega el modelo al auto-migrate:

```go
// Dentro de la función de migración
db.AutoMigrate(&models.User{}, &models.Post{})
```

### 3. Crear el servicio

Crea `internal/services/post_service.go` con la lógica de negocio:

```go
// internal/services/post_service.go
package services

import (
    "go-api-template/internal/models"
    "gorm.io/gorm"
)

type PostService struct {
    db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
    return &PostService{db: db}
}

func (s *PostService) GetAll(page, perPage int) ([]models.Post, int64, error) {
    var posts []models.Post
    var total int64
    s.db.Model(&models.Post{}).Count(&total)
    s.db.Offset((page - 1) * perPage).Limit(perPage).Find(&posts)
    return posts, total, nil
}

func (s *PostService) Create(post *models.Post) error {
    return s.db.Create(post).Error
}
```

### 4. Crear el handler

Crea `internal/handlers/post_handler.go`:

```go
// internal/handlers/post_handler.go
package handlers

import (
    "net/http"
    "go-api-template/internal/services"
    "go-api-template/internal/utils"
    "github.com/gin-gonic/gin"
)

type PostHandler struct {
    postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
    return &PostHandler{postService: postService}
}

func (h *PostHandler) GetAll(c *gin.Context) {
    posts, total, err := h.postService.GetAll(1, 20)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
        return
    }
    utils.PaginatedResponse(c, posts, total, 1, 20)
}
```

### 5. Registrar las rutas

En `internal/router/router.go`, inyecta el handler y registra las rutas:

```go
postService := services.NewPostService(db)
postHandler := handlers.NewPostHandler(postService)

// Dentro del grupo /api/v1 protegido
posts := api.Group("/posts")
posts.Use(middleware.AuthMiddleware(cfg))
{
    posts.GET("", postHandler.GetAll)
}
```

---

## Docker

El proyecto incluye un `Dockerfile` con multi-stage build para generar una imagen mínima y segura, y un `docker-compose.yml` listo para levantar el servicio.

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

# golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

---

## Licencia

Este proyecto está bajo la licencia **MIT**. Consulta el archivo [LICENSE](LICENSE) para más detalles.