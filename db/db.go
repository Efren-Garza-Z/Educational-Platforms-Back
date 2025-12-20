package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var dsn string

	// K_SERVICE es una variable automática de Cloud Run
	if os.Getenv("K_SERVICE") != "" {
		// CONFIGURACIÓN PARA PRODUCCIÓN (GCP - Unix Sockets)
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		instanceConnectionName := os.Getenv("INSTANCE_CONNECTION_NAME")

		// Agregamos TimeZone=UTC para consistencia total
		dsn = fmt.Sprintf("user=%s password=%s dbname=%s host=/cloudsql/%s sslmode=disable TimeZone=UTC",
			dbUser, dbPass, dbName, instanceConnectionName)
	} else {
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
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	// Aseguramos el schema para tus modelos personalizados
	if err := DB.Exec("CREATE SCHEMA IF NOT EXISTS service;").Error; err != nil {
		log.Printf("Aviso: No se pudo crear el schema 'service' (puede que ya exista o no tengas permisos): %v", err)
	}

	log.Println("Conexión a la base de datos exitosa")
}
