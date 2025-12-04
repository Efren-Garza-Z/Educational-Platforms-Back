package controllers

import (
	"net/http"

	"github.com/Efren-Garza-Z/go-api-gemini/domain/models"
	"github.com/Efren-Garza-Z/go-api-gemini/services"
	"github.com/gin-gonic/gin"
)

// AuthController maneja las solicitudes relacionadas con la autenticación (Login/Logout).
type AuthController struct {
	// Necesita el UserService para verificar credenciales
	userService services.UserService
	// Necesitarás un servicio de JWT o una función aquí para generar el token (Fase 1.2/1.3 del plan).
	// Por ahora, lo dejaremos simple.
}

func NewAuthController(us services.UserService) *AuthController {
	return &AuthController{userService: us}
}

// @Summary Iniciar sesión de usuario
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.LoginInput true "Credenciales de inicio de sesión"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON de entrada inválido o faltante"})
		return
	}

	// 1. Llamar al servicio para verificar credenciales (Busca por email y compara el hash)
	user, err := ac.userService.Login(input.Email, input.Password)

	if err != nil {
		// El servicio devuelve un error genérico si las credenciales son malas o hay un error interno.
		// Es buena práctica no especificar si falló el email o la contraseña.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email o contraseña inválidos"})
		return
	}

	token, err := ac.userService.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el token de sesión"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token:  token,
		UserID: user.ID,
	})
}
