package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizzy-design/book-corner-api/internal/models"
)

func GetAllPublishedBooks(pool *pgxpool.Pool, filter models.BookFilter)([]models.Book, int, error){
	var ctx context.Context
	var cancel context.CancelFunc

	// Stop endpoint after 5 secs if it is faulty
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second) 
	defer cancel()

	// Total matching published books for pagination
	var countQuery string = `
		SELECT COUNT(*) -- COUNT is a built-in SQL function that counts how many rows match the WHERE filter
		FROM books b
		WHERE b.status = 'published'
		AND ($1 = '' 
			OR b.search @@ websearch_to_tsquery('english', $1) 
			OR b.title ILIKE '%' || $1 || '%' 
			OR b.author ILIKE '%' || $1 || '%'
		) -- "Does this book's search column match the user's search words?"
		AND ($2 = '' OR EXISTS (
			SELECT 1 FROM book_tags bt
			INNER JOIN tags t ON t.id = bt.tag_id
			WHERE bt.book_id = b.id AND LOWER(t.name) = LOWER($2) AND t.kind = 'genre'
		))
		AND ($3 = '' OR EXISTS(
			SELECT 1 FROM book_tags bt
			INNER JOIN tags t ON t.id = bt.tag_id
			WHERE bt.book_id = b.id AND LOWER(t.name) = LOWER($3)
		))
	`

	var total int
	var err error = pool.QueryRow(ctx, countQuery, filter.Query, filter.Genre, filter.Tag).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch paginated books, ranked by search relevance when query is present
	var offset int = (filter.Page - 1) * filter.Limit

	var query string = `
		SELECT b.id, b.title, b.author, b.description, b.price, b.stock, b.image_url, b.status, b.created_at, b.updated_at
		FROM books b
		WHERE b.status = 'published'
			AND ($1 = '' OR b.search @@ websearch_to_tsquery('english', $1) OR b.title ILIKE '%' || $1 || '%' 
			OR b.author ILIKE '%' || $1 || '%')
			AND ($2 = '' OR EXISTS(
				SELECT 1 FROM book_tags bt
				INNER JOIN tags t ON t.id = bt.tag_id
				WHERE bt.book_id = b.id AND LOWER(t.name) = LOWER($2) AND t.kind = 'genre'
			))
			AND ($3 = '' OR EXISTS(
				SELECT 1 FROM book_tags bt
				INNER JOIN tags t ON t.id = bt.tag_id
				WHERE bt.book_id = b.id AND LOWER(t.name) = lower($3)
			))
			ORDER BY CASE WHEN $1 <> '' THEN ts_rank(b.search, websearch_to_tsquery('english', $1)) ELSE 0 END DESC,
			b.created_at DESC
			LIMIT $4 OFFSET $5
	`

	var rows, queryErr = pool.Query(ctx, query, filter.Query, filter.Genre, filter.Tag, filter.Limit, offset)

	if queryErr != nil {
		return nil, 0, queryErr
	}

	defer rows.Close()

	var books []models.Book = []models.Book{}

	for rows.Next() {
		var book models.Book

		err = rows.Scan(
			&book.ID,
			&book.Title,       
			&book.Author,      
			&book.Description, 
			&book.Price,       
			&book.Stock,       
			&book.ImageUrl,    
			&book.Status,      
			&book.CreatedAt,   
			&book.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func GetBookByID (pool *pgxpool.Pool, id int) (*models.Book, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, title, author, description, price, stock, image_url, status, created_at, updated_at
		FROM books
		WHERE id = $1 AND status = $2
	`

	var book models.Book

	var err error = pool.QueryRow(ctx, query, id, "published").Scan(
		&book.ID,
		&book.Title,       
		&book.Author,      
		&book.Description, 
		&book.Price,       
		&book.Stock,       
		&book.ImageUrl,    
		&book.Status,      
		&book.CreatedAt,   
		&book.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// connects the tags table to the book_tags table by matching their tag IDs.
	var tagQuery string = `
		SELECT t.id, t.name, t.kind
		FROM tags t
		INNER JOIN book_tags bt ON t.id = bt.tag_id
		WHERE bt.book_id = $1
	`

	var rows, tagErr = pool.Query(ctx, tagQuery, id)
	if tagErr != nil {
		return nil, tagErr
	}
	defer rows.Close()

	book.Tags = []models.Tag{}

	for rows.Next(){
		var tag models.Tag
		if err = rows.Scan(&tag.ID, &tag.Name, &tag.Kind); err != nil {
			return nil, err
		}

		book.Tags = append(book.Tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &book, nil
}

func GetAllBooks(pool *pgxpool.Pool, filter models.BookFilter) ([]models.Book, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Total matching books (counting drafts & published)
	countQuery := `
		SELECT COUNT(*)
		FROM books b
		WHERE ($1 = '' OR b.status = $1)
		AND ($2 = '' 
			OR b.search @@ websearch_to_tsquery('english', $2) 
			OR b.title ILIKE '%' || $2 || '%' 
			OR b.author ILIKE '%' || $2 || '%'
		)
		AND ($3 = '' OR EXISTS (
			SELECT 1 FROM book_tags bt
			INNER JOIN tags t ON t.id = bt.tag_id
			WHERE bt.book_id = b.id AND LOWER(t.name) = LOWER($3) AND t.kind = 'genre'
		))
		AND ($4 = '' OR EXISTS (
			SELECT 1 FROM book_tags bt
			INNER JOIN tags t ON t.id = bt.tag_id
			WHERE bt.book_id = b.id AND LOWER(t.name) = LOWER($4)
		))
	`

	var total int
	err := pool.QueryRow(ctx, countQuery, filter.Status, filter.Query, filter.Genre, filter.Tag).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 2. Fetch paginated books
	offset := (filter.Page - 1) * filter.Limit

	query := `
		SELECT b.id, b.title, b.author, b.description, b.price, b.stock, b.image_url, b.status, b.created_at, b.updated_at
		FROM books b
		WHERE ($1 = '' OR b.status = $1)
			AND ($2 = '' OR b.search @@ websearch_to_tsquery('english', $2) OR b.title ILIKE '%' || $2 || '%' OR b.author ILIKE '%' || $2 || '%')
			AND ($3 = '' OR EXISTS(
				SELECT 1 FROM book_tags bt
				INNER JOIN tags t ON t.id = bt.tag_id
				WHERE bt.book_id = b.id AND LOWER(t.name) = LOWER($3) AND t.kind = 'genre'
			))
			AND ($4 = '' OR EXISTS(
				SELECT 1 FROM book_tags bt
				INNER JOIN tags t ON t.id = bt.tag_id
				WHERE bt.book_id = b.id AND LOWER(t.name) = LOWER($4)
			))
		ORDER BY CASE WHEN $2 <> '' THEN ts_rank(b.search, websearch_to_tsquery('english', $2)) ELSE 0 END DESC,
			b.created_at DESC
		LIMIT $5 OFFSET $6
	`

	rows, err := pool.Query(ctx, query, filter.Status, filter.Query, filter.Genre, filter.Tag, filter.Limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	books := []models.Book{}
	for rows.Next() {
		var book models.Book
		err = rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Description,
			&book.Price,
			&book.Stock,
			&book.ImageUrl,
			&book.Status,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return books, total, nil
}