package controllers

import (
	"net/http"

	"github.com/Efren-Garza-Z/go-api-gemini/domain/models"
	"github.com/Efren-Garza-Z/go-api-gemini/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LearningController struct {
	geminiService   services.GeminiService
	userService     services.UserService
	progressService services.ProgressService
}

func NewLearningController(
	gs services.GeminiService,
	us services.UserService,
	ps services.ProgressService,
) *LearningController {
	return &LearningController{
		geminiService:   gs,
		userService:     us,
		progressService: ps,
	}
}

// @Summary Iniciar tutoría de conversación con IA
// @Tags learning
// @Accept json
// @Produce json
// @Param input body models.PromptRequest true "Mensaje del estudiante"
// @Security ApiKeyAuth
// @Success 202 {object} models.GeminiProcessingIDResponse
// @Router /learning/chat [post]
func (lc *LearningController) ChatWithTutor(c *gin.Context) {

	// 1️⃣ Obtener userID del JWT
	val, _ := c.Get("userID")
	userID := val.(uint)

	// 2️⃣ Obtener usuario real
	user, err := lc.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// 3️⃣ Leer request
	var req models.PromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mensaje inválido"})
		return
	}

	// 4️⃣ Idioma y nivel (fallbacks)
	lang := user.TargetLanguage
	if lang == "" {
		lang = "English"
	}

	lvl := user.LanguageLevel
	if lvl == "" {
		lvl = "A1"
	}

	// 5️⃣ ConversationID (nuevo o existente)
	conversationID := req.ConversationID
	if conversationID == "" {
		conversationID = uuid.New().String()
	}

	// 6️⃣ Llamar al service (SIN contextualPrompt)
	id, err := lc.geminiService.ProcessChatAsync(
		userID,
		conversationID,
		lang,
		lvl,
		req.Prompt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar con Gemini"})
		return
	}

	// 7️⃣ Responder
	c.JSON(http.StatusAccepted, models.GeminiProcessingIDResponse{
		GeminiProcessingID: id,
	})
}

// GetHistory recupera todas las interacciones de aprendizaje del usuario logueado.
// @Summary Obtener historial de aprendizaje
// @Tags learning
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.LearningInteractionDB
// @Router /learning/history [get]
func (lc *LearningController) GetHistory(c *gin.Context) {
	// 1. Obtener el UserID del token JWT
	val, _ := c.Get("userID")
	userID := val.(uint)

	// 2. Llamar al servicio de progreso
	history, err := lc.progressService.GetHistoryByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo recuperar el historial"})
		return
	}

	// 3. Retornar la lista (puede ser una lista vacía [] si no hay registros)
	c.JSON(http.StatusOK, history)
}
