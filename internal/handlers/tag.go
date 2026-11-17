package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizzy-design/book-corner-api/internal/models"
	"github.com/wizzy-design/book-corner-api/internal/repository"
)

func GetTagsHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context){
		tags, err := repository.GetAllTags(pool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code": "internal_error",
					"message": "Failed to fetch tags",
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": tags,
		})
	}
}

func CreateTagHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateTagRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": "validation_failed",
					"message": "Tag name is required",
				},
			})
			return
		}

		tag, err := repository.CreateTag(pool, req.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code": "internal_error",
					"message": "Failed to create tag (tag may already exist)",
				},
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"data": tag,
		})
	}
}

func ReplaceBookTagsHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		bookID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": "invalid_id",
					"message": "Invalid book ID",
				},
			})
			return
		}

		var req models.ReplaceBookTagsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": "validation_failed",
					"message": "'tag_ids' array is required",
				},
			})
			return
		}

		tags, err := repository.ReplaceBookTags(pool, bookID, req.TagIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code": "internal_error",
					"message": "Failed to update book tags",
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": tags,
		})
	}
}