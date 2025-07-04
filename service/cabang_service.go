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
}

type serviceCabang struct {
	repositoryCabang repository.RepositoryCabang
}

func NewCabangService(repositoryCabang repository.RepositoryCabang) *serviceCabang {
	return &serviceCabang{repositoryCabang}
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
