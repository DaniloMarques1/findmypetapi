package repository

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/danilomarques1/findmypetapi/model"
	"github.com/danilomarques1/findmypetapi/util"
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

	err = stmt.QueryRow(post.Id, post.Author.Id, post.Title, post.Description, post.ImageUrl).Scan(
		&post.Status, &post.CreatedAt)
	if err != nil {
		log.Printf("Error executing statement %v\n", err)
		return err
	}

	author, err := pr.FindPostAuthor(post.Author.Id)
	if err != nil {
		log.Printf("Error finding author %v\n", err)
		return err
	}
	post.Author = author

	return nil
}

func (pr *PostRepositorySql) Update(post *model.Post) error {
	stmt, err := pr.db.Prepare(`
		update post
		set title = $1, description = $2, status = $3
		where id = $4
	`)
	if err != nil {
		log.Printf("Error creating statement update %v\n", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.Title, post.Description, post.Status, post.Id)
	if err != nil {
		log.Printf("Error executing update %v\n", err)
		return err
	}

	return nil
}

func (pr *PostRepositorySql) FindById(id string) (*model.Post, error) {
	stmt, err := pr.db.Prepare(`
		select id, author_id, title, description, image_url, status, created_at
		from post
		where id = $1`)

	if err != nil {
		log.Printf("Error building query %v\n", err)
		return nil, err
	}
	defer stmt.Close()

	var post model.Post
	err = stmt.QueryRow(id).Scan(&post.Id, &post.Author.Id, &post.Title,
		&post.Description, &post.ImageUrl, &post.Status, &post.CreatedAt)
	if err != nil {
		log.Printf("Error querying row %v\n", err)
		return nil, util.NewApiError("No post found", http.StatusNotFound)
	}

	return &post, nil
}

func (pr *PostRepositorySql) FindPostByAuthor(authorId, postId string) (*model.Post, error) {
	stmt, err := pr.db.Prepare(`
		select id, author_id, title, description, created_at
		from post
		where author_id = $1 and id = $2
	`)
	if err != nil {
		log.Printf("Error preparing statement post %v\n", err)
		return nil, err
	}
	defer stmt.Close()

	var post model.Post
	err = stmt.QueryRow(authorId, postId).Scan(&post.Id, &post.Author.Id,
		&post.Title, &post.Description, &post.CreatedAt)
	if err != nil {
		log.Printf("Error querying post %v\n", err)
		return nil, err
	}

	return &post, nil
}

func (pr *PostRepositorySql) FindAll() ([]model.Post, error) {
	stmt, err := pr.db.Prepare(`
		select id, author_id, title, description, image_url, status, created_at
		from post
		order by created_at desc`)
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
		err = rows.Scan(&post.Id, &post.Author.Id, &post.Title, &post.Description,
			&post.ImageUrl, &post.Status, &post.CreatedAt)
		if err != nil {
			log.Printf("Error scanning %v\n", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostRepositorySql) FindPostsByAuthor(authorId string) ([]model.Post, error) {
	stmt, err := pr.db.Prepare(`
		select id, author_id, title, description, image_url, status, created_at
		from post
		where author_id = $1
		order by created_at desc
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]model.Post, 0)
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.Id, &post.Author.Id, &post.Title, &post.Description, &post.ImageUrl,
			&post.Status, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostRepositorySql) FindPostAuthor(authorId string) (*model.User, error) {
	stmt, err := pr.db.Prepare(`
		select id, name, email, password_hash
		from userpet
		where id = $1
	`)
	if err != nil {
		log.Printf("Error Preparing statement %v\n", err)
		return nil, err
	}
	defer stmt.Close()
	var user model.User
	err = stmt.QueryRow(authorId).Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		log.Printf("Error querying user %v\n", err)
		return nil, err
	}

	return &user, nil 
}
