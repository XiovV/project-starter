package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/XiovV/starter-template/repository"
	"github.com/XiovV/starter-template/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("starting up on port", os.Getenv("PORT"))

	db := repository.NewPostgres()
	userRepository := repository.NewUserRepository(db)
	postRepository := repository.NewPostRepository(db)

	srv := &server.Server{
		UserRepository: userRepository,
		PostRepository: postRepository,
	}

	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler:      srv.New(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
