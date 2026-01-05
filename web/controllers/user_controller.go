package controllers

import (
	"net/http"
	"strconv"

	"github.com/Efren-Garza-Z/go-api-gemini/domain/models"
	"github.com/Efren-Garza-Z/go-api-gemini/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	service services.UserService
	db      *gorm.DB
}

func NewUserController(s services.UserService, db *gorm.DB) *UserController {
	return &UserController{service: s, db: db}
}

// @Summary Crear usuario
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.CreateUserInput true "Datos para crear usuario"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	var input models.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := uc.service.CreateUser(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear usuario"})
		return
	}
	c.JSON(http.StatusCreated, u.ToPublic())
}

// @Summary Obtener todos los usuarios
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func (uc *UserController) GetAll(c *gin.Context) {
	users, err := uc.service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios"})
		return
	}
	var out []models.User
	for _, u := range users {
		out = append(out, u.ToPublic())
	}
	c.JSON(http.StatusOK, out)
}

// @Summary Obtener usuario por ID
// @Tags users
// @Param id path int true "ID del usuario"
// @Produce json
// @Success 200 {object} models.User
// @Router /users/{id} [get]
// @security ApiKeyAuth
func (uc *UserController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	u, err := uc.service.GetUserByID(uint(id64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u.ToPublic())
}

// @Summary Obtener usuario por Email
// @Tags users
// @Param email formData string true "Email del usuario"
// @Produce json
// @Success 200 {object} models.User
// @Router /users/{email} [get]
// @security ApiKeyAuth
func (uc *UserController) GetByEmail(c *gin.Context) {
	// 1. Obtener el email de la URL (el parámetro debe coincidir con el nombre en la ruta)
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El email es requerido"})
		return
	}

	u, err := uc.service.GetUser(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u.ToPublic())
}

// @Summary Actualizar usuario
// @Tags users
// @Accept json
// @Produce json
// @Param email path string true "Email del usuario"
// @Param input body models.CreateUserInput true "Datos para actualizar usuario"
// @Success 200 {object} models.User
// @Router /users/email/{email} [put]
// @security ApiKeyAuth
func (uc *UserController) Update(c *gin.Context) {
	// 1. Obtener el email de la URL (el parámetro debe coincidir con el nombre en la ruta)
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El email es requerido"})
		return
	}

	// 2. Validar el cuerpo de la solicitud (JSON)
	var input models.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos: " + err.Error()})
		return
	}

	// 3. Llamar al servicio usando el email
	u, err := uc.service.UpdateUser(email, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Responder con el modelo público (para no mostrar contraseñas)
	c.JSON(http.StatusOK, u.ToPublic())
}

// @Summary Eliminar usuario
// @Tags users
// @Param id path int true "ID del usuario"
// @Success 204
// @Router /users/{id} [delete]
// @security ApiKeyAuth
func (uc *UserController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	if err := uc.service.DeleteUser(uint(id64)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
