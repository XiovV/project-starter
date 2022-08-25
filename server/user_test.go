package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/XiovV/starter-template/mock"
	"github.com/XiovV/starter-template/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockUser := models.User{Username: "test", Password: "test"}

		mockUserRepository := mock.UserRepository{}
		mockUserRepository.On("Create").Return(mockUser, nil)
		server := NewMockServer(&mockUserRepository, nil)

		w := httptest.NewRecorder()

		var request struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		request.Username = mockUser.Username
		request.Password = mockUser.Password

		j, _ := json.Marshal(request)

		req, _ := http.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer(j))
		server.ServeHTTP(w, req)

		var response struct {
			AccessToken string `json:"access_token"`
		}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			assert.NotNil(t, err)
		}

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Greater(t, len(response.AccessToken), 0)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		server := NewMockServer(&mock.UserRepository{}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/v1/users/register", nil)
		server.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "{\"error\":\"invalid json\"}", w.Body.String())
	})
}

func TestLoginUser(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockUser := models.User{Username: "test", Password: "test"}

		mockUserRepository := mock.UserRepository{}
		mockUserRepository.On("FindByUsername").Return(mockUser, nil)
		server := NewMockServer(&mockUserRepository, nil)

		w := httptest.NewRecorder()

		var request struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		request.Username = mockUser.Username
		request.Password = mockUser.Password

		j, _ := json.Marshal(request)

		req, _ := http.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewBuffer(j))
		server.ServeHTTP(w, req)

		var response struct {
			AccessToken string `json:"access_token"`
		}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			assert.NotNil(t, err)
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Greater(t, len(response.AccessToken), 0)
	})

	t.Run("User not found", func(t *testing.T) {
		mockUser := models.User{Username: "test", Password: "test"}

		mockUserRepository := mock.UserRepository{}
		mockUserRepository.On("FindByUsername").Return(models.User{}, errors.New("db user not found"))
		server := NewMockServer(&mockUserRepository, nil)

		w := httptest.NewRecorder()

		var request struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		request.Username = mockUser.Username
		request.Password = mockUser.Password

		j, _ := json.Marshal(request)

		req, _ := http.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewBuffer(j))
		server.ServeHTTP(w, req)

		var response struct {
			AccessToken string `json:"access_token"`
		}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			assert.NotNil(t, err)
		}

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Equal(t, 0, len(response.AccessToken))
	})

	t.Run("Invalid Password", func(t *testing.T) {
		mockUser := models.User{Username: "test", Password: "test"}

		mockUserRepository := mock.UserRepository{}
		mockUserRepository.On("FindByUsername").Return(mockUser, nil)
		server := NewMockServer(&mockUserRepository, nil)

		w := httptest.NewRecorder()

		var request struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		request.Username = mockUser.Username
		request.Password = "someRandomPassword"

		j, _ := json.Marshal(request)

		req, _ := http.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewBuffer(j))
		server.ServeHTTP(w, req)

		var response struct {
			AccessToken string `json:"access_token"`
		}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			assert.NotNil(t, err)
		}

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Equal(t, 0, len(response.AccessToken))
	})
}
