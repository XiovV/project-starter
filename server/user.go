package server

import (
	"fmt"
	"github.com/XiovV/starter-template/jwt"
	"github.com/XiovV/starter-template/models"
	"github.com/XiovV/starter-template/validator"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (s *Server) registerUserHandler(c *gin.Context) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type response struct {
		AccessToken string `json:"access_token"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		s.invalidJSONResponse(c)
		return
	}

	v := validator.New()

	v.Required(request.Username, "username")
	v.Required(request.Password, "password")

	if !v.IsValid() {
		s.invalidInputResponse(c, v)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		s.internalServerErrorResponse(c, err)
		return
	}

	user := &models.User{Username: request.Username, Password: string(hashedPassword)}

	createdUser, err := s.UserRepository.Create(user)
	if err != nil {
		c.Error(fmt.Errorf("registerUserHandler: %w", err))
		return
	}

	token, err := jwt.New(createdUser.ID, createdUser.Username)
	if err != nil {
		s.internalServerErrorResponse(c, err)
		return
	}

	res := response{AccessToken: token}

	c.JSON(http.StatusCreated, res)
}

func (s *Server) loginUserHandler(c *gin.Context) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type response struct {
		AccessToken string `json:"access_token"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		s.invalidJSONResponse(c)
		return
	}

	v := validator.New()

	v.Required(request.Username, "username")
	v.Required(request.Password, "password")

	if !v.IsValid() {
		s.invalidInputResponse(c, v)
		return
	}

	user, err := s.UserRepository.FindByUsername(request.Username)
	if err != nil {
		fmt.Println("user not found")
		c.JSON(http.StatusForbidden, gin.H{"error": "username or password is incorrect"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "username or password is incorrect"})
		return
	}

	token, err := jwt.New(user.ID, user.Username)
	if err != nil {
		s.internalServerErrorResponse(c, err)
		return
	}

	res := response{AccessToken: token}

	c.JSON(http.StatusOK, res)
}

func (s *Server) getPostsHandler(c *gin.Context) {
	type response struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	page, limit, err := s.validatePageAndLimit(c)
	if err != nil {
		s.badRequestResponse(c, err.Error())
		return
	}

	user := s.getUserFromContext(c)

	posts, err := s.PostRepository.FindByUserID(user.ID, page, limit)
	if err != nil {
		c.Error(fmt.Errorf("getPostsHandler: %w", err))
		return
	}

	var res []response
	for _, post := range posts {
		res = append(res, response{Title: post.Title, Body: post.Body})
	}

	c.JSON(http.StatusOK, gin.H{"posts": res})
}
