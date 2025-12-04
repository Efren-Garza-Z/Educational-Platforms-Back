package routes

import (
	"github.com/Efren-Garza-Z/go-api-gemini/web/controllers"
	"github.com/Efren-Garza-Z/go-api-gemini/web/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, uc *controllers.UserController) {
	users := r.Group("/users")
	{
		users.POST("", uc.CreateUser)
		users.GET("", uc.GetAll)

		authenticated := users.Group("/")
		authenticated.Use(middleware.AuthRequired()) // Aplicar el middleware a este grupo
		{
			authenticated.GET("/:id", uc.GetByID)
			authenticated.PUT("/:id", uc.Update)
			authenticated.DELETE("/:id", uc.Delete)
		}
	}
}
