package server

import (
	"fmt"
	"github.com/XiovV/starter-template/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	MaxLimitValue = 100
	MinLimitValue = 1
	MinPageValue  = 1
)

func (s *Server) getUserFromContext(c *gin.Context) models.User {
	userCtx, exists := c.Get("user")
	if !exists {
		s.internalServerErrorResponse(c, fmt.Errorf("user not found in context"))
		return models.User{}
	}

	return userCtx.(models.User)
}

func (s *Server) validatePageAndLimit(c *gin.Context) (int, int, error) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		return 0, 0, fmt.Errorf("page must be an integer")
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		return 0, 0, fmt.Errorf("limit must be an integer")
	}

	if page < MinPageValue {
		return 0, 0, fmt.Errorf("page must be greater than 0")
	}

	if limit < MinLimitValue {
		return 0, 0, fmt.Errorf("limit must be greater than 0")
	}

	if limit > MaxLimitValue {
		return 0, 0, fmt.Errorf("maximum limit size is %d", MaxLimitValue)
	}

	return page, limit, nil
}

func (s *Server) newError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}
