package services

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/Efren-Garza-Z/go-api-gemini/domain/models"
	"github.com/Efren-Garza-Z/go-api-gemini/domain/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(input models.CreateUserInput) (*models.UserDB, error)
	GetAllUsers() ([]models.UserDB, error)
	GetUserByID(id uint) (*models.UserDB, error)
	GetUser(email string) (*models.UserDB, error)
	FindUserByEmail(email string) (*models.UserDB, error)
	UpdateUser(email string, input models.CreateUserInput) (*models.UserDB, error)
	DeleteUser(id uint) error
	Login(email, password string) (*models.UserDB, error)
	GenerateJWT(user *models.UserDB) (string, error)
	UpdateLanguage(email string, input models.UpdateLanguageInput) (*models.UserDB, error)
}

type userService struct {
	repo      repositories.UserRepository
	jwtSecret string
}

func NewUserService(r repositories.UserRepository) UserService {

	_ = godotenv.Load()
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		// MUY IMPORTANTE: Terminar la aplicación si la clave no existe
		log.Fatal("FATAL: JWT_SECRET_KEY no está configurada en el entorno.")
	}

	return &userService{
		repo:      r,
		jwtSecret: secret,
	}
}

func (s *userService) CreateUser(input models.CreateUserInput) (*models.UserDB, error) {
	hashedPassword, err := s.hashPassword(input.Password)
	if err != nil {
		return nil, err // Error al hashear
	}
	user := &models.UserDB{
		FullName:       input.FullName,
		Email:          input.Email,
		Password:       hashedPassword, // ideal: hash aquí
		TargetLanguage: input.TargetLanguage,
		LanguageLevel:  input.LanguageLevel,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetAllUsers() ([]models.UserDB, error) {
	return s.repo.FindAll()
}

func (s *userService) GetUserByID(id uint) (*models.UserDB, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("usuario no encontrado")
	}
	return u, nil
}

func (s *userService) GetUser(email string) (*models.UserDB, error) {
	u, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("usuario no encontrado")
	}
	return u, nil
}

func (s *userService) UpdateUser(email string, input models.CreateUserInput) (*models.UserDB, error) {
	u, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("usuario no encontrado")
	}
	u.FullName = input.FullName
	u.Email = input.Email
	u.LanguageLevel = input.LanguageLevel
	u.TargetLanguage = input.TargetLanguage

	if input.Password != "" {
		hashedPassword, err := s.hashPassword(input.Password)
		if err != nil {
			return nil, err // Error al hashear
		}
		u.Password = hashedPassword
	}
	// opcional: actualizar password si se envía
	if err := s.repo.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *userService) FindUserByEmail(email string) (*models.UserDB, error) {
	return s.repo.FindUserByEmail(email)
}

// Implementación del Login
func (s *userService) Login(email, password string) (*models.UserDB, error) {
	// 1. Buscar el usuario por email (Necesitarás un método en el repositorio)
	user, err := s.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("credenciales inválidas")
	}

	// 2. Comparar la contraseña de texto plano con el hash guardado
	if !s.checkPasswordHash(password, user.Password) {
		return nil, errors.New("credenciales inválidas")
	}

	// 3. Éxito
	return user, nil
}

func (s *userService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}

// hashPassword toma una contraseña en texto plano y devuelve su hash.
func (s *userService) hashPassword(password string) (string, error) {
	// Generar hash con costo 14 (o el que desees)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("error al hashear la contraseña")
	}
	return string(hashedBytes), nil
}

// checkPasswordHash compara una contraseña de texto plano con un hash existente.
func (s *userService) checkPasswordHash(password, hash string) bool {
	// Si la comparación falla, devuelve un error
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT crea un token JWT firmado para el usuario dado.
func (s *userService) GenerateJWT(user *models.UserDB) (string, error) {
	// Definir el tiempo de expiración (ej. 24 horas)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Crear las claims (cargas útiles)
	claims := &models.JWTClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			// Emitido en: ahora
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// Expira en: 24 horas a partir de ahora
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declarar el token con el algoritmo de firma (HS256) y las claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token usando la clave secreta
	tokenString, err := token.SignedString([]byte(s.jwtSecret))

	if err != nil {
		return "", errors.New("error al firmar el token JWT")
	}

	return tokenString, nil
}

func (s *userService) UpdateLanguage(email string, input models.UpdateLanguageInput) (*models.UserDB, error) {
	u, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("usuario no encontrado")
	}

	// Actualizamos solo los campos específicos
	u.TargetLanguage = input.TargetLanguage
	u.LanguageLevel = input.LanguageLevel

	if err := s.repo.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}
