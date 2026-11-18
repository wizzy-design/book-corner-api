package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizzy-design/book-corner-api/internal/models"
	"github.com/wizzy-design/book-corner-api/internal/quiz"
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

func RecommendBooksHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context){
		var req models.QuizRecommendRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": "validation_failed",
					"message": "Invalid input from user",
				},
			})
			return
		}

		// Collect option IDs from answers
		var optionIDs []int
		for _, ans := range req.Answers {
			optionIDs = append(optionIDs, ans.OptionID)
		}

		// Fetch tag weights for chosen options
		weights, err := repository.GetOptionTagWeights(pool, optionIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code": "internal_error",
					"message": "Failed to calculate option weights",
				},
			})
			return
		}

		// Fetch all published books with their tags
		books, err := repository.GetPublishedBooksForScoring(pool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "internal_error",
					"message": "Failed to fetch books for recommendation",
				},
			})
			return
		}

		// Score and rank top 5 books
		recommendations := quiz.Score(books, weights, 5)
		if recommendations == nil {
			recommendations = []quiz.ScoredBook{}
		}

		c.JSON(http.StatusOK, gin.H{
			"data": recommendations,
		})
	}
}