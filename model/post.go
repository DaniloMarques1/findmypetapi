package model

import "time"

type Post struct {
	Id          string    `json:"id"`
	Author      *User     `json:"author"`
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
	FindPostByAuthor(authorId, postId string) (*Post, error)
	FindPostsByAuthor(authorId string) ([]Post, error)
	FindAll() ([]Post, error)
}
