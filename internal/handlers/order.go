package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizzy-design/book-corner-api/internal/models"
	"github.com/wizzy-design/book-corner-api/internal/paystack"
	"github.com/wizzy-design/book-corner-api/internal/repository"
)

func CreateOrderHandler(pool *pgxpool.Pool, paystackClient *paystack.Client) gin.HandlerFunc {
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

		if req.PaymentMethod == models.PaymentMethodPaystack &&
		req.Email == nil || *req.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": "validation_failed",
					"message": "Email is required for Paystack payments",
				},
			})
		}

		// Run database transaction (deduct stock, insert order)
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

		// If payment is Paystack, initialize checkout session
		if req.PaymentMethod == models.PaymentMethodPaystack {
			ref := fmt.Sprintf("bc_%d_%d", order.ID, time.Now().Unix())

			paystackResp, err := paystackClient.InitializeTransaction(c.Request.Context(), *req.Email, order.Total, ref)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{
						"code": "payment_initialization_failed",
						"message": err.Error(),
					},
				})
				return
			}

			// Persist the reference on the order in the database
			if err := repository.SetOrderPaymentReference(pool, order.ID, ref); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{
						"code": "internal_error",
						"message": "Failed to record payment reference",
					},
				})
				return
			}

			order.PaymentReference = &ref

			c.JSON(http.StatusCreated, gin.H{
				"data": gin.H{
					"order": order,
					"checkout_url": paystackResp.Data.AuthorizationURL,
				},
			})
		}

		c.JSON(http.StatusCreated, gin.H{
			"data": order,
		})
	}
}