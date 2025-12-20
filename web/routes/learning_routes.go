package routes

import (
	"github.com/Efren-Garza-Z/go-api-gemini/web/controllers"
	"github.com/Efren-Garza-Z/go-api-gemini/web/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterLearningRoutes(r *gin.Engine, lc *controllers.LearningController) {
	// Creamos un grupo protegido por JWT
	learning := r.Group("/learning")
	learning.Use(middleware.AuthRequired()) // Obligatorio estar logueado
	{
		// Endpoint de conversación
		learning.POST("/chat", lc.ChatWithTutor)
		learning.GET("/history", lc.GetHistory)
	}
}
