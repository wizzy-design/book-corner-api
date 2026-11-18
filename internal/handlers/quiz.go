package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizzy-design/book-corner-api/internal/repository"
)

func GetQuizQuestionsHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context){
		questions, err := repository.GetQuizQuestions(pool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code": "internal_error",
					"message": "An error occurred while fetching questions",
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": questions, 
		})
	}
}