package main

import "time"

type Movie struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Duration    int       `json:"duration"`
	ReleaseYear int       `json:"release_year"`
	Director    string    `json:"director"`
	Rating      float64   `json:"rating"`
	Categories  []string  `json:"categories"`
	CreatedAt   time.Time `json:"-"`
	ModifiedAt  time.Time `json:"-"`
}

type MovieCreate struct {
	Name        string   `json:"name"`
	Duration    int      `json:"duration"`
	ReleaseYear int      `json:"release_year"`
	Director    string   `json:"director"`
	Rating      float64  `json:"rating"`
	Categories  []string `json:"categories"`
}

type MoviePatch struct {
	Name        *string   `json:"name,omitempty"`
	Duration    *int      `json:"duration,omitempty"`
	ReleaseYear *int      `json:"release_year,omitempty"`
	Director    *string   `json:"director,omitempty"`
	Rating      *float64  `json:"rating,omitempty"`
	Categories  *[]string `json:"categories,omitempty"`
}
