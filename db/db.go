package db

import (
	"fmt"
	"log"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	log.Println("==========================================")
	log.Println("🔍 INICIANDO CONEXIÓN A BASE DE DATOS")
	log.Println("==========================================")
	
	var dsn string
	
	// K_SERVICE es una variable automática de Cloud Run
	if os.Getenv("K_SERVICE") != "" {
		log.Println("☁️ Modo: CLOUD RUN (Producción)")
		
		// CONFIGURACIÓN PARA PRODUCCIÓN (GCP - Unix Sockets)
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		instanceConnectionName := os.Getenv("INSTANCE_CONNECTION_NAME")
		
		// Log de verificación (sin mostrar password completa)
		log.Printf("📋 DB_USER: %s", dbUser)
		log.Printf("📋 DB_NAME: %s", dbName)
		log.Printf("📋 INSTANCE_CONNECTION_NAME: %s", instanceConnectionName)
		log.Printf("📋 K_SERVICE: %s", os.Getenv("K_SERVICE"))
		
		// Validación de variables críticas
		if dbUser == "" {
			return fmt.Errorf("❌ DB_USER no está configurado")
		}
		if dbPass == "" {
			return fmt.Errorf("❌ DB_PASSWORD no está configurado")
		}
		if dbName == "" {
			return fmt.Errorf("❌ DB_NAME no está configurado")
		}
		if instanceConnectionName == "" {
			return fmt.Errorf("❌ INSTANCE_CONNECTION_NAME no está configurado")
		}
		
		dsn = fmt.Sprintf("user=%s password=%s dbname=%s host=/cloudsql/%s sslmode=disable TimeZone=UTC",
			dbUser, dbPass, dbName, instanceConnectionName)
		
		log.Printf("🔗 DSN Cloud Run: user=%s dbname=%s host=/cloudsql/%s", dbUser, dbName, instanceConnectionName)
		
	} else {
		log.Println("💻 Modo: LOCAL (Desarrollo)")
		
		// CONFIGURACIÓN PARA LOCAL
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}
		user := os.Getenv("DB_USER")
		if user == "" {
			user = "edgz"
		}
		password := os.Getenv("DB_PASSWORD")
		if password == "" {
			password = "1234"
		}
		dbname := os.Getenv("DB_NAME")
		if dbname == "" {
			dbname = "edgz"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}
		
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			host, user, password, dbname, port)
		
		log.Printf("🔗 DSN Local: host=%s user=%s dbname=%s port=%s", host, user, dbname, port)
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
