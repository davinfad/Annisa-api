package service

import (
	"annisa-api/models"
	"annisa-api/repository"
	"errors"
	"fmt"
	"time"
)

type ServiceMember interface {
	Create(input *models.CreateMemberDTO) (*models.Member, error)
	GetByID(ID int) (*models.Member, error)
	GetAll() ([]*models.Member, error)
	GetMemberByCabangID(IDCabang int) ([]*models.Member, error)
	Update(ID int, input *models.CreateMemberDTO) (*models.Member, error)
	Delete(ID int) (*models.Member, error)
}

type serviceMember struct {
	repositoryMember repository.RepositoryMember
}

func NewMemberService(repositoryMember repository.RepositoryMember) ServiceMember {
	return &serviceMember{repositoryMember}
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func (s *serviceMember) Create(input *models.CreateMemberDTO) (*models.Member, error) {
	tglLahir, err := parseDate(input.TanggalLahir)
	if err != nil {
		return nil, errors.New("invalid tanggal_lahir format, expected YYYY-MM-DD")
	}

	tglDaftar, err := parseDate(input.TanggalDaftar)
	if err != nil {
		return nil, errors.New("invalid tanggal_daftar format, expected YYYY-MM-DD")
	}

	member := &models.Member{
		NomorPelanggan: input.NomorPelanggan,
		NamaMember:     input.NamaMember,
		NomorTelepon:   input.NomorTelepon,
		Alamat:         input.Alamat,
		TanggalLahir:   tglLahir,
		TanggalDaftar:  tglDaftar,
		IDCabang:       input.IDCabang,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return s.repositoryMember.Create(member)
}

func (s *serviceMember) GetByID(ID int) (*models.Member, error) {
	member, err := s.repositoryMember.Get(ID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errors.New("member not found")
	}
	return member, nil
}

func (s *serviceMember) GetAll() ([]*models.Member, error) {
	get, err := s.repositoryMember.GetAll()
	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *serviceMember) GetMemberByCabangID(IDCabang int) ([]*models.Member, error) {
	get, err := s.repositoryMember.GetMemberByIDCabang(IDCabang)

	if err != nil {
		return get, err
	}

	if len(get) == 0 {
		return nil, fmt.Errorf("tidak ditemukan member untuk id_cabang %d", IDCabang)
	}

	return get, nil
}

func (s *serviceMember) Update(ID int, input *models.CreateMemberDTO) (*models.Member, error) {
	existing, err := s.repositoryMember.Get(ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("member not found")
	}

	tglLahir, err := parseDate(input.TanggalLahir)
	if err != nil {
		return nil, errors.New("invalid tanggal_lahir format, expected YYYY-MM-DD")
	}

	tglDaftar, err := parseDate(input.TanggalDaftar)
	if err != nil {
		return nil, errors.New("invalid tanggal_daftar format, expected YYYY-MM-DD")
	}

	existing.NomorPelanggan = input.NomorPelanggan
	existing.NamaMember = input.NamaMember
	existing.NomorTelepon = input.NomorTelepon
	existing.Alamat = input.Alamat
	existing.TanggalLahir = tglLahir
	existing.TanggalDaftar = tglDaftar
	existing.IDCabang = input.IDCabang
	existing.UpdatedAt = time.Now()

	return s.repositoryMember.Update(existing)
}

func (s *serviceMember) Delete(ID int) (*models.Member, error) {
	existing, err := s.repositoryMember.Get(ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("member not found")
	}
	return s.repositoryMember.Delete(ID)
}
