package repository

import (
	"database/sql"
	"errors"
	"github.com/XiovV/starter-template/models"
	"github.com/jmoiron/sqlx"
)

type PostRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (pr *PostRepository) Create(post *models.Post) (models.Post, error) {
	var p models.Post

	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	err := pr.db.QueryRowxContext(ctx, "INSERT INTO posts (user_id, title, body) VALUES ($1, $2, $3) RETURNING *;", post.UserID, post.Title, post.Body).StructScan(&p)
	if err != nil {
		return models.Post{}, err
	}

	return p, nil
}

func (pr *PostRepository) FindByUserID(userId, page, limit int) ([]models.Post, error) {
	var p []models.Post

	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	err := pr.db.SelectContext(ctx, &p, "SELECT * FROM posts WHERE user_id = $1 LIMIT $2 OFFSET $3", userId, limit, calculateOffset(page, limit))
	if err != nil {
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, notFoundErr("this user has no posts")
			}
		}
	}

	return p, nil
}

func (pr *PostRepository) FindByPostID(postId int) (models.Post, error) {
	var p models.Post

	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	err := pr.db.SelectContext(ctx, &p, "SELECT * FROM posts WHERE id = $1", postId)
	if err != nil {
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return models.Post{}, notFoundErr("post does not exist")
			}
		}
	}

	return p, nil
}

//func (pr *PostRepository) DeleteByPostID(postId int) error {
//	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
//	defer cancel()
//
//	return pr.db.ExecContext(ctx, )
//}
