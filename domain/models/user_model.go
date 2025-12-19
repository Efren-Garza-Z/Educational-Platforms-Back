package models

import "time"

// UserDB es el modelo para GORM (tabla service.users)
type UserDB struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	FullName string `json:"full_name" gorm:"not null" example:"Efren David"`
	Email    string `json:"email" gorm:"uniqueIndex;not null" example:"efren@example.com"`
	Password string `json:"password" gorm:"not null" example:"miPasswordSeguro123"`

	// CAMPOS DE PERSONALIZACIÓN PARA LA IA
	TargetLanguage string `json:"target_language" gorm:"default:'English'"`
	LanguageLevel  string `json:"language_level" gorm:"default:'A1'"`
}

func (UserDB) TableName() string {
	return "service.users"
}

// User representa la vista pública (no expone password) — usado por Swagger.
type User struct {
	ID       uint   `json:"id" example:"1"`
	FullName string `json:"full_name" example:"Efren David"`
	Email    string `json:"email" example:"efren@example.com"`
}

// CreateUserInput es el payload esperado para crear usuarios.
type CreateUserInput struct {
	FullName       string `json:"full_name" binding:"required" example:"Efren David"`
	Email          string `json:"email" binding:"required,email" example:"efren@example.com"`
	Password       string `json:"password" binding:"required" example:"miPasswordSeguro123"`
	TargetLanguage string `json:"target_language" example:"English"`
	LanguageLevel  string `json:"language_level" example:"A1"`
}

// ToPublic convierte UserDB a User (oculta password)
func (u *UserDB) ToPublic() User {
	return User{
		ID:       u.ID,
		FullName: u.FullName,
		Email:    u.Email,
	}
}
