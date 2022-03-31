package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

const (
	MaxOpenConns        = 25
	MaxIdleConns        = 25
	MaxIdleTime         = "15m"
	DefaultQueryTimeout = 5
)

func NewPostgres() *sqlx.DB {
	db, err := sqlx.Connect("postgres", os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatalln(err)
	}

	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)

	duration, err := time.ParseDuration(MaxIdleTime)
	if err != nil {
		log.Fatalln(err)
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

func newBackgroundContext(duration int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
}

func calculateOffset(page, limit int) int {
	return (page - 1) * limit
}
