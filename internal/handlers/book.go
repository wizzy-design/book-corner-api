package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizzy-design/book-corner-api/internal/repository"
)

func GetBooksHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context){
		books, err := repository.GetAllPublishedBooks(pool)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code": "internal_error",
					"message": "Failed to fetch books",
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": books,
		})
	}
}

func GetBookByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context){
		idStr := c.Param("id")

		// Converts id to an int
		id, err := strconv.Atoi(idStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "invalid_id",
					"message": "Invalid book ID",
				},
			})
			return
		}

		book, err := repository.GetBookByID(pool, id)

		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "not_found",
					"message": "Book not found",
				},
			})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "internal_error",
					"message": "Failed to fetch book",
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": book,
		})
	}
}