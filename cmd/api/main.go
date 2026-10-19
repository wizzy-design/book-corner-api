package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizzy-design/book-corner-api/internal/config"
	"github.com/wizzy-design/book-corner-api/internal/database"
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
	router.GET("/", func(c *gin.Context){
		// A context is Go's automatic cancel-switch. 
		// It cancels database queries immediately if a client disconnects or times out. 
		if err := pool.Ping(c.Request.Context());
		err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "database": "disconnected"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "healthy", "database": "connected"})
	})

	// Tell router to listen on the port:
	router.Run(":" + cfg.Port)
}