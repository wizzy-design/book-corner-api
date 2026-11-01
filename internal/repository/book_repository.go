package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizzy-design/book-corner-api/internal/models"
)

func GetAllPublishedBooks(pool *pgxpool.Pool)([]models.Book, error){
	var ctx context.Context
	var cancel context.CancelFunc

	// Stop endpoint after 5 secs if it is faulty
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second) 
	defer cancel()

	var query string = `
		SELECT id, title, author, description, price, stock, image_url, status, created_at, updated_at
		FROM books
		WHERE status = $1
	`

	var rows, err = pool.Query(ctx, query, "published")

	if err != nil {
		return nil, err
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
			return nil, err
		}

		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
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