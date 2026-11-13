package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizzy-design/book-corner-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func AdminLoginHandler(pool *pgxpool.Pool, jwtSecret string) gin.HandlerFunc{
	return func(c *gin.Context){
		var req loginRequest
		if err := c.ShouldBindJSON(&req);

		err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": "validation_failed",
					"message": "Valid email and password are required",
				},
			})
			return
		}

		admin, err := repository.GetAdminByEmail(pool, req.Email)
		if errors.Is(err, pgx.ErrNoRows){
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code": "invalid_credentials",
					"message": "Invalid email or password",
				},
			})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code": "internal_error",
					"message": "Failed to process login",
				},
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password)); err != nil{
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code": "invalid_credentials",
					"message": "Invalid email or password",
				},
			})
			return
		}

		// claims are the payload data baked into the token
		claims := jwt.MapClaims{
			"sub": admin.ID,
			"exp": time.Now().Add(24 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // builds unsigned token
		signedToken, err := token.SignedString([]byte(jwtSecret))  // signs the token using my secret key 
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code": "internal_error",
					"message": "Failed to generate token",
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": signedToken,
		})

	}
}