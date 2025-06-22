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
	cabangService  ServiceCabang
}

func NewUserService(repositoryUser repository.RepositoryUser, cabangService ServiceCabang) *serviceUser {
	return &serviceUser{repositoryUser, cabangService}
}

func (s *serviceUser) RegisterUser(inputUser models.UserRegisterDTO) (*models.User, error) {
	var idCabang *int = inputUser.IDCabang

	// Buat cabang baru jika tidak ada ID
	if idCabang == nil && inputUser.CabangName != "" {
		newCabang := &models.CabangDTO{
			NamaCabang: inputUser.CabangName,
			KodeCabang: inputUser.KodeCabang,
			JamBuka:    inputUser.JamBuka,
			JamTutup:   inputUser.JamTutup,
		}
		cabang, err := s.cabangService.Create(newCabang)
		if err != nil {
			return nil, err
		}
		idCabang = &cabang.IDCabang
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(inputUser.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:   inputUser.Username,
		Password:   string(passwordHash),
		AccessCode: inputUser.AccessCode,
		IDCabang:   idCabang,
	}

	_, err = s.repositoryUser.Create(user)
	if err != nil {
		return nil, err
	}

	// ✅ Ambil ulang user agar Cabangs keisi
	fullUser, err := s.repositoryUser.FindByUsername(user.Username)
	if err != nil {
		return nil, err
	}

	return fullUser, nil
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
