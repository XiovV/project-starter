package server

import (
	"github.com/XiovV/starter-template/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	LOCAL_ENV      = "LOCAL"
	STAGING_ENV    = "STAGING"
	PRODUCTION_ENV = "PROD"
)

type Server struct {
	UserRepository models.UserService
	PostRepository models.PostService
	Logger         *zap.Logger
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

	postsAuth := v1.Group("/posts")
	postsAuth.Use(s.userAuth)
	{
		postsAuth.GET("/", s.getPostsHandler)
		postsAuth.POST("/", s.createPostHandler)
		postsAuth.DELETE("/:postId", s.deletePostHandler)
	}

	return router
}

func NewMockServer(userRepository models.UserService, postRepository models.PostService) *gin.Engine {
	server := Server{UserRepository: userRepository}

	return server.New()
}
