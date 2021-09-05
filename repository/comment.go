package repository

import (
	"database/sql"

	"github.com/danilomarques1/findmypetapi/model"
)

type CommentRepositorySql struct {
	db *sql.DB
}

func NewCommentRepositorySql(db *sql.DB) *CommentRepositorySql {
	return &CommentRepositorySql{
		db: db,
	}
}

func (cr *CommentRepositorySql) Save(comment *model.Comment) error {
	stmt, err := cr.db.Prepare(`
		insert into comment(id, author_id, post_id, comment_text)
		values($1, $2, $3, $4)
		returning created_at
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(comment.Id, comment.AuthorId,
		comment.PostId, comment.CommentText).Scan(&comment.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (cr *CommentRepositorySql) FindAll(postId string) ([]model.Comment, error) {
	return nil, nil
}
