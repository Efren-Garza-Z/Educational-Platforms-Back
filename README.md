# Educational Platforms Backend - Gemini API

> Una API RESTful moderna construida en Go que integra gestión de usuarios con el poder de la inteligencia artificial de Google Gemini para procesamiento asíncrono de prompts y análisis de archivos.

---

## 📋 Tabla de Contenidos

- [Descripción del Proyecto](#descripción-del-proyecto)
- [Características](#características)
- [Arquitectura](#arquitectura)
- [Requisitos Previos](#requisitos-previos)
- [Instalación](#instalación)
- [Configuración del Entorno](#configuración-del-entorno)
- [Ejecución](#ejecución)
- [Documentación de la API](#documentación-de-la-api)
- [Estructura del Proyecto](#estructura-del-proyecto)
- [Endpoints Disponibles](#endpoints-disponibles)
- [Modelos de Datos](#modelos-de-datos)
- [Tecnologías Utilizadas](#tecnologías-utilizadas)
- [Notas de Desarrollo](#notas-de-desarrollo)

---

## 🎯 Descripción del Proyecto

Esta es una plataforma educativa backend que proporciona dos funcionalidades principales:

1. **Gestión de Usuarios**: Creación, lectura, actualización y eliminación de usuarios
2. **Procesamiento Asíncrono con Gemini AI**: Envía prompts o archivos a Google Gemini para procesamiento inteligente con manejo asíncrono de tareas

La API está diseñada siguiendo el patrón **MVC (Model-View-Controller)** con separación clara entre capas de modelos, servicios y controladores, lo que garantiza un código mantenible y escalable.

---

## ✨ Características

- ✅ **API RESTful completa** con Gin Framework
- ✅ **Autenticación y gestión de usuarios** con validación de email
- ✅ **Procesamiento asíncrono** de prompts de IA
- ✅ **Soporte para archivos** (PDF, PNG, JPG) en procesamiento Gemini
- ✅ **Sistema de estado de tareas** (pendiente, en proceso, completado, error)
- ✅ **Documentación automática con Swagger/OpenAPI**
- ✅ **Base de datos PostgreSQL** con ORM GORM
- ✅ **Generación de IDs únicos** con UUID
- ✅ **Migraciones automáticas** de modelos

---

## 🏗️ Arquitectura

El proyecto sigue una arquitectura en capas con separación de responsabilidades:

```
┌─────────────────────────────────────┐
│      API HTTP (Gin Router)          │
├─────────────────────────────────────┤
│      Controllers (Request Handlers)  │
│  ├── UserController                 │
│  └── GeminiController               │
├─────────────────────────────────────┤
│      Services (Business Logic)       │
│  ├── UserService                    │
│  └── GeminiService                  │
├─────────────────────────────────────┤
│      Repositories (Data Access)      │
│  ├── UserRepository                 │
│  └── GeminiRepository               │
├─────────────────────────────────────┤
│      Domain Models (Entities)        │
│  ├── UserDB                         │
│  ├── GeminiProcessingDB             │
│  └── GeminiProcessingFileDB         │
├─────────────────────────────────────┤
│      Database (PostgreSQL + GORM)   │
└─────────────────────────────────────┘
```

**Beneficios:**
- Código modular y reutilizable
- Fácil de testear cada capa
- Cambios en la BD no afectan la lógica de negocio
- Clara separación de responsabilidades

---

## 📋 Requisitos Previos

Antes de comenzar, asegúrate de tener instalado:

- **Go** 1.25.0 o superior ([descargar](https://golang.org/dl/))
- **PostgreSQL** 12 o superior ([descargar](https://www.postgresql.org/download/))
- **Git** ([descargar](https://git-scm.com/))
- **API Key de Google Gemini** ([obtener aquí](https://ai.google.dev/))

**Verificar instalación de Go:**
```bash
go version
```

---

## 🚀 Instalación

### 1. Clonar el repositorio

```bash
git clone https://github.com/Efren-Garza-Z/go-api-gemini.git
cd go-api-gemini
```

### 2. Instalar dependencias

```bash
go mod tidy
```

Este comando descargará e instalará todas las dependencias especificadas en `go.mod`.

---

## ⚙️ Configuración del Entorno

### Crear archivo `.env`

En la raíz del proyecto, crea un archivo `.env` con la siguiente estructura:

```ini
# Base de datos PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name

# Google Gemini API
GEMINI_API_KEY=your_gemini_api_key

# Puerto de la aplicación (opcional)
PORT=8080
```

**Variables explicadas:**

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| `DB_HOST` | Host del servidor PostgreSQL | `localhost` |
| `DB_PORT` | Puerto del servidor PostgreSQL | `5432` |
| `DB_USER` | Usuario de PostgreSQL | `postgres` |
| `DB_PASSWORD` | Contraseña de PostgreSQL | `secure_password` |
| `DB_NAME` | Nombre de la base de datos | `gemini_db` |
| `GEMINI_API_KEY` | Clave API de Google Gemini | `AIzaSy...` |
| `PORT` | Puerto en el que corre la app | `8080` |

### Crear base de datos en PostgreSQL

```bash
psql -U postgres
```

```sql
CREATE DATABASE your_db_name;
CREATE SCHEMA service;
```

---

## 🎬 Ejecución

### Ejecutar en modo desarrollo

```bash
go run main.go
```

La aplicación se iniciará en `http://localhost:8080`

**Salida esperada:**
```
Servidor corriendo en http://localhost:8080
```

### Compilar a ejecutable

```bash
go build -o api-gemini
./api-gemini
```

---

## 📖 Documentación de la API

La API utiliza **Swagger/OpenAPI** para generar documentación interactiva automáticamente.

### Instalar herramienta Swagger

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

**Nota para usuarios de Linux (Pop!_OS):** Si el comando `swag` no se reconoce directamente, ejecuta:
```bash
$(go env GOPATH)/bin/swag init
```

### Generar documentación de Swagger

```bash
swag init
```

Este comando:
- Lee las anotaciones en los controladores
- Genera/actualiza `docs/swagger.json`
- Prepara la documentación interactiva

### Acceder a Swagger UI

Una vez que la aplicación está corriendo, abre tu navegador en:

```
http://localhost:8080/swagger/index.html
```

Aquí podrás:
- Ver todos los endpoints disponibles
- Probar cada endpoint directamente
- Ver ejemplos de request/response
- Explorar modelos de datos

---

## 📁 Estructura del Proyecto

```
.
├── main.go                          # Punto de entrada de la aplicación
├── go.mod                           # Dependencias del proyecto
├── go.sum                           # Checksums de dependencias
├── .env                             # Variables de entorno (no versionar)
├── README.md                        # Este archivo
│
├── db/
│   └── db.go                        # Conexión y configuración de GORM
│
├── domain/
│   ├── models/
│   │   ├── user_model.go            # Modelo de Usuario (UserDB, User, CreateUserInput)
│   │   ├── gemini_processing.go     # Modelo de Procesamiento Gemini
│   │   └── gemini_processing_file.go # Modelo de Procesamiento con Archivos
│   └── repositories/
│       ├── user_repository.go       # Interfaz y implementación de acceso a usuarios
│       └── gemini_repository.go     # Interfaz y implementación para Gemini
│
├── services/
│   ├── user_service.go              # Lógica de negocio para usuarios
│   └── gemini_service.go            # Lógica de negocio para Gemini
│
├── web/
│   ├── controllers/
│   │   ├── user_controller.go       # Manejo de requests de usuarios
│   │   └── gemini_controller.go     # Manejo de requests de Gemini
│   └── routes/
│       ├── user_routes.go           # Rutas de usuarios
│       └── gemini_routes.go         # Rutas de Gemini
│
└── docs/
    ├── docs.go                      # Generado automáticamente por Swagger
    ├── swagger.json                 # Especificación OpenAPI (generado)
    └── swagger.yaml                 # Especificación OpenAPI (generado)
```

---

## 🔌 Endpoints Disponibles

### 👥 Usuarios

#### Crear usuario
```
POST /users
Content-Type: application/json

{
  "full_name": "Juan Pérez",
  "email": "juan@example.com",
  "password": "SecurePassword123"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "full_name": "Juan Pérez",
  "email": "juan@example.com"
}
```

#### Obtener todos los usuarios
```
GET /users
```

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "full_name": "Juan Pérez",
    "email": "juan@example.com"
  }
]
```

#### Obtener usuario por ID
```
GET /users/{id}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "full_name": "Juan Pérez",
  "email": "juan@example.com"
}
```

#### Actualizar usuario
```
PUT /users/{id}
Content-Type: application/json

{
  "full_name": "Juan Pérez Actualizado",
  "email": "juannuevo@example.com",
  "password": "NewPassword123"
}
```

#### Eliminar usuario
```
DELETE /users/{id}
```

---

### 🤖 Gemini

#### Procesar prompt (Asíncrono)
```
POST /gemini/process
Content-Type: application/json

{
  "prompt": "¿Cuáles son las mejores universidades de Estados Unidos?"
}
```

**Response (202 Accepted):**
```json
{
  "task_id": "8b9a1d2e-3c4f-5a6b-7c8d-9e0f1a2b3c4d"
}
```

Utiliza este `task_id` para consultar el estado.

#### Obtener estado de procesamiento
```
GET /gemini/status/{gemini_processing_id}
```

**Response (200 OK):**
```json
{
  "id": "8b9a1d2e-3c4f-5a6b-7c8d-9e0f1a2b3c4d",
  "status": "finalizado",
  "result": "Las mejores universidades de EE.UU. incluyen...",
  "error": null
}
```

**Posibles estados:**
- `pendiente`: Esperando ser procesado
- `en_proceso`: Actualmente siendo procesado
- `finalizado`: Completado exitosamente
- `error`: Ocurrió un error durante el procesamiento

#### Procesar archivo con prompt
```
POST /gemini/process-file
Content-Type: multipart/form-data

prompt: "Analiza este documento y extrae puntos clave"
file: [archivo PDF, PNG o JPG]
```

**Response (202 Accepted):**
```json
{
  "file_processing_id": "abc123def456..."
}
```

---

## 📊 Modelos de Datos

### Usuario (UserDB)

```go
type UserDB struct {
  ID        uint      `gorm:"primaryKey"`      // ID único (autoincremental)
  CreatedAt time.Time                          // Fecha de creación
  UpdatedAt time.Time                          // Fecha de última actualización
  FullName  string    `gorm:"not null"`        // Nombre completo
  Email     string    `gorm:"uniqueIndex"`     // Email único
  Password  string    `gorm:"not null"`        // Contraseña
}
```

**Tabla:** `service.users`

---

### Procesamiento Gemini (GeminiProcessingDB)

```go
type GeminiProcessingDB struct {
  ID        string                 `gorm:"primaryKey"` // UUID único
  CreatedAt time.Time
  UpdatedAt time.Time
  Status    GeminiProcessingStatus                     // Estado del procesamiento
  Result    string                 `gorm:"type:text"`  // Resultado de Gemini
  Error     string                 `gorm:"type:text"`  // Mensaje de error (si aplica)
  Prompt    string                 `gorm:"type:text"`  // Prompt original
}
```

**Tabla:** `service.gemini_processing`

**Estados permitidos:**
- `pendiente`
- `en_proceso`
- `finalizado`
- `error`

---

### Procesamiento con Archivo (GeminiProcessingFileDB)

```go
type GeminiProcessingFileDB struct {
  ID              string                 `gorm:"primaryKey"`
  CreatedAt       time.Time
  UpdatedAt       time.Time
  Status          GeminiProcessingStatus
  Result          string                 `gorm:"type:text"`
  Error           string                 `gorm:"type:text"`
  Prompt          string                 `gorm:"type:text"`
  FileName        string                 // Nombre del archivo
  FileContentType string                 // Tipo MIME (pdf, image/png, etc)
}
```

**Tabla:** `service.gemini_processing_files`

---

## 🛠️ Tecnologías Utilizadas

| Tecnología | Versión | Propósito |
|-----------|---------|----------|
| **Go** | 1.25.0 | Lenguaje de programación |
| **Gin Framework** | v1.10.1 | Framework web HTTP |
| **GORM** | v1.30.3 | ORM para Go |
| **PostgreSQL** | 12+ | Base de datos principal |
| **Google Gemini SDK** | v1.23.0 | API de inteligencia artificial |
| **Swagger/OpenAPI** | v1.6.1 | Documentación interactiva |
| **UUID** | v1.6.0 | Generación de IDs únicos |
| **godotenv** | v1.5.1 | Carga de variables de entorno |

---

## 📝 Notas de Desarrollo

### Migraciones automáticas
Al iniciar la aplicación, se ejecutan automáticamente las migraciones para crear/actualizar las tablas:

```go
db.DB.AutoMigrate(&models.UserDB{}, 
                  &models.GeminiProcessingDB{}, 
                  &models.GeminiProcessingFileDB{})
```

### Validación de entrada
Los modelos incluyen etiquetas `binding` para validación automática con Gin:

```go
type CreateUserInput struct {
  Email string `binding:"required,email"` // Valida que sea email válido
}
```

### Seguridad
⚠️ **Importante:** 
- Las contraseñas se almacenan en texto plano (TODO: implementar hashing con bcrypt)
- No exponer el archivo `.env` en repositorios públicos
- Usar variables de entorno para credenciales sensibles

### Logging
La aplicación usa el logging estándar de Go. Para más detalles, usa:

```go
log.Printf("Tu mensaje: %v", valor)
```

### Rate Limiting (Recomendado para producción)
Considera agregar rate limiting con:
```bash
go get github.com/gin-contrib/ratelimit
```

---

## 🤝 Contribuciones

Para contribuir al proyecto:

1. Fork el repositorio
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

---

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver archivo `LICENSE` para más detalles.

---

## 📞 Contacto

**Autor:** Efren David Garza Z.  
**Email:** [Tu email aquí]  
**GitHub:** [@Efren-Garza-Z](https://github.com/Efren-Garza-Z)

---

## 🎓 Próximas Mejoras

- [ ] Implementar hashing de contraseñas con bcrypt
- [ ] Agregar autenticación JWT
- [ ] Implementar rate limiting
- [ ] Agregar tests unitarios
- [ ] Agregar logs estructurados
- [ ] Configuración de CORS
- [ ] Validación más robusta de archivos
- [ ] Caché de resultados
- [ ] Paginación en endpoints de listado
