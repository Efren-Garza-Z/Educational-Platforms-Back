package main

import (
	"log"
	"os"

	"github.com/Efren-Garza-Z/go-api-gemini/db"
	_ "github.com/Efren-Garza-Z/go-api-gemini/docs" // swag docs
	"github.com/Efren-Garza-Z/go-api-gemini/domain/models"
	"github.com/Efren-Garza-Z/go-api-gemini/domain/repositories"
	service "github.com/Efren-Garza-Z/go-api-gemini/services"
	controllers "github.com/Efren-Garza-Z/go-api-gemini/web/controllers"
	"github.com/Efren-Garza-Z/go-api-gemini/web/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Escriba 'Bearer' seguido de un espacio y luego el token JWT. Ejemplo: "Bearer eyJhbGciOiJIUzI1NiIsInR5c..."
func main() {
    log.Println("🚀 Iniciando aplicación...")
    
    // 1. Conectar a la base de datos
    log.Println("📡 Conectando a la base de datos...")
    if err := db.Connect(); err != nil {
        log.Fatalf("❌ Error conectando a la base de datos: %v", err)
    }
    log.Println("✅ Conexión a base de datos exitosa")
    
    // 2. Migraciones
    log.Println("🔄 Ejecutando migraciones...")
    if err := db.DB.AutoMigrate(
        &models.UserDB{}, 
        &models.GeminiProcessingDB{}, 
        &models.GeminiProcessingFileDB{}, 
        &models.LearningInteractionDB{},
    ); err != nil {
        log.Fatalf("❌ Error al migrar modelos: %v", err)
    }
    log.Println("✅ Migraciones completadas")
    
    // 3. Inicializar repositorios y servicios
    log.Println("🏗️ Inicializando servicios...")
    userRepo := repositories.NewUserRepository(db.DB)
    gemRepo := repositories.NewGeminiRepository(db.DB)
    proRepo := repositories.NewProgressRepository(db.DB)
    
    userSvc := service.NewUserService(userRepo)
    proSvc := service.NewProgressService(proRepo)
    gemSvc := service.NewGeminiService(gemRepo, proSvc)
    
    // 4. Inicializar controladores
    userCtrl := controllers.NewUserController(userSvc, db.DB)
    gemCtrl := controllers.NewGeminiController(gemSvc)
    authCtrl := controllers.NewAuthController(userSvc)
    proCtrl := controllers.NewLearningController(gemSvc, userSvc, proSvc)
    
    // 5. Configurar Gin
    r := gin.Default()
    
    // Health check endpoint (IMPORTANTE para Cloud Run)
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // Swagger
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    // Routes
    routes.RegisterUserRoutes(r, userCtrl)
    routes.RegisterGeminiRoutes(r, gemCtrl)
    routes.RegisterAuthRoutes(r, authCtrl)
    routes.RegisterLearningRoutes(r, proCtrl)
    
    // 6. Obtener puerto
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("🌐 Servidor iniciando en puerto %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatalf("❌ Error al iniciar servidor: %v", err)
    }
}
