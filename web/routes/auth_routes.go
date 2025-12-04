package routes

import (
	"github.com/Efren-Garza-Z/go-api-gemini/web/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, ac *controllers.AuthController) {
	auth := r.Group("/auth")
	{
		// Ruta de Login
		auth.POST("/login", ac.Login)

		// La ruta de Registro (CreateUser) ya existe en UserController,
		// pero podrías moverla aquí si lo deseas para agrupar mejor la autenticación.
		// Por ahora, la dejamos en /users.
	}
}
