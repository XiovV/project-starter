package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Active   bool   `db:"active"`
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Create(user *User) (User, error) {
	var u User

	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	err := ur.db.QueryRowxContext(ctx, "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING *;", user.Username, user.Password).StructScan(&u)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return User{}, fmt.Errorf("ur.Create: %w", alreadyExistsErr("a user with this username already exists"))
			}
		}

		return User{}, err
	}

	return u, nil
}

func (ur *UserRepository) FindByID(id int) (User, error) {
	var u User

	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	err := ur.db.GetContext(ctx, &u, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, fmt.Errorf("ur.FindByID: %w", notFoundErr("user not found"))
		}
	}

	return u, nil
}

func (ur *UserRepository) FindByUsername(username string) (User, error) {
	var u User

	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	err := ur.db.GetContext(ctx, &u, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, fmt.Errorf("ur.FindByUsername: %w", notFoundErr("user not found"))
		}
	}

	return u, nil
}
