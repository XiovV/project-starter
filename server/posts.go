package server

import (
	"fmt"
	"github.com/XiovV/starter-template/models"
	"github.com/XiovV/starter-template/validator"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (s *Server) createPostHandler(c *gin.Context) {
	user := s.getUserFromContext(c)

	var request struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	type response struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		s.invalidJSONResponse(c)
		return
	}

	v := validator.New()

	v.Required(request.Title, "title")
	v.Required(request.Body, "body")

	if !v.IsValid() {
		s.invalidInputResponse(c, v)
		return
	}

	post := &models.Post{
		UserID: user.ID,
		Title:  request.Title,
		Body:   request.Body,
	}

	createdPost, err := s.PostRepository.Create(post)
	if err != nil {
		c.Error(fmt.Errorf("createPostHandler: %w", err))
		return
	}

	res := response{
		ID:    createdPost.ID,
		Title: createdPost.Title,
		Body:  createdPost.Body,
	}

	c.JSON(http.StatusCreated, gin.H{"post": res})
}

func (s *Server) deletePostHandler(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		s.badRequestResponse(c, "post id needs to be an integer")
		return
	}

	post, err := s.PostRepository.FindByPostID(postId)
	if err != nil {
		c.Error(fmt.Errorf("deletePostHandler: %w", err))
		return
	}

	fmt.Println(post)
}
