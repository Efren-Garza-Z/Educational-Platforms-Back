package services

import (
	"github.com/Efren-Garza-Z/go-api-gemini/domain/models"
	"github.com/Efren-Garza-Z/go-api-gemini/domain/repositories"
)

// ProgressService define los métodos de negocio para el progreso del usuario.
type ProgressService interface {
	SaveInteraction(input models.LearningInteractionInput) (*models.LearningInteractionDB, error)
	GetHistoryByUserID(userID uint) ([]models.LearningInteractionDB, error)
	BuildConversationContext(
		userID uint,
		conversationID string,
	) (string, error)
}

type progressService struct {
	repo repositories.ProgressRepository
}

func NewProgressService(r repositories.ProgressRepository) ProgressService {
	return &progressService{repo: r}
}

// SaveInteraction toma el input del controlador, lo mapea al modelo de DB y lo persiste.
func (s *progressService) SaveInteraction(input models.LearningInteractionInput) (*models.LearningInteractionDB, error) {
	interaction := &models.LearningInteractionDB{
		ConversationID:  input.ConversationID,
		UserID:          input.UserID,
		InteractionType: input.InteractionType,
		Language:        input.Language,
		Level:           input.Level,
		Prompt:          input.Prompt,
		Response:        input.Response,
	}

	// Aquí podrías agregar más lógica de negocio, como validar el nivel o tipo antes de guardar.

	if err := s.repo.Create(interaction); err != nil {
		return nil, err
	}
	return interaction, nil
}

// GetHistoryByUserID recupera todas las interacciones de aprendizaje de un usuario.
func (s *progressService) GetHistoryByUserID(userID uint) ([]models.LearningInteractionDB, error) {
	return s.repo.FindAllByUserID(userID)
}

func (s *progressService) BuildConversationContext(
	userID uint,
	conversationID string,
) (string, error) {

	history, err := s.repo.FindByConversationID(userID, conversationID)
	if err != nil {
		return "", err
	}

	var context string
	for _, h := range history {
		context += "Student: " + h.Prompt + "\n"
		context += "Tutor: " + h.Response + "\n"
	}

	return context, nil
}
