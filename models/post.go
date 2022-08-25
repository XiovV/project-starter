package models

type Post struct {
	ID     int    `db:"id"`
	UserID int    `db:"user_id"`
	Title  string `db:"title"`
	Body   string `db:"body"`
}

type PostService interface {
	Create(post *Post) (Post, error)
	FindByUserID(userId, page, limit int) ([]Post, error)
	FindByPostID(postId int) (Post, error)
}
