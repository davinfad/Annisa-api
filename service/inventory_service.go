package service

import (
	"annisa-api/models"
	"annisa-api/repository"
	"errors"
	"time"
)

type ServiceInventory interface {
	GetByCabang(idCabang int) ([]*models.InventoryStokView, error)
	AdjustStok(idCabang, idItem, delta int) (*models.Inventory, error)
	SetStok(idCabang, idItem, stok int) (*models.Inventory, error)
}

type serviceInventory struct {
	repositoryInventory repository.RepositoryInventory
	repositoryItem      repository.RepositoryItem
}

func NewInventoryService(repositoryInventory repository.RepositoryInventory, repositoryItem repository.RepositoryItem) ServiceInventory {
	return &serviceInventory{repositoryInventory, repositoryItem}
}

func nowWIB() time.Time {
	return time.Now().In(time.FixedZone("WIB", 7*3600))
}

func (s *serviceInventory) GetByCabang(idCabang int) ([]*models.InventoryStokView, error) {
	return s.repositoryInventory.GetByCabang(idCabang)
}

func (s *serviceInventory) AdjustStok(idCabang, idItem, delta int) (*models.Inventory, error) {
	if err := s.ensureItemExists(idItem); err != nil {
		return nil, err
	}
	return s.repositoryInventory.AdjustStok(idItem, idCabang, delta, nowWIB())
}

func (s *serviceInventory) SetStok(idCabang, idItem, stok int) (*models.Inventory, error) {
	if stok < 0 {
		return nil, errors.New("stok cannot be negative")
	}
	if err := s.ensureItemExists(idItem); err != nil {
		return nil, err
	}
	return s.repositoryInventory.SetStok(idItem, idCabang, stok, nowWIB())
}

func (s *serviceInventory) ensureItemExists(idItem int) error {
	item, err := s.repositoryItem.GetByID(idItem)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("item not found")
	}
	return nil
}
