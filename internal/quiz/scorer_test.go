package quiz_test

import (
	"testing"
	"time"

	"github.com/wizzy-design/book-corner-api/internal/models"
	"github.com/wizzy-design/book-corner-api/internal/quiz"
)

func TestScore(t *testing.T) {
	now := time.Now()
	older := now.Add(-24 * time.Hour)

	// Sample test books
	bookHighMatch := models.Book{
		ID: 1,
		Title: "High Match Thriller",
		Tags: []models.Tag{
			{ID: 10, Name: "Mystery"},
			{ID: 20, Name: "Thriller"},
		},
		Stock: 5,
		CreatedAt: now,
	}

	bookLowMatch := models.Book{
		ID:    2,
		Title: "Low Match Mystery",
		Tags: []models.Tag{
			{ID: 10, Name: "Mystery"},
		},
		Stock:     3,
		CreatedAt: now,
	}

	bookZeroMatch := models.Book{
		ID:    3,
		Title: "Baking Bread",
		Tags: []models.Tag{
			{ID: 99, Name: "Cooking"},
		},
		Stock:     10,
		CreatedAt: now,
	}

	bookTiedOutOfStock := models.Book{
		ID:    4,
		Title: "Tied But Out of Stock",
		Tags: []models.Tag{
			{ID: 10, Name: "Mystery"},
		},
		Stock:     0, // Out of stock
		CreatedAt: now,
	}

	bookTiedOlder := models.Book{
		ID:    5,
		Title: "Tied But Older",
		Tags: []models.Tag{
			{ID: 10, Name: "Mystery"},
		},
		Stock:     3,
		CreatedAt: older, // Older creation date
	}

	weights := map[int]int{
		10: 5, // Mystery: 5 pts
		20: 4, // Thriller: 4 pts
	}

	tests := []struct {
		name        string
		books       []models.Book
		weights     map[int]int
		limit       int
		expectedIDs []int
	}{
		{
			name:        "highest score comes first and zero matches are excluded",
			books:       []models.Book{bookLowMatch, bookZeroMatch, bookHighMatch},
			weights:     weights,
			limit:       10,
			expectedIDs: []int{1, 2}, // 1 has 9 pts, 2 has 5 pts, 3 has 0 pts (excluded)
		},
		{
			name:        "tie-breaker prefers in-stock over out-of-stock",
			books:       []models.Book{bookTiedOutOfStock, bookLowMatch},
			weights:     weights,
			limit:       10,
			expectedIDs: []int{2, 4}, // Both have 5 pts, but 2 has stock > 0, 4 has stock == 0
		},
		{
			name:        "tie-breaker prefers newer book when scores and stock are equal",
			books:       []models.Book{bookTiedOlder, bookLowMatch},
			weights:     weights,
			limit:       10,
			expectedIDs: []int{2, 5}, // Both have 5 pts and stock 3, but 2 is newer than 5
		},
		{
			name:        "limit caps the number of results returned",
			books:       []models.Book{bookHighMatch, bookLowMatch, bookTiedOlder},
			weights:     weights,
			limit:       1,
			expectedIDs: []int{1}, // Only top 1 returned
		},
		{
			name:        "empty weights returns no results",
			books:       []models.Book{bookHighMatch, bookLowMatch},
			weights:     map[int]int{},
			limit:       10,
			expectedIDs: []int{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			results := quiz.Score(tc.books, tc.weights, tc.limit)

			if len(results) != len(tc.expectedIDs) {
				t.Fatalf("expected %d results, got %d", len(tc.expectedIDs), len(results))
			}

			for i, expectedID := range tc.expectedIDs {
				if results[i].Book.ID != expectedID {
					t.Errorf("at index %d: expected book ID %d, got %d", i, expectedID, results[i].Book.ID)
				}
			}
		})
	}
}