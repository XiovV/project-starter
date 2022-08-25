package mock

import (
	"github.com/XiovV/starter-template/models"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Create(user *models.User) (models.User, error) {
	args := m.Called()
	result := args.Get(0)
	return result.(models.User), nil
}

func (m *UserRepository) FindByID(id int) (models.User, error) {
	args := m.Called()
	result := args.Get(0)
	return result.(models.User), nil
}

func (m *UserRepository) FindByUsername(username string) (models.User, error) {
	args := m.Called()
	resultUser := args.Get(0).(models.User)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(resultUser.Password), bcrypt.DefaultCost)

	resultUser.Password = string(hashedPassword)
	
	return resultUser, nil
}
