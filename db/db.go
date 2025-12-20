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

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if os.Getenv("K_SERVICE") != "" {
		// CLOUD RUN (Cloud SQL vía TCP)
		dsn = fmt.Sprintf(
			"host=127.0.0.1 user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=UTC",
			dbUser, dbPass, dbName,
		)
	} else {
		// LOCAL
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}

		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			host, dbUser, dbPass, dbName, port,
		)
	}

	log.Println("Conectando a DB...")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error DB: %v", err)
	}

	log.Println("DB conectada correctamente")
}
