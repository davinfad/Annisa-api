package service

import (
	"annisa-api/models"
	"annisa-api/repository"
	"errors"
	"time"
)

type ServiceInventory interface {
	Create(input *models.CreateInventoryDTO) (*models.Inventory, error)
	GetByID(ID int) (*models.Inventory, error)
	GetByCabang(idCabang int) ([]*models.Inventory, error)
	Update(ID int, input *models.CreateInventoryDTO) (*models.Inventory, error)
	Delete(ID int) (*models.Inventory, error)
	AdjustStok(ID, delta int) (*models.Inventory, error)
}

type serviceInventory struct {
	repositoryInventory repository.RepositoryInventory
}

func NewInventoryService(repositoryInventory repository.RepositoryInventory) ServiceInventory {
	return &serviceInventory{repositoryInventory}
}

func nowWIB() time.Time {
	return time.Now().In(time.FixedZone("WIB", 7*3600))
}

func validateBatas(batasBawah, batasAtas, stok int) error {
	if batasBawah < 0 || batasAtas < 0 || stok < 0 {
		return errors.New("batas_bawah, batas_atas, and stok cannot be negative")
	}
	if batasAtas > 0 && batasBawah > batasAtas {
		return errors.New("batas_bawah cannot be greater than batas_atas")
	}
	return nil
}

func (s *serviceInventory) Create(input *models.CreateInventoryDTO) (*models.Inventory, error) {
	if err := validateBatas(input.BatasBawah, input.BatasAtas, input.Stok); err != nil {
		return nil, err
	}

	now := nowWIB()
	inv := &models.Inventory{
		IDCabang:   input.IDCabang,
		NamaItem:   input.NamaItem,
		BatasBawah: input.BatasBawah,
		BatasAtas:  input.BatasAtas,
		Stok:       input.Stok,
		Satuan:     input.Satuan,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	return s.repositoryInventory.Create(inv)
}

func (s *serviceInventory) GetByID(ID int) (*models.Inventory, error) {
	inv, err := s.repositoryInventory.GetByID(ID)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, errors.New("inventory not found")
	}
	return inv, nil
}

func (s *serviceInventory) GetByCabang(idCabang int) ([]*models.Inventory, error) {
	return s.repositoryInventory.GetByCabang(idCabang)
}

func (s *serviceInventory) Update(ID int, input *models.CreateInventoryDTO) (*models.Inventory, error) {
	if err := validateBatas(input.BatasBawah, input.BatasAtas, input.Stok); err != nil {
		return nil, err
	}

	existing, err := s.repositoryInventory.GetByID(ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("inventory not found")
	}

	existing.IDCabang = input.IDCabang
	existing.NamaItem = input.NamaItem
	existing.BatasBawah = input.BatasBawah
	existing.BatasAtas = input.BatasAtas
	existing.Stok = input.Stok
	existing.Satuan = input.Satuan
	existing.UpdatedAt = nowWIB()

	return s.repositoryInventory.Update(existing)
}

func (s *serviceInventory) Delete(ID int) (*models.Inventory, error) {
	inv, err := s.repositoryInventory.GetByID(ID)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, errors.New("inventory not found")
	}
	return s.repositoryInventory.Delete(ID)
}

func (s *serviceInventory) AdjustStok(ID, delta int) (*models.Inventory, error) {
	existing, err := s.repositoryInventory.GetByID(ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("inventory not found")
	}

	affected, err := s.repositoryInventory.AdjustStok(ID, delta, nowWIB())
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		// item exists (checked above), so the change would make stok negative
		return nil, errors.New("insufficient stock")
	}

	return s.repositoryInventory.GetByID(ID)
}
