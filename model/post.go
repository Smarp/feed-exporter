package model

import "time"

type Post struct {
	Id            string     `json:"_id"`
	Body          string     `json:"body"`
	Title         string     `json:"title"`
	Image         string     `json:"imageUrl"`
	PublishedDate *time.Time `json:"publishedAt"`
}
