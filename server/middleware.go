package server

import (
	"errors"
	"fmt"
	"github.com/XiovV/starter-template/jwt"
	"github.com/XiovV/starter-template/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (s *Server) errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()

		if err != nil {
			var alreadyExistsErr *repository.AlreadyExistsErr
			var notFoundErr *repository.NotFoundErr

			switch {
			case errors.As(err, &alreadyExistsErr):
				s.newError(c, http.StatusConflict, alreadyExistsErr)
			case errors.As(err, &notFoundErr):
				s.newError(c, http.StatusNotFound, notFoundErr)
			default:
				s.internalServerErrorResponse(c, err)
			}
		}
	}
}

func (s *Server) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func (s *Server) userAuth(c *gin.Context) {
	tokenHeader := c.GetHeader("Authorization")

	if len(tokenHeader) == 0 {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "did not receive Authorization header"})
		return
	}

	authorizationHeaderSplit := strings.Split(tokenHeader, " ")
	if len(authorizationHeaderSplit) != 2 {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "wrong Authorization header format"})
		return
	}

	if authorizationHeaderSplit[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "wrong Authorization header format"})
		return
	}

	authToken := authorizationHeaderSplit[1]

	token, err := jwt.Validate(authToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid token"})
		return
	}

	userId := jwt.GetClaimInt(token, "id")

	user, err := s.UserRepository.FindByID(userId)
	if err != nil {
		c.Error(fmt.Errorf("userAuth: %w", err))
		return
	}

	if !user.Active {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user is inactive"})
		return
	}

	c.Set("user", user)

	c.Next()
}
