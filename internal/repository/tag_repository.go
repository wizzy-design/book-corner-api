package repository

import (
	"context"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizzy-design/book-corner-api/internal/models"
)

func GetAllTags(pool *pgxpool.Pool) ([]models.Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := pool.Query(ctx, `SELECT id, name FROM tags ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []models.Tag{}
	for rows.Next() {
		var tag models.Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

func CreateTag(pool *pgxpool.Pool, name string) (*models.Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tag models.Tag
	err := pool.QueryRow(ctx, `
		INSERT INTO tags (name)
		VALUES ($1)
		RETURNING id, name
	`, name).Scan(&tag.ID, &tag.Name)

	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func ReplaceBookTags(pool *pgxpool.Pool, bookID int, tagIDs []int) ([]models.Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Delete existing associations
	_, err = tx.Exec(ctx, `DELETE FROM book_tags WHERE book_id = $1`, bookID)
	if err != nil {
		return nil, err
	}

	// Insert new associations
	for _, tagID := range tagIDs {
		_, err := tx.Exec(ctx, `
			INSERT INTO book_tags (book_id, tag_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, bookID, tagID)
		if err != nil {
			return nil, err	
		}
	}

	// Return refreshed tags
	tags := []models.Tag{}
	if len(tagIDs) > 0 {
		rows, err := tx.Query(ctx, `
			SELECT t.id, t.name
			FROM tags t
			INNER JOIN book_tags bt ON bt.tag_id = t.id
			WHERE bt.book_id = $1
			ORDER BY t.name
		`, bookID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var tag models.Tag
			if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
				return nil, err
			}
			tags = append(tags, tag)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return tags, nil
}