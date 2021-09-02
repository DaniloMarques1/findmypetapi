package model

import "time"

type Post struct {
	Id          string    `json:"id"`
	AuthorId    string    `json:"author_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageUrl    string    `json:"image_url"`
	Status      string    `json:"status"` // TODO enum like structure [missing, found]
	CreatedAt   time.Time `json:"created_at"`
}

type PostRepository interface {
	Save(*Post) error
	FindById(id string) (*Post, error)
	Update(*Post) error
	FindPostByAuthor(author_id string) ([]Post, error)
}
