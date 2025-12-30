package main

import (
	"log"
	"os"
	
	"github.com/Efren-Garza-Z/go-api-gemini/db"
	_ "github.com/Efren-Garza-Z/go-api-gemini/docs"
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
func main() {
	log.Println("==========================================")
	log.Println("🚀 INICIANDO EDUCATIONAL PLATFORMS BACKEND")
	log.Println("==========================================")
	
	// ✅ AHORA SÍ MANEJAMOS EL ERROR
	if err := db.Connect(); err != nil {
		log.Fatalf("❌ FALLO CRÍTICO EN CONEXIÓN DB: %v", err)
	}
	
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
	
	// Repositorios
	log.Println("🏗️ Inicializando repositorios...")
	userRepo := repositories.NewUserRepository(db.DB)
	gemRepo := repositories.NewGeminiRepository(db.DB)
	proRepo := repositories.NewProgressRepository(db.DB)
	
	// Services
	log.Println("🛠️ Inicializando servicios...")
	userSvc := service.NewUserService(userRepo)
	proSvc := service.NewProgressService(proRepo)
	gemSvc := service.NewGeminiService(gemRepo, proSvc)
	
	// Controllers
	log.Println("🎮 Inicializando controladores...")
	userCtrl := controllers.NewUserController(userSvc, db.DB)
	gemCtrl := controllers.NewGeminiController(gemSvc)
	authCtrl := controllers.NewAuthController(userSvc)
	proCtrl := controllers.NewLearningController(gemSvc, userSvc, proSvc)
	
	// Gin
	log.Println("🌐 Configurando servidor Gin...")
	r := gin.Default()
	
	// Health check endpoint (IMPORTANTE para Cloud Run)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "educational-platforms-back",
		})
	})
	log.Println("✅ Health check en /health")
	
	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("✅ Swagger en /swagger")
	
	// Routes
	log.Println("🛣️ Registrando rutas...")
	routes.RegisterUserRoutes(r, userCtrl)
	routes.RegisterGeminiRoutes(r, gemCtrl)
	routes.RegisterAuthRoutes(r, authCtrl)
	routes.RegisterLearningRoutes(r, proCtrl)
	log.Println("✅ Rutas registradas")
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Println("==========================================")
	log.Printf("✅ SERVIDOR LISTO EN PUERTO %s", port)
	log.Println("==========================================")
	
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Error al iniciar servidor: %v", err)
	}
}
