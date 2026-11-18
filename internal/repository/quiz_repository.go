package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizzy-design/book-corner-api/internal/models"
)

func GetQuizQuestions(pool *pgxpool.Pool) ([]models.QuizQuestion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	questionQuery := `
		SELECT id, prompt, position
		FROM quiz_questions
		ORDER BY position ASC
	`

	// Fetch all questions ordered by position
	rows, err := pool.Query(ctx, questionQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questions := []models.QuizQuestion{}
	questionIndexMap := make(map[int]int)

	for rows.Next() {
		var q models.QuizQuestion
		if err := rows.Scan(&q.ID, &q.Prompt, &q.Position);
		err != nil {
			return nil, err
		}

		// Initialize Options to an empty slice so JSON outputs [] instead of null
		q.Options = []models.QuizOption{}
		questionIndexMap[q.ID] = len(questions)
		questions = append(questions, q)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	optionQuery := `
		SELECT id, question_id, label
		FROM quiz_options
		ORDER BY id ASC
	`

	optRows, err := pool.Query(ctx, optionQuery)
	if err != nil {
		return nil, err
	}
	defer optRows.Close()

	for optRows.Next(){
		var opt models.QuizOption
		if err := optRows.Scan(&opt.ID, &opt.QuestionID, &opt.Label); err != nil {
			return nil, err
		}

		if idx, ok := questionIndexMap[opt.QuestionID]; ok {
			questions[idx].Options = append(questions[idx].Options, opt)
		}
	}

	if err := optRows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}

func GetOptionTagWeights(pool *pgxpool.Pool, optionIDs []int) (map[int]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT tag_id, SUM(weight)
		FROM option_tag_weights
		WHERE option_id = ANY($1)
		GROUP BY tag_id
	`

	rows, err := pool.Query(ctx, query, optionIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	weights := make(map[int]int)
	for rows.Next() {
		var tagID, totalWeight int
		if err := rows.Scan(&tagID, &totalWeight); err != nil {
			return nil, err
		}
		weights[tagID] = totalWeight
	}

	return weights, rows.Err()
}

func GetPublishedBooksForScoring(pool *pgxpool.Pool) ([]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	booksQuery := `
		SELECT id, title, author, description, price, stock, image_url, status, created_at, updated_at
		FROM books
		WHERE status = 'published'
	`

	// Fetch published books
	rows, err := pool.Query(ctx, booksQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	bookIndexMap := make(map[int]int)

	for rows.Next() {
		var b models.Book
		err := rows.Scan(
			&b.ID, &b.Title, &b.Author, &b.Description,
			&b.Price, &b.Stock, &b.ImageUrl, &b.Status,
			&b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		b.Tags = []models.Tag{}
		bookIndexMap[b.ID] = len(books)
		books = append(books, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Fetch tags for published books
	tagQuery := `
		SELECT bt.book_id, t.id, t.name
		FROM book_tags bt
		INNER JOIN tags t ON t.id = bt.tag_id
		INNER JOIN books b ON b.id = bt.book_id
		WHERE b.status = 'published'
	`
	tagRows, err := pool.Query(ctx, tagQuery)
	if err != nil {
		return nil, err
	}
	defer tagRows.Close()

	for tagRows.Next() {
		var bookID int
		var tag models.Tag
		if err := tagRows.Scan(&bookID, &tag.ID, &tag.Name); err != nil {
			return nil, err
		}

		if idx, ok := bookIndexMap[bookID]; ok {
			books[idx].Tags = append(books[idx].Tags, tag)
		}
	}

	if err := tagRows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}