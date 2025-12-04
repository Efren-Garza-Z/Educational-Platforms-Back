package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/Efren-Garza-Z/go-api-gemini/domain/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// AuthRequired es un middleware de Gin que verifica un token JWT y extrae el UserID.
func AuthRequired() gin.HandlerFunc {
	// 1. Obtener la clave secreta del entorno
	_ = godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		// En una aplicación real, esto ya debería haberse manejado en main.go
		// para un cierre seguro, pero lo verificamos aquí como fallback.
		panic("JWT_SECRET_KEY no está configurada en el entorno.")
	}
	secretKey := []byte(jwtSecret)

	return func(c *gin.Context) {
		// 2. Extraer el token del encabezado Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Se requiere encabezado Authorization"})
			c.Abort() // Abortar procesamiento y no ir al controlador
			return
		}

		// El formato es "Bearer <token>", separamos la palabra clave.
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido. Use Bearer <token>"})
			c.Abort()
			return
		}
		tokenString := parts[1]

		// 3. Parsear y validar el token
		claims := &models.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Verificar que el método de firma sea el esperado (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("método de firma inesperado")
			}
			// Devolver la clave secreta
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			// Error: expirado, inválido, firma incorrecta, etc.
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			c.Abort()
			return
		}

		// 4. Guardar el UserID en el contexto de Gin
		// Esto permite que el controlador acceda al ID del usuario logueado.
		c.Set("userID", claims.UserID)

		// Continuar con el siguiente handler
		c.Next()
	}
}
