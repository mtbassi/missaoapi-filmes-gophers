package main

import "time"

type Movie struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Duration  int       `json:"duration"`
	CreatedAt time.Time `json:"-"`
}

type MovieCreate struct {
	Name     string `json:"name"`
	Duration int    `json:"duration"`
}

type MoviePatch struct {
	Name     *string `json:"name,omitempty"`
	Duration *int    `json:"duration,omitempty"`
}
