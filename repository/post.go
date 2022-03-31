package repository

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type Post struct {
	ID     int    `db:"id"`
	UserID int    `db:"user_id"`
	Title  string `db:"title"`
	Body   string `db:"body"`
}

type PostRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (pr *PostRepository) FindByUserID(userId, page, limit int) ([]Post, error) {
	var p []Post

	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	err := pr.db.SelectContext(ctx, &p, "SELECT * FROM posts WHERE user_id = $1 LIMIT $2 OFFSET $3", userId, limit, calculateOffset(page, limit))
	if err != nil {
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, notFoundErr("post")
			}
		}
	}

	return p, nil
}
