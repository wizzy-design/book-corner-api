package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizzy-design/book-corner-api/internal/models"
)

var (
	ErrBookNotFound = errors.New("book not found")
	ErrBookNotPublished = errors.New("book is not available for purchase")
	ErrInsufficientStock = errors.New("insufficient stock")
)

func CreateOrder(pool *pgxpool.Pool, req models.CreateOrderRequest) (models.Order, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel();

	tx, err := pool.Begin(ctx)
	if err != nil {
		return models.Order{}, err
	}
	defer tx.Rollback(ctx)

	/* 
		FOR UPDATE locks that specific book row in Postgres. If two customers try to buy the
		last copy at the same millisecond, the second customer's transaction waits until the first transaction finishes. This physically guarantees no negative stock!
	*/
	query := `
		SELECT title, price, stock, status
		FROM books WHERE id = $1 FOR UPDATE 
	`

	var total int
	var orderItems []models.OrderItem

	// verify books, deducted stock, calc server total
	for _, item := range req.Items {
		var title, status string
		var price, stock int

		err = tx.QueryRow(ctx, query, item.BookID).Scan(
			&title,
			&price,
			&stock,
			&status,
		)
		if errors.Is(err, pgx.ErrNoRows){
			return models.Order{}, fmt.Errorf("%w: book id %d", ErrBookNotFound, item.BookID)
		}
		if err != nil {
			return models.Order{}, err
		}

		// Can only buy published books
		if status != "published" {
			return models.Order{}, fmt.Errorf("%w: '%s'", ErrBookNotPublished, title)
		}

		// Can only buy books in stock
		if item.Quantity > stock {
			return models.Order{}, fmt.Errorf("%w: '%s' (only %d left)", ErrInsufficientStock, title, stock)
		}

		// Deduct stock from the books table
		updateStockQuery := `
			UPDATE books
			SET stock = stock - $1, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2
		`

		_, err = tx.Exec(ctx, updateStockQuery, item.Quantity, item.BookID)
		if err != nil {
			return models.Order{}, fmt.Errorf("failed to update stock: %w", err) 
		}

		// Accumulate total on server
		total += price * item.Quantity

		// Stash the item data for inserting into order_items
		orderItems = append(orderItems, models.OrderItem{
			BookID: item.BookID,
			BookTitle: title,
			Quantity: item.Quantity,
			UnitPrice: price,
		})
	}

	insertOrderQuery := `
		INSERT INTO orders (
			customer_name, phone, email, address, note, status, payment_status, payment_method, total
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	var created models.Order

	created = models.Order{
		CustomerName:  req.CustomerName,
		Phone:         req.Phone,
		Email:         req.Email,
		Address:       req.Address,
		Note:          req.Note,
		Status:        models.OrderStatusPending,
		PaymentStatus: models.PaymentStatusPending,
		PaymentMethod: req.PaymentMethod,
		Total:         total,
	}

	err = tx.QueryRow(ctx, insertOrderQuery, 
		created.CustomerName,
		created.Phone,
		created.Email,
		created.Address,
		created.Note,
		created.Status,
		created.PaymentStatus,
		created.PaymentMethod,
		created.Total,
	).Scan(&created.ID, &created.CreatedAt)

	if err != nil {
		return models.Order{}, fmt.Errorf("failed to create order: %w", err)
	}

	insertItemQuery := `
		INSERT INTO order_items (order_id, book_id, quantity, unit_price)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	for i := range orderItems {
		orderItems[i].OrderID = created.ID
		err = tx.QueryRow(ctx, insertItemQuery, 
			created.ID,
			orderItems[i].BookID,
			orderItems[i].Quantity,
			orderItems[i].UnitPrice,
		).Scan(&orderItems[i].ID)

		if err != nil {
			return models.Order{}, fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	created.Items = orderItems

	if err = tx.Commit(ctx); err != nil {
		return models.Order{}, fmt.Errorf("failed to commit order transaction: %w", err)
	}

	return created, nil
}

func SetOrderPaymentReference(pool *pgxpool.Pool, orderID int, reference string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE orders
		SET payment_reference = $1
		WHERE id = $2
	`

	_, err := pool.Exec(ctx, query, reference, orderID)
	return err
}