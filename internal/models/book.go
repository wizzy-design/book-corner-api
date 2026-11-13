package models

import "time"

type Tag struct {
	ID int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Kind string `json:"kind" db:"kind"`
}

type Book struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Author      string `json:"author" db:"author"`
	Description string `json:"description" db:"description"`
	Price       int    `json:"price" db:"price"`
	Stock       int    `json:"stock" db:"stock"`
	ImageUrl    string `json:"image_url" db:"image_url"`
	Status      string `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Tags []Tag `json:"tags,omitempty"`
}

type BookFilter struct {
	Query string
	Genre string
	Tag string
	Status string
	Page int
	Limit int
}