package main

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Movie struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Duration    int       `json:"duration"`
	ReleaseYear int       `json:"release_year"`
	Director    string    `json:"director"`
	Rating      float64   `json:"rating"`
	Categories  []string  `json:"categories"`
	Comments    []Comment `json:"comments"`
	CreatedAt   time.Time `json:"-"`
	ModifiedAt  time.Time `json:"-"`
}

type Comment struct {
	ID        string    `json:"-"`
	MovieID   string    `json:"-"`
	UserID    string    `json:"-"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"-"`
}

type MovieCreate struct {
	Name        string    `json:"name" validate:"required,min=3,max=150"`
	Duration    int       `json:"duration" validate:"required,min=15,max=360"`
	ReleaseYear int       `json:"release_year" validate:"required,yearValid"`
	Director    string    `json:"director" validate:"required,min=5,max=30"`
	Rating      float64   `json:"rating" validate:"required,min=1,max=5"`
	Categories  []string  `json:"categories" validate:"required,min=1,max=5"`
	Comments    []Comment `json:"comments" validate:"omitempty,min=1,max=100"`
}

type MoviePatch struct {
	Name        *string   `json:"name,omitempty" validate:"omitempty,min=3,max=150"`
	Duration    *int      `json:"duration,omitempty" validate:"omitempty,min=15,max=360"`
	ReleaseYear *int      `json:"release_year,omitempty" validate:"omitempty,yearValid"`
	Director    *string   `json:"director,omitempty" validate:"omitempty,min=5,max=30"`
	Rating      *float64  `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	Categories  *[]string `json:"categories,omitempty" validate:"omitempty,min=1,max=5"`
}

func yearValid(fl validator.FieldLevel) bool {
	year := fl.Field().Int()
	currentYear := int64(time.Now().Year())
	return year >= 1888 && year <= currentYear
}
