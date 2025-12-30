package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	log.Println("==========================================")
	log.Println("🔍 INICIANDO CONEXIÓN A BASE DE DATOS")
	log.Println("==========================================")
	_ = godotenv.Load()

	var dsn string
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	log.Println(host)

	if os.Getenv("K_SERVICE") != "" {
		log.Println("☁️ Modo: CLOUD RUN (Producción con Neon)")
		// NEON requiere sslmode=require
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=UTC",
			host, dbUser, dbPass, dbName, port)
	} else {
		log.Println("💻 Modo: LOCAL (Probando Neon)")
		// En local también puedes usar sslmode=require para probar la conexión real
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=UTC",
			host, dbUser, dbPass, dbName, port)
	}

	log.Println("🔌 Intentando conectar con PostgreSQL...")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("❌ ERROR al conectar: %v", err)
		return fmt.Errorf("error al conectar a la base de datos: %w", err)
	}
	log.Println("✅ Conexión GORM exitosa")

	// Verificar que la conexión funciona
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("❌ ERROR al obtener sqlDB: %v", err)
		return fmt.Errorf("error al obtener DB instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Printf("❌ ERROR en Ping: %v", err)
		return fmt.Errorf("error haciendo ping a la base de datos: %w", err)
	}
	log.Println("✅ Ping exitoso a PostgreSQL")

	// Crear schema si es necesario
	if err := DB.Exec("CREATE SCHEMA IF NOT EXISTS service;").Error; err != nil {
		log.Printf("⚠️ Aviso: No se pudo crear schema 'service': %v", err)
	} else {
		log.Println("✅ Schema 'service' verificado")
	}

	log.Println("==========================================")
	log.Println("✅ CONEXIÓN COMPLETADA EXITOSAMENTE")
	log.Println("==========================================")
	return nil
}
