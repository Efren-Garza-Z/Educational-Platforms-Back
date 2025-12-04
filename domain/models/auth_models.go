package models

import "github.com/golang-jwt/jwt/v5"

// LoginInput es el payload esperado para iniciar sesión.
type LoginInput struct {
	Email    string `json:"email" binding:"required,email" example:"efren@example.com"`
	Password string `json:"password" binding:"required" example:"miPasswordSeguro123"`
}

// AuthResponse es el payload devuelto tras un login exitoso.
type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	// Opcional: Podrías incluir datos del usuario (nombre, email) aquí.
	UserID uint `json:"user_id" example:"1"`
}

// JWTClaims define los claims personalizados para nuestro token.
// Debe incluir los campos estándar de jwt.RegisteredClaims.
type JWTClaims struct {
	UserID uint `json:"user_id"` // Nuestro claim personalizado
	jwt.RegisteredClaims
}
