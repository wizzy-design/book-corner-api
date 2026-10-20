package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Tip: handlers process incoming web requests and send back responses

func HealthCheck(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context){
		if err := pool.Ping(c.Request.Context());
		err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "database": "disconnected"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "healthy", "database": "connected"})
	}
}