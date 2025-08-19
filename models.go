package main

import "time"

type Movie struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Duration   int       `json:"duration"`
	Categories []string  `json:"categories"`
	CreatedAt  time.Time `json:"-"`
}

type MovieCreate struct {
	Name       string   `json:"name"`
	Duration   int      `json:"duration"`
	Categories []string `json:"categories"`
}

type MoviePatch struct {
	Name       *string   `json:"name,omitempty"`
	Duration   *int      `json:"duration,omitempty"`
	Categories *[]string `json:"categories,omitempty"`
}
