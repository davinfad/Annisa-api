package service

import (
	"annisa-api/models"
	"annisa-api/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type ServiceUser interface {
	RegisterUser(inputUser models.UserRegisterDTO) (*models.User, error)
	LoginUser(inputUser models.UserLoginDTO) (*models.User, error)
	IsUsernameAvailability(input string) (bool, error)
	GetUserByUsername(username string) (*models.User, error)
}

type serviceUser struct {
	repositoryUser repository.RepositoryUser
}

func NewUserService(repositoryUser repository.RepositoryUser) *serviceUser {
	return &serviceUser{repositoryUser}
}

func (s *serviceUser) RegisterUser(inputUser models.UserRegisterDTO) (*models.User, error) {
	user := &models.User{
		Username: inputUser.Username,
		Password: inputUser.Password,
		IDCabang: inputUser.IDCabang,
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(inputUser.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.Password = string(passwordHash)

	createUser, err := s.repositoryUser.Create(user)
	if err != nil {
		return createUser, err
	}
	return createUser, nil
}

func (s *serviceUser) LoginUser(inputUser models.UserLoginDTO) (*models.User, error) {
	username := inputUser.Username
	password := inputUser.Password

	checkUser, err := s.repositoryUser.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	if checkUser == nil {
		return nil, errors.New("user not found with that username")
	}

	err = bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(password))
	if err != nil {
		return nil, errors.New("incorrect password")
	}

	return checkUser, nil
}

func (s *serviceUser) IsUsernameAvailability(input string) (bool, error) {
	user, err := s.repositoryUser.FindByUsername(input)
	if err != nil {
		return false, err
	}

	if user == nil {
		return true, nil
	}

	return false, nil
}

func (s *serviceUser) GetUserByUsername(username string) (*models.User, error) {
	user, err := s.repositoryUser.FindByUsername(username)

	if err != nil {
		return user, err
	}

	if user == nil {
		return nil, errors.New("user not found with that username")
	}

	return user, nil
}
