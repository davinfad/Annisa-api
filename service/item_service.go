package service

import (
	"annisa-api/models"
	"annisa-api/repository"
	"errors"
)

type ServiceItem interface {
	Create(input *models.CreateItemDTO) (*models.Item, error)
	GetByID(ID int) (*models.Item, error)
	GetAll() ([]*models.Item, error)
	Update(ID int, input *models.CreateItemDTO) (*models.Item, error)
	Delete(ID int) (*models.Item, error)
}

type serviceItem struct {
	repositoryItem repository.RepositoryItem
}

func NewItemService(repositoryItem repository.RepositoryItem) ServiceItem {
	return &serviceItem{repositoryItem}
}

func validateBatas(batasBawah, batasAtas int) error {
	if batasBawah < 0 || batasAtas < 0 {
		return errors.New("batas_bawah and batas_atas cannot be negative")
	}
	if batasAtas > 0 && batasBawah > batasAtas {
		return errors.New("batas_bawah cannot be greater than batas_atas")
	}
	return nil
}

func (s *serviceItem) Create(input *models.CreateItemDTO) (*models.Item, error) {
	if err := validateBatas(input.BatasBawah, input.BatasAtas); err != nil {
		return nil, err
	}

	now := nowWIB()
	item := &models.Item{
		NamaItem:   input.NamaItem,
		Satuan:     input.Satuan,
		BatasBawah: input.BatasBawah,
		BatasAtas:  input.BatasAtas,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	return s.repositoryItem.Create(item)
}

func (s *serviceItem) GetByID(ID int) (*models.Item, error) {
	item, err := s.repositoryItem.GetByID(ID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("item not found")
	}
	return item, nil
}

func (s *serviceItem) GetAll() ([]*models.Item, error) {
	return s.repositoryItem.GetAll()
}

func (s *serviceItem) Update(ID int, input *models.CreateItemDTO) (*models.Item, error) {
	if err := validateBatas(input.BatasBawah, input.BatasAtas); err != nil {
		return nil, err
	}

	existing, err := s.repositoryItem.GetByID(ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("item not found")
	}

	existing.NamaItem = input.NamaItem
	existing.Satuan = input.Satuan
	existing.BatasBawah = input.BatasBawah
	existing.BatasAtas = input.BatasAtas
	existing.UpdatedAt = nowWIB()

	return s.repositoryItem.Update(existing)
}

func (s *serviceItem) Delete(ID int) (*models.Item, error) {
	item, err := s.repositoryItem.GetByID(ID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("item not found")
	}
	return s.repositoryItem.Delete(ID)
}
