package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizzy-design/book-corner-api/internal/config"
	"github.com/wizzy-design/book-corner-api/internal/database"
	"github.com/wizzy-design/book-corner-api/internal/handlers"
	"github.com/wizzy-design/book-corner-api/internal/middleware"
	"github.com/wizzy-design/book-corner-api/migrations"
)

func main() {
	var cfg *config.Config
	var err error
	cfg, err = config.Load()

	// Close program if no env variables exist
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Setup database connection pool
	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)
	
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Cleans up pool connection when main function exits
	defer pool.Close()

	// Execute migrations when server starts up
	if err := database.Migrate(pool, migrations.FS); err != nil {
		log.Fatal("Failed to run database migrations:", err)
	}

	// Start http server
	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)

	// Middleware
	router.Use(middleware.CORSMiddleware(cfg.AllowedOrigin))

	// Public routes
	router.GET("/health", handlers.HealthCheck(pool))
	router.GET("/books", handlers.GetBooksHandler(pool))
	router.GET("/books/:id", handlers.GetBookByIDHandler(pool))

	// Protected admin routes
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Checkpoint route: verifies auth works
		admin.GET("/books", handlers.GetAdminBooksHandler(pool))
	}
	router.POST("/admin/login", handlers.AdminLoginHandler(pool, cfg.JWTSecret))

	// Tell router to listen on the port:
	router.Run(":" + cfg.Port)
}