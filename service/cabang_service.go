package service

import (
	"annisa-api/models"
	"annisa-api/repository"
	"fmt"
	"time"
)

type ServiceCabang interface {
	Create(cabang *models.CabangDTO) (*models.Cabang, error)
	GetByID(id int) (*models.Cabang, error)
	GetAll() ([]*models.Cabang, error)
	Update(id int, cabang *models.CabangDTO, user *models.UpdateUserDTO) (*models.Cabang, error)
	Delete(id int) error
}

type serviceCabang struct {
	repositoryCabang repository.RepositoryCabang
	repositoryUser   repository.RepositoryUser
}

func NewCabangService(repositoryCabang repository.RepositoryCabang, repositoryUser repository.RepositoryUser) *serviceCabang {
	return &serviceCabang{repositoryCabang, repositoryUser}
}

func (s *serviceCabang) GetAll() ([]*models.Cabang, error) {
	return s.repositoryCabang.GetAll()
}

func (s *serviceCabang) GetByID(id int) (*models.Cabang, error) {
	cabang, err := s.repositoryCabang.GetByID(id)
	if err != nil {
		return nil, err
	}
	if cabang == nil {
		return nil, fmt.Errorf("cabang not found")
	}
	return cabang, nil
}

func (s *serviceCabang) Create(cabang *models.CabangDTO) (*models.Cabang, error) {
	layout := "15:04"

	jamBukaParsed, err := time.Parse(layout, cabang.JamBuka)
	if err != nil {
		return nil, fmt.Errorf("invalid jam_buka: %v", err)
	}

	jamTutupParsed, err := time.Parse(layout, cabang.JamTutup)
	if err != nil {
		return nil, fmt.Errorf("invalid jam_tutup: %v", err)
	}

	// ✅ Fix: Gunakan tahun 2000 agar valid
	jamBuka := time.Date(2000, time.January, 1, jamBukaParsed.Hour(), jamBukaParsed.Minute(), 0, 0, time.Local)
	jamTutup := time.Date(2000, time.January, 1, jamTutupParsed.Hour(), jamTutupParsed.Minute(), 0, 0, time.Local)

	input := &models.Cabang{
		NamaCabang: cabang.NamaCabang,
		KodeCabang: cabang.KodeCabang,
		JamBuka:    jamBuka,
		JamTutup:   jamTutup,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	create, err := s.repositoryCabang.Create(input)
	if err != nil {
		return nil, err
	}
	return create, nil
}

func (s *serviceCabang) Update(id int, cabang *models.CabangDTO, user *models.UpdateUserDTO) (*models.Cabang, error) {
	layout := "15:04"

	jamBukaParsed, err := time.Parse(layout, cabang.JamBuka)
	if err != nil {
		return nil, fmt.Errorf("invalid jam_buka: %v", err)
	}
	jamTutupParsed, err := time.Parse(layout, cabang.JamTutup)
	if err != nil {
		return nil, fmt.Errorf("invalid jam_tutup: %v", err)
	}

	cabang.JamBuka = time.Date(2000, time.January, 1, jamBukaParsed.Hour(), jamBukaParsed.Minute(), 0, 0, time.Local).Format("15:04:05")
	cabang.JamTutup = time.Date(2000, time.January, 1, jamTutupParsed.Hour(), jamTutupParsed.Minute(), 0, 0, time.Local).Format("15:04:05")

	updated, err := s.repositoryCabang.Update(id, cabang)
	if err != nil {
		return nil, err
	}

	if user != nil && user.Username != "" {
		if err := s.repositoryUser.UpdateByCabang(id, user); err != nil {
			return nil, err
		}
	}

	return updated, nil
}

func (s *serviceCabang) Delete(id int) error {
	cabang, err := s.repositoryCabang.GetByID(id)
	if err != nil || cabang == nil {
		return fmt.Errorf("cabang not found")
	}

	if err := s.repositoryUser.DeleteByCabang(id); err != nil {
		return fmt.Errorf("gagal menghapus user: %v", err)
	}

	return s.repositoryCabang.Delete(id)
}
