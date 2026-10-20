package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizzy-design/book-corner-api/internal/config"
	"github.com/wizzy-design/book-corner-api/internal/database"
	"github.com/wizzy-design/book-corner-api/internal/handlers"
)

func main() {
	var cfg *config.Config
	var err error
	cfg, err = config.Load()

	// Close program if no env variables exist
	if err != nil {
		log.Fatal("Failed to load configuratyion:", err)
	}

	// Setup database connection pool
	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)
	
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Cleans up pool connection when main function exits
	defer pool.Close()

	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)

	// Smoke test
	router.GET("/health", handlers.HealthCheck(pool))

	// Tell router to listen on the port:
	router.Run(":" + cfg.Port)
}