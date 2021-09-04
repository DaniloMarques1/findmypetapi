package repository

import (
	"database/sql"
	"log"

	"github.com/danilomarques1/findmypetapi/model"
)

type PostRepositorySql struct {
	db *sql.DB
}

func NewPostRepositorySql(db *sql.DB) *PostRepositorySql {
	return &PostRepositorySql{
		db: db,
	}
}

func (pr *PostRepositorySql) Save(post *model.Post) error {
	stmt, err := pr.db.Prepare(`
		insert into post(id, author_id, title, description, image_url)
		values($1, $2, $3, $4, $5)
		returning status, created_at`)
	if err != nil {
		log.Printf("Error creating statement %v\n", err)
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(post.Id, post.AuthorId, post.Title, post.Description, post.ImageUrl).Scan(
		&post.Status, &post.CreatedAt)
	if err != nil {
		log.Printf("Error executing statement %v\n", err)
		return err
	}

	return nil
}

func (pr *PostRepositorySql) Update(post *model.Post) error {
	return nil
}

func (pr *PostRepositorySql) FindById(id string) (*model.Post, error) {
	return nil, nil
}

func (pr *PostRepositorySql) FindPostByAuthor(author_id string) ([]model.Post, error) {
	return nil, nil
}
