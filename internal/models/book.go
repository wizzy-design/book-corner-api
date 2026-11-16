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

type BookLookupResult struct {
	Title         string `json:"title"`
	Author        string `json:"author"` 
	Description   string `json:"description"`
	ImageUrl      string `json:"image_url"`
	SuggestedTags []string `json:"suggested_tags"`
}

type CreateBookRequest struct {
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       int    `json:"price" binding:"required,min=0"`
	Stock       *int   `json:"stock" binding:"required,min=0"`
	ImageUrl    string `json:"image_url"`
	Status      string `json:"status" binding:"omitempty,oneof=draft published"`
	TagIDs      []int  `json:"tag_ids"`
}