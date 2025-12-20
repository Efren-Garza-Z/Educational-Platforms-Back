package models

import "time"

// LearningInteractionDB es el modelo para GORM (tabla service.learning_interactions)
type LearningInteractionDB struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	ConversationID string `gorm:"index"` // 👈 NUEVO

	UserID          uint   `json:"user_id" gorm:"not null"`                               // Clave Foránea al usuario
	InteractionType string `json:"interaction_type" gorm:"not null" example:"Correction"` // Ejemplo: Conversation, Grammar, Exercise
	Language        string `json:"language" gorm:"not null" example:"French"`
	Level           string `json:"level" gorm:"not null" example:"B2"`
	Prompt          string `json:"prompt" gorm:"type:text" example:"Write a dialogue about a train ticket."`
	Response        string `json:"response" gorm:"type:text" example:"Bonjour, je voudrais acheter un billet."`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (LearningInteractionDB) TableName() string {
	return "service.learning_interactions"
}

// LearningInteractionInput es el payload para guardar una nueva interacción.
type LearningInteractionInput struct {
	UserID          uint   `json:"user_id" binding:"required"`
	ConversationID  string `json:"conversation_id" binding:"required"`
	InteractionType string `json:"interaction_type" binding:"required" example:"Conversation"`
	Language        string `json:"language" binding:"required" example:"French"`
	Level           string `json:"level" binding:"required" example:"B2"`
	Prompt          string `json:"prompt" binding:"required"`
	Response        string `json:"response" binding:"required"`
}
