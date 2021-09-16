package repository

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/danilomarques1/findmypetapi/dto"
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
	// TODO find a way to verify if post exist first
	stmt, err := cr.db.Prepare(`
		select id, author_id, post_id, comment_text, created_at
		from comment
		where post_id = $1
	`)
	if err != nil {
		log.Printf("Error scanning %v\n", err)
		return nil, err
	}
	defer stmt.Close()

	comments := make([]model.Comment, 0)
	rows, err := stmt.Query(postId)
	if err != nil {
		log.Printf("Error scanning %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		log.Printf("Reading a new line...\n")
		var comment model.Comment
		err = rows.Scan(&comment.Id, &comment.AuthorId, &comment.PostId,
			&comment.CommentText, &comment.CreatedAt)
		if err != nil {
			log.Printf("Error scanning %v\n", err)
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (cr *CommentRepositorySql) GetCommentNotificationMessage(postId, commentId string) ([]byte, error) {
	stmt, err := cr.db.Prepare(`
		select pu.email, pu.name, cou.email, cou.name
		from comment as c
		join post as p on c.post_id = p.id
		join userpet as cou on c.author_id = cou.id
		join userpet as pu on p.author_id = pu.id
		where c.id = $1
	`)
	if err != nil {
		log.Printf("Error statement %v\n", err)
		return nil, err
	}
	defer stmt.Close()

	var msg dto.CommentNotification
	err = stmt.QueryRow(commentId).Scan(&msg.PostAuthorEmail,
		&msg.PostAuthorName, &msg.CommentAuthorEmail,
		&msg.CommentAuthorName)
	msg.PostId = postId
	if err != nil {
		log.Printf("Error querying %v\n", err)
		return nil, err
	}
	if msg.PostAuthorEmail == msg.CommentAuthorEmail {
		return []byte(""), nil
	}

	mBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	log.Printf("%v\n", string(mBytes))

	return mBytes, nil
}
