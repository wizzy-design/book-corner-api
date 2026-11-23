package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizzy-design/book-corner-api/internal/models"
	"github.com/wizzy-design/book-corner-api/internal/repository"
)

func CreateOrderHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context){
		var req models.CreateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": "validation_failed",
					"message": "Failed to bind into JSON",
				},
			})
			return
		}

		order, err := repository.CreateOrder(pool, req); 
		if err != nil {
			if errors.Is(err, repository.ErrBookNotFound){
				c.JSON(http.StatusNotFound, gin.H{
					"error": gin.H{
						"code": "not_found",
						"message": err.Error(),
					},
				})
				return
			}

			if errors.Is(err, repository.ErrBookNotPublished){
				c.JSON(http.StatusBadRequest, gin.H{
					"error": gin.H{
						"code": "book_unavailable",
						"message": err.Error(),
					},
				})
				return
			}

			if errors.Is(err, repository.ErrInsufficientStock) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": gin.H{
						"code": "insufficient_stock",
						"message": err.Error(),
					},
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "internal_error",
					"message": "Failed to create order",
				},
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"data": order,
		})
	}
}