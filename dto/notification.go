package dto

// represents the message that will be sent to
// the notification service about a new comment
// on a post
type CommentNotification struct {
	PostId             string `json:"post_id"`
	PostAuthorEmail    string `json:"post_author_email"`
	PostAuthorName     string `json:"post_author_name"`
	CommentAuthorEmail string `json:"comment_author_email"`
	CommentAuthorName  string `json:"comment_author_name"`
}

// represents the message that will be sen to
// the notification service about a status update
// on a post (found pet)
type StatusChangeNotification struct {
	PostId string `json:"post_id"`
}
