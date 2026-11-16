package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wizzy-design/book-corner-api/internal/models"
)

type googleBooksResponse struct {
	Items []struct {
		VolumeInfo struct {
			Title               string   `json:"title"`
			Authors             []string `json:"authors"`
			Description         string   `json:"description"`
			Categories          []string `json:"categories"`
			IndustryIdentifiers []struct {
				Type       string `json:"type"`
				Identifier string `json:"identifier"`
			} `json:"industryIdentiers"`
			ImageLinks struct {
				Thumbnail string `json:"thumbnail"`
			} `json:"imageLinks"`
		} `json:"volumeInfo"`
	} `json:"items"`
}

func BookLookupHandler() gin.HandlerFunc {
	client := &http.Client{Timeout: 5 * time.Second}

	return func(c *gin.Context) {
		q := strings.TrimSpace(c.Query("q"))
		if q == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": "validation_failed",
					"message": "Query parameter 'q' is required",
				},
			})
			return
		}

		apiURL := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=%s&maxResults=5", url.QueryEscape(q))

		resp, err := client.Get(apiURL)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": gin.H{
					"code": "external_service_error",
					"message": "Failed to reach external book metadata service",
				},
			})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": gin.H{
					"code": "external_service_error",
					"message": "External book metadata service returned an error",
				},
			})
			return
		}

		var gbResp googleBooksResponse
		if err := json.NewDecoder(resp.Body).Decode(&gbResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code": "internal_error",
					"message": "Failed to parse book metadata",
				},
			})
			return
		}

		results := make([]models.BookLookupResult, 0, len(gbResp.Items))

		for _, item := range gbResp.Items {
			info := item.VolumeInfo

			// 1. Extract ISBN (prefer ISBN_13, fallback to ISBN_10)
			var isbn string
			for _, id := range info.IndustryIdentifiers {
				if id.Type == "ISBN_13" {
					isbn = id.Identifier
					break
				}
				if id.Type == "ISBN_10" && isbn == ""{
					isbn = id.Identifier
				}
			}

			// 2. High-res Open Library cover by ISBN; fallback to Google thumbnail
			imageURL := ""
			if isbn != "" {
				imageURL = fmt.Sprintf("https://covers.openlibrary.org/b/isbn/%s-L.jpg", isbn)
			} else if info.ImageLinks.Thumbnail != "" {
				imageURL = info.ImageLinks.Thumbnail
				// If imageURL starts with "http://", found is true and after holds everything after "http://".
				if after, found := strings.CutPrefix(imageURL, "http://"); found {
					imageURL = "https://" + after
				}
			}

			author := strings.Join(info.Authors, ", ")

			results = append(results, models.BookLookupResult{
				Title: info.Title,
				Author: author,
				Description: info.Description,
				ImageUrl: imageURL,
				SuggestedTags: info.Categories,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"data": results,
		})
	}
}