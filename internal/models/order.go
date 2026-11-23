package models

import "time"

// Order and Payment Constants
const (
	OrderStatusPending        = "pending"
	OrderStatusConfirmed      = "confirmed"
	OrderStatusOutForDelivery = "out_for_delivery"
	OrderStatusDelivered      = "delivered"
	OrderStatusCancelled      = "cancelled"

	PaymentStatusPending = "pending"
	PaymentStatusPaid    = "paid"

	PaymentMethodPayOnDelivery = "pay_on_delivery"
	PaymentMethodPaystack      = "paystack"
)

// OrderItem represents a line item in an order
type OrderItem struct {
	ID        int    `json:"id" db:"id"`
	OrderID   int    `json:"order_id" db:"order_id"`
	BookID    int    `json:"book_id" db:"book_id"`
	BookTitle string `json:"book_title,omitempty" db:"book_title"` // Helpful when returning order details
	Quantity  int    `json:"quantity" db:"quantity"`
	UnitPrice int    `json:"unit_price" db:"unit_price"`
}

// Order represents an order in the database
type Order struct {
	ID               int         `json:"id" db:"id"`
	CustomerName     string      `json:"customer_name" db:"customer_name"`
	Phone            string      `json:"phone" db:"phone"`
	Email            *string     `json:"email,omitempty" db:"email"`
	Address          string      `json:"address" db:"address"`
	Note             *string     `json:"note,omitempty" db:"note"`
	Status           string      `json:"status" db:"status"`
	PaymentStatus    string      `json:"payment_status" db:"payment_status"`
	PaymentMethod    string      `json:"payment_method" db:"payment_method"`
	PaymentReference *string     `json:"payment_reference,omitempty" db:"payment_reference"`
	Total            int         `json:"total" db:"total"`
	CreatedAt        time.Time   `json:"created_at" db:"created_at"`
	Items            []OrderItem `json:"items,omitempty"`
}

// DTOs (Data Transfer Objects) for incoming HTTP requests

type CreateOrderItemRequest struct {
	BookID   int `json:"book_id" binding:"required,gt=0"`
	Quantity int `json:"quantity" binding:"required,gt=0"`
}

type CreateOrderRequest struct {
	CustomerName  string                   `json:"customer_name" binding:"required"`
	Phone         string                   `json:"phone" binding:"required"`
	Email         *string                  `json:"email" binding:"omitempty,email"`
	Address       string                   `json:"address" binding:"required"`
	Note          *string                  `json:"note"`
	PaymentMethod string                   `json:"payment_method" binding:"required,oneof=pay_on_delivery paystack"`
	Items         []CreateOrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed out_for_delivery delivered cancelled"`
}