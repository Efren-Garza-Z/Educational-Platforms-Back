package repositories

import (
	"github.com/Efren-Garza-Z/go-api-gemini/domain/models"
	"gorm.io/gorm"
)

// ProgressRepository define la interfaz para la persistencia de datos de progreso.
type ProgressRepository interface {
	Create(interaction *models.LearningInteractionDB) error
	FindAllByUserID(userID uint) ([]models.LearningInteractionDB, error)
	FindByConversationID(
		userID uint,
		conversationID string,
	) ([]models.LearningInteractionDB, error)
}

type progressRepository struct {
	db *gorm.DB
}

func NewProgressRepository(db *gorm.DB) ProgressRepository {
	return &progressRepository{db: db}
}

// Create guarda una nueva interacción en la base de datos.
func (r *progressRepository) Create(interaction *models.LearningInteractionDB) error {
	return r.db.Create(interaction).Error
}

// FindAllByUserID recupera todas las interacciones de un usuario específico.
func (r *progressRepository) FindAllByUserID(userID uint) ([]models.LearningInteractionDB, error) {
	var interactions []models.LearningInteractionDB
	if err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&interactions).Error; err != nil {
		return nil, err
	}
	return interactions, nil
}

func (r *progressRepository) FindByConversationID(
	userID uint,
	conversationID string,
) ([]models.LearningInteractionDB, error) {

	var interactions []models.LearningInteractionDB

	err := r.db.
		Where("user_id = ? AND conversation_id = ?", userID, conversationID).
		Order("created_at asc").
		Find(&interactions).Error

	return interactions, err
}
