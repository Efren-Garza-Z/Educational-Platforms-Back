# Project Analysis & Architecture Documentation

## 📊 Comprehensive Project Overview

### Project Summary
**Educational Platforms Backend - Gemini API** is a Go-based RESTful API that combines user management with Google Gemini AI integration for asynchronous processing of prompts and file analysis.

---

## 🔍 Detailed Architecture Analysis

### 1. **Layered Architecture Pattern**

The project implements a clean **4-layer architecture**:

#### Layer 1: Presentation Layer (Controllers)
- **Location:** `web/controllers/`
- **Components:**
  - `UserController` - Handles HTTP requests related to user operations
  - `GeminiController` - Manages Gemini AI requests
- **Responsibility:** 
  - Parse incoming HTTP requests
  - Call appropriate services
  - Format and return HTTP responses
  - Handle HTTP status codes

#### Layer 2: Business Logic Layer (Services)
- **Location:** `services/`
- **Components:**
  - `UserService` - User-related business operations
  - `GeminiService` - AI processing logic
- **Responsibility:**
  - Implement core business rules
  - Orchestrate repository calls
  - Handle validation
  - Process async tasks

#### Layer 3: Data Access Layer (Repositories)
- **Location:** `domain/repositories/`
- **Components:**
  - `UserRepository` - Database operations for users
  - `GeminiRepository` - Database operations for Gemini processing
- **Responsibility:**
  - Abstract database interactions
  - Provide CRUD operations
  - Implement data persistence logic

#### Layer 4: Domain Layer (Models)
- **Location:** `domain/models/`
- **Components:**
  - `UserDB` - User entity with GORM annotations
  - `GeminiProcessingDB` - Gemini task entity
  - `GeminiProcessingFileDB` - File processing entity
- **Responsibility:**
  - Define data structures
  - Include validation tags
  - Represent domain entities

---

## 📦 Dependency Injection & Initialization

### Main Function Flow
```
main.go
  ↓
db.Connect()          [Establishes PostgreSQL connection]
  ↓
AutoMigrate()         [Creates/updates database tables]
  ↓
Repository Init       [UserRepository, GeminiRepository]
  ↓
Service Init          [UserService, GeminiService]
  ↓
Controller Init       [UserController, GeminiController]
  ↓
Route Registration    [Sets up Gin routes]
  ↓
Swagger Setup         [Initializes documentation]
  ↓
Server Start          [Listening on configured port]
```

---

## 🔗 Component Interaction Flow

### User CRUD Operation Example
```
HTTP Request (POST /users)
    ↓
UserController.CreateUser()
    ↓
UserService.CreateUser()
    ↓
UserRepository.Create()
    ↓
PostgreSQL Database
    ↓
[Response Returns Through Same Path]
```

### Async Gemini Processing Example
```
HTTP Request (POST /gemini/process)
    ↓
GeminiController.ProcessPrompt()
    ↓
GeminiService.ProcessPromptAsync()
    │
    ├─→ Generate UUID
    ├─→ Create GeminiProcessingDB record (status: "pendiente")
    └─→ Return Task ID immediately (202 Accepted)
    
Background Processing:
    ↓
Update Status to "en_proceso"
    ↓
Call Google Gemini API
    ↓
Update Status to "finalizado" or "error"
    ↓
Store Result in Database
    
Status Check:
HTTP Request (GET /gemini/status/{id})
    ↓
GeminiController.GetTaskStatus()
    ↓
GeminiService.GetProcessStatus()
    ↓
GeminiRepository.FindByID()
    ↓
Return Current Status & Result
```

---

## 📋 Data Models Deep Dive

### User Model Structure
```
UserDB Table Structure (service.users):
┌─────────────┬──────────────┬─────────────────────────────┐
│ Column      │ Type         │ Constraints                 │
├─────────────┼──────────────┼─────────────────────────────┤
│ id          │ SERIAL       │ PRIMARY KEY, AUTO INCREMENT │
│ created_at  │ TIMESTAMP    │ Auto set by GORM            │
│ updated_at  │ TIMESTAMP    │ Auto set by GORM            │
│ full_name   │ VARCHAR      │ NOT NULL                    │
│ email       │ VARCHAR      │ NOT NULL, UNIQUE INDEX      │
│ password    │ VARCHAR      │ NOT NULL                    │
└─────────────┴──────────────┴─────────────────────────────┘

View Model (User):
- Does NOT include password
- Used for API responses to prevent information leakage
```

### Gemini Processing Model
```
GeminiProcessingDB Table Structure (service.gemini_processing):
┌─────────────┬──────────────────┬─────────────────────────────┐
│ Column      │ Type             │ Constraints                 │
├─────────────┼──────────────────┼─────────────────────────────┤
│ id          │ VARCHAR(36)      │ PRIMARY KEY (UUID)          │
│ created_at  │ TIMESTAMP        │ Auto set by GORM            │
│ updated_at  │ TIMESTAMP        │ Auto set by GORM            │
│ status      │ VARCHAR(20)      │ NOT NULL                    │
│ result      │ TEXT             │ NULL (populated on success) │
│ error       │ TEXT             │ NULL (populated on error)   │
│ prompt      │ TEXT             │ NOT NULL                    │
└─────────────┴──────────────────┴─────────────────────────────┘

Status Enum:
- pendiente     → Task created, waiting for processing
- en_proceso    → Currently being processed by Gemini
- finalizado    → Successfully completed
- error         → Error occurred during processing
```

---

## 🔐 Security Considerations

### Current Implementation
1. **Input Validation**
   - Email validation using Gin's `binding:"email"`
   - Required field validation using `binding:"required"`
   - Type safety with Go's strong typing

2. **Data Protection**
   - Passwords NOT hashed (⚠️ SECURITY RISK)
   - User model hides password in public responses

### Security Recommendations
1. **Priority 1 - Critical**
   - Implement password hashing with `golang.org/x/crypto/bcrypt`
   - Add JWT-based authentication
   - Implement HTTPS/TLS
   - Add CORS configuration

2. **Priority 2 - Important**
   - Add rate limiting
   - Implement request logging
   - Add input sanitization
   - File upload validation

3. **Priority 3 - Nice to Have**
   - Add request tracing
   - Implement audit logging
   - Add API versioning
   - Implement caching strategies

---

## 🚀 API Endpoint Mapping

### User Management Endpoints
```
POST   /users                  → Create new user
GET    /users                  → List all users
GET    /users/{id}             → Get user by ID
PUT    /users/{id}             → Update user
DELETE /users/{id}             → Delete user
```

### Gemini Processing Endpoints
```
POST   /gemini/process         → Submit prompt for async processing
GET    /gemini/status/{id}     → Check processing status
POST   /gemini/process-file    → Submit file with prompt
GET    /gemini/status/{id}     → Check file processing status
```

### Documentation & Health
```
GET    /swagger/*              → Swagger UI documentation
GET    /swagger/index.html     → Direct access to Swagger
```

---

## 🛠️ Technology Stack Analysis

### Backend Framework
- **Gin Framework v1.10.1**
  - Lightweight HTTP web framework
  - Fast routing with radix tree
  - Middleware support
  - Built-in validation
  - Error handling

### Database Layer
- **GORM v1.30.3**
  - Object-Relational Mapping
  - Automatic migrations
  - Relationship management
  - Query builder
  
- **PostgreSQL Driver**
  - Native support through `gorm.io/driver/postgres`
  - Connection pooling
  - Transaction support

### AI Integration
- **Google Genai SDK v1.23.0**
  - Official Google AI SDK
  - Gemini model support
  - File upload capability
  - Streaming responses

### Documentation
- **Swagger/OpenAPI v1.6.1**
  - Auto-generated from code annotations
  - Interactive testing UI
  - Schema documentation
  - Multi-format export (JSON, YAML)

### Utilities
- **UUID v1.6.0** - Unique identifier generation
- **godotenv v1.5.1** - Environment variable loading

---

## 📊 Database Schema

### Schema Diagram
```
Database: your_db_name
Schema: service

Tables:
┌─────────────────────────┐
│     service.users       │
├─────────────────────────┤
│ id (PK)                 │
│ created_at              │
│ updated_at              │
│ full_name               │
│ email (UNIQUE)          │
│ password                │
└─────────────────────────┘

┌────────────────────────────────┐
│ service.gemini_processing      │
├────────────────────────────────┤
│ id (PK) - UUID                 │
│ created_at                     │
│ updated_at                     │
│ status                         │
│ result                         │
│ error                          │
│ prompt                         │
└────────────────────────────────┘

┌────────────────────────────────────┐
│ service.gemini_processing_files    │
├────────────────────────────────────┤
│ id (PK) - UUID                     │
│ created_at                         │
│ updated_at                         │
│ status                             │
│ result                             │
│ error                              │
│ prompt                             │
│ file_name                          │
│ file_content_type                  │
└────────────────────────────────────┘
```

---

## 🔄 Request/Response Flow

### User Creation Flow
```
Request:
POST /users HTTP/1.1
Content-Type: application/json

{
  "full_name": "John Doe",
  "email": "john@example.com",
  "password": "secure123"
}

Processing:
1. Controller receives request
2. Binds JSON to CreateUserInput struct
3. Validates input (required fields, email format)
4. Calls UserService.CreateUser()
5. Service creates UserDB instance
6. Service calls UserRepository.Create()
7. Repository executes SQL INSERT
8. Returns created user

Response:
HTTP/1.1 201 Created
Content-Type: application/json

{
  "id": 1,
  "full_name": "John Doe",
  "email": "john@example.com"
  // Password NOT included in response
}
```

### Async Gemini Processing Flow
```
Request:
POST /gemini/process HTTP/1.1
Content-Type: application/json

{
  "prompt": "What are the best universities in USA?"
}

Immediate Response:
HTTP/1.1 202 Accepted
Content-Type: application/json

{
  "task_id": "550e8400-e29b-41d4-a716-446655440000"
}

Background Processing (Client polls):
GET /gemini/status/550e8400-e29b-41d4-a716-446655440000

Response 1 (Still processing):
HTTP/1.1 200 OK

{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "en_proceso",
  "result": null,
  "error": null
}

Response 2 (After completion):
HTTP/1.1 200 OK

{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "finalizado",
  "result": "The best universities in the USA include Harvard, MIT, Stanford...",
  "error": null
}
```

---

## 📈 Scalability Considerations

### Current Limitations
1. **Synchronous operations** for user management
2. **Single instance** deployment
3. **No caching layer**
4. **No rate limiting**
5. **Direct database calls** without connection pooling optimization

### Recommendations for Scaling
1. **Horizontal Scaling**
   - Add load balancer (Nginx, HAProxy)
   - Implement database connection pooling
   - Use distributed caching (Redis)

2. **Asynchronous Processing**
   - Message queue (RabbitMQ, Kafka) for Gemini tasks
   - Worker processes for batch processing
   - Job scheduling system

3. **Caching Strategy**
   - Cache frequently accessed users
   - Cache Gemini results
   - Implement cache invalidation

4. **Database Optimization**
   - Implement indexing strategy
   - Connection pooling
   - Query optimization
   - Database replication

---

## 🧪 Testing Strategy Recommendations

### Unit Testing
```go
// Test UserService logic
func TestCreateUser(t *testing.T) {
  // Mock UserRepository
  // Call CreateUser
  // Verify results
}

// Test GeminiService logic
func TestProcessPromptAsync(t *testing.T) {
  // Mock GeminiRepository
  // Call ProcessPromptAsync
  // Verify task creation
}
```

### Integration Testing
```go
// Test UserController endpoints
func TestUserCreateEndpoint(t *testing.T) {
  // Start test server
  // Make HTTP request
  // Verify database changes
}
```

### End-to-End Testing
```go
// Test full workflows
func TestCompleteUserFlow(t *testing.T) {
  // Create user
  // Retrieve user
  // Update user
  // Delete user
  // Verify all operations
}
```

---

## 📝 Code Quality Observations

### Strengths
✅ Clear separation of concerns  
✅ Dependency injection pattern  
✅ Consistent naming conventions  
✅ Use of interfaces for repositories  
✅ Proper HTTP status codes  
✅ Environment variable configuration  

### Areas for Improvement
⚠️ No error handling middleware  
⚠️ Password stored in plaintext  
⚠️ No input sanitization  
⚠️ Limited logging  
⚠️ No tests included  
⚠️ No transaction management  
⚠️ No pagination on list endpoints  

---

## 🎯 Development Priorities

### Short Term (Week 1-2)
- [ ] Add password hashing
- [ ] Add JWT authentication
- [ ] Implement error handling middleware
- [ ] Add request logging

### Medium Term (Week 3-4)
- [ ] Add unit tests (target 80% coverage)
- [ ] Implement pagination
- [ ] Add rate limiting
- [ ] Implement caching

### Long Term (Month 2+)
- [ ] Message queue integration
- [ ] Distributed processing
- [ ] Advanced monitoring
- [ ] Performance optimization

---

## 📚 References & Resources

### Go Packages Used
- Gin: https://github.com/gin-gonic/gin
- GORM: https://gorm.io/
- Google Genai: https://github.com/google/generative-ai-go
- UUID: https://github.com/google/uuid
- Godotenv: https://github.com/joho/godotenv

### Best Practices
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Gin Documentation](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/docs/)

---

*This analysis was generated on January 29, 2026*
