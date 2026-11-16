package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
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

func CreateBook(pool *pgxpool.Pool, book models.Book, tagIDs []int) (*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	/* 
		transaction (pool.Begin) is used here so that inserting the book
		into books and linking its tag_ids in book_tags happen together. 
		If inserting tags fails, the book is rolled back automatically.
	*/
	transaction, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer transaction.Rollback(ctx)

	bookQuery := `
		INSERT INTO books (title, author, description, price, stock, image_url, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, title, author, description, price, stock, image_url, status, created_at, updated_at
	`

	var created models.Book
	err = transaction.QueryRow(ctx, bookQuery, 
		book.Title,
		book.Author,
		book.Description,
		book.Price,
		book.Stock,
		book.ImageUrl,
		book.Status,
	).Scan(
		&created.ID,
		&created.Title,
		&created.Author,
		&created.Description,
		&created.Price,
		&created.Stock,
		&created.ImageUrl,
		&created.Status,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// link tags in book_tags
	created.Tags = []models.Tag{}
	for _, tagID := range tagIDs {
		_, err := transaction.Exec(ctx, `
			INSERT INTO book_tags (book_id, tag_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, created.ID, tagID)
		if err != nil {
			return nil, err
		}
	}

	// 3. Fetch full tag objects to populate in the returned book
	if len(tagIDs) > 0 {
		tagRows, err := transaction.Query(ctx, `
			SELECT t.id, t.name, t.kind
			FROM tags t
			INNER JOIN book_tags bt ON bt.tag_id = t.id
			WHERE bt.book_id = $1
		`, created.ID)
		if err != nil {
			return nil, err
		}
		defer tagRows.Close()

		for tagRows.Next() {
			var tag models.Tag
			if err := tagRows.Scan(&tag.ID, &tag.Name, &tag.Kind);
			err != nil {
				return nil, err
			}
			created.Tags = append(created.Tags, tag)
		}
	}

	if err := transaction.Commit(ctx); err != nil {
		return nil, err
	}

	return &created, nil
}

func UpdateBook(pool *pgxpool.Pool, id int, book models.Book, tagIDs []int) (*models.Book, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Update book details
	query := `
		UPDATE books
		SET title = $1, author = $2, description = $3, price = $4, stock = $5, image_url = $6, status = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		RETURNING id, title, author, description, price, stock, image_url, status, created_at, updated_at
	`

	var updated models.Book
	err = tx.QueryRow(ctx, query, 
		book.Title,
		book.Author,
		book.Description,
		book.Price,
		book.Stock,
		book.ImageUrl,
		book.Status,
		id,
	).Scan(
		&updated.ID,
		&updated.Title,
		&updated.Author,
		&updated.Description,
		&updated.Price,
		&updated.Stock,
		&updated.ImageUrl,
		&updated.Status,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Replace old tags, insert new ones
	_, err = tx.Exec(ctx, `DELETE FROM book_tags WHERE book_id = $1`, id)
	if err != nil {
		return nil, err
	}

	updated.Tags = []models.Tag{}
	for _, tagID := range tagIDs {
		_, err := tx.Exec(ctx, `
			INSERT INTO books_tags (book_id, tag_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, id, tagID)
		if err != nil {
			return nil, err
		}
	}

	// Query the refreshed tags for the response
	if len(tagIDs) > 0 {
		tagRows, err := tx.Query(ctx, `
			SELECT t.id, t.name, t.kind
			FROM tags t
			INNER JOIN book_tags bt ON bt.tag_id = t.id
			WHERE bt.book_id = $1
		`, id)
		if err != nil {
			return nil, err
		}
		defer tagRows.Close()

		for tagRows.Next() {
			var tag models.Tag
			if err := tagRows.Scan(&tag.ID, &tag.Name, &tag.Kind);
			err != nil {
				return nil, err
			}
			updated.Tags = append(updated.Tags, tag)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &updated, nil
}

func DeleteBook(pool *pgxpool.Pool, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmdTag, err := pool.Exec(ctx, `DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		return err
	}

	// If no rows were affected, the book ID didn't exist
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func UpdateBookStatus(pool *pgxpool.Pool, id int, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE books
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	cmdTag, err := pool.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}