package quiz

import (
	"sort"

	"github.com/wizzy-design/book-corner-api/internal/models"
)

type ScoredBook struct {
	Book models.Book `json:"book"`
	Score int `json:"score"`
}

// Score takes a list of books and a map of tag_id -> weight,
// calculates scores, sorts them (score -> stock -> newest), and returns up to limit results.
func Score(books []models.Book, tagWeights map[int]int, limit int) []ScoredBook {
	var scored []ScoredBook

	for _, book := range books {
		score := 0
		for _, tag := range book.Tags {
			if weight, exists := tagWeights[tag.ID]; exists{
				score += weight
			}
		}

		// Only recommend books that matched at least one tag preference
		if score > 0 {
			scored = append(scored, ScoredBook{
				Book: book,
				Score: score,
			})
		}
	}

	// Sort: 1. Score descending, 2. Stock descending, 3. CreatedAt descending
	sort.Slice(scored, func(i, j int) bool {
		// Primary sort: Score
		if scored[i].Score != scored[j].Score {
			return scored[i].Score > scored[j].Score
		}

		// Tie-break (In-stock first)
		iInStock := scored[i].Book.Stock > 0
		jInStock := scored[j].Book.Stock > 0
		if iInStock != jInStock {
			return iInStock // true comes before false
		}

		// Tie-break (Newest first)
		return scored[i].Book.CreatedAt.After(scored[j].Book.CreatedAt)
	})

	if limit > 0 && len(scored) > limit {
		return scored[:limit]
	}

	return scored
}