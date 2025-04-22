package service

import (
	"annisa-api/models"
	"annisa-api/repository"
	"errors"
	"time"
)

type ServiceLayanan interface {
	Create(input *models.CreateLayananDTO) (*models.Layanan, error)
	GetByID(ID int) (*models.Layanan, error)
	GetAll() ([]*models.Layanan, error)
	Delete(ID int) (*models.Layanan, error)
	Update(ID int, input *models.CreateLayananDTO) (*models.Layanan, error)
}

type serviceLayanan struct {
	repositoryLayanan repository.RepositoryLayanan
}

func NewLayananService(repositoryLayanan repository.RepositoryLayanan) ServiceLayanan {
	return &serviceLayanan{repositoryLayanan}
}

func (s *serviceLayanan) Create(input *models.CreateLayananDTO) (*models.Layanan, error) {
	layanan := &models.Layanan{
		NamaLayanan:         input.NamaLayanan,
		PersenKomisi:        input.PersenKomisi,
		PersenKomisiLuarJam: input.PersenKomisiLuarJam,
		Kategori:            input.Kategori,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	return s.repositoryLayanan.Create(layanan)
}

func (s *serviceLayanan) GetByID(ID int) (*models.Layanan, error) {
	layanan, err := s.repositoryLayanan.GetByID(ID)
	if err != nil {
		return nil, err
	}
	if layanan == nil {
		return nil, errors.New("layanan not found")
	}
	return layanan, nil
}

func (s *serviceLayanan) GetAll() ([]*models.Layanan, error) {
	val, err := s.repositoryLayanan.GetAll()
	if err != nil {
		return val, err
	}

	return val, nil
}

func (s *serviceLayanan) Delete(ID int) (*models.Layanan, error) {
	layanan, err := s.repositoryLayanan.GetByID(ID)
	if err != nil {
		return nil, err
	}
	if layanan == nil {
		return nil, errors.New("layanan not found")
	}
	return s.repositoryLayanan.Delete(ID)
}

func (s *serviceLayanan) Update(ID int, input *models.CreateLayananDTO) (*models.Layanan, error) {
	existing, err := s.repositoryLayanan.GetByID(ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("layanan not found")
	}

	existing.NamaLayanan = input.NamaLayanan
	existing.PersenKomisi = input.PersenKomisi
	existing.PersenKomisiLuarJam = input.PersenKomisiLuarJam
	existing.Kategori = input.Kategori
	existing.UpdatedAt = time.Now()

	return s.repositoryLayanan.Update(existing)
}
