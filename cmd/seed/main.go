package main

import (
	"flag"
	"log"

	"github.com/wizzy-design/book-corner-api/internal/config"
	"github.com/wizzy-design/book-corner-api/internal/database"
	"github.com/wizzy-design/book-corner-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	email := flag.String("email", "", "Admin email address")
	password := flag.String("password", "", "Admin password")
	flag.Parse()

	if *email == "" || *password == "" {
		log.Fatal("Usage: go run ./cmd/seed -email=admin@example.com -password=yourpassword")
	}

	// Read .env
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	// Initiate database connection pool
	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pool.Close()

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Use repository to create admin
	admin, err := repository.CreateAdmin(pool, *email, string(hash))
	if err != nil {
		log.Fatal("Failed to create admin: ", err)
	}

	log.Printf("Admin seeded successfully: id=%d email=%s", admin.ID, admin.Email)
}