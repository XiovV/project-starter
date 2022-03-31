package server

import (
	"github.com/XiovV/starter-template/repository"
	"github.com/gin-gonic/gin"
)

type Server struct {
	UserRepository *repository.UserRepository
	PostRepository *repository.PostRepository
}

func (s *Server) New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), s.CORS(), s.errorHandler())

	v1 := router.Group("/v1")

	usersPublic := v1.Group("/users")
	{
		usersPublic.POST("/register", s.registerUserHandler)
		usersPublic.POST("/login", s.loginUserHandler)
	}

	usersAuth := v1.Group("/users")
	usersAuth.Use(s.userAuth)
	{
		usersAuth.GET("/posts", s.getPostsHandler)
	}

	return router
}
