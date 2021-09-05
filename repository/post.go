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
	stmt, err := pr.db.Prepare(`
		select id, author_id, title, description, image_url, status, created_at
		from post
		where id = $1`)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var post model.Post
	err = stmt.QueryRow(id).Scan(&post.Id, &post.AuthorId, &post.Title,
		&post.Description, &post.ImageUrl, &post.Status, &post.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (pr *PostRepositorySql) FindPostByAuthor(author_id string) ([]model.Post, error) {
	return nil, nil
}

func (pr *PostRepositorySql) FindAll() ([]model.Post, error) {
	stmt, err := pr.db.Prepare(`
		select id, author_id, title, description, image_url, status, created_at
		from post`)
	if err != nil {
		log.Printf("Error preparing statement %v\n", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Printf("Error querying %v\n", err)
		return nil, err
	}
	defer rows.Close()
	var posts []model.Post
	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.Id, &post.AuthorId, &post.Title, &post.Description,
			&post.ImageUrl, &post.Status, &post.CreatedAt)
		if err != nil {
			log.Printf("Error scanning %v\n", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
