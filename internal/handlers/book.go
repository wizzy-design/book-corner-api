package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizzy-design/book-corner-api/internal/models"
	"github.com/wizzy-design/book-corner-api/internal/repository"
)

func GetBooksHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context){
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}

		limit, err := strconv.Atoi(c.DefaultQuery("limit", "12"))
		if err != nil || limit < 1 {
			limit = 12
		}
		if limit > 50 {
			limit = 50
		}

		filter := models.BookFilter{
			Query: c.Query("query"),
			Genre: c.Query("genre"),
			Tag: c.Query("tag"),
			Page: page,
			Limit: limit,
		}

		books, total, err := repository.GetAllPublishedBooks(pool, filter)

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
			"page": filter.Page,
			"limit": filter.Limit,
			"total": total,
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

func GetAdminBooksHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}

		limit, err := strconv.Atoi(c.DefaultQuery("limit", "12"))
		if err != nil || limit < 1 {
			limit = 12
		}
		if limit > 50 {
			limit = 50
		}

		filter := models.BookFilter{
			Query:  c.Query("query"),
			Genre:  c.Query("genre"),
			Tag:    c.Query("tag"),
			Status: c.Query("status"),
			Page:   page,
			Limit:  limit,
		}

		books, total, err := repository.GetAllBooks(pool, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "internal_error",
					"message": "Failed to fetch books",
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  books,
			"page":  filter.Page,
			"limit": filter.Limit,
			"total": total,
		})
	}
}