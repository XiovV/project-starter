package server

import (
	"github.com/XiovV/starter-template/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (s *Server) successResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func (s *Server) badRequestResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": msg})
}

func (s *Server) invalidJSONResponse(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
}

func (s *Server) invalidInputResponse(c *gin.Context, v *validator.Validator) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": v.Errors()})
}

func (s *Server) internalServerErrorResponse(c *gin.Context, err error) {
	s.Logger.Error("INTERNAL SERVER ERROR", zap.String("err", err.Error()))
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
