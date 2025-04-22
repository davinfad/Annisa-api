package service

import (
	"annisa-api/models"
	"annisa-api/repository"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type ServiceKaryawan interface {
	Create(karyawan *models.CreateKaryawanDTO) (*models.Karyawan, error)
	GetByID(ID int) (*models.Karyawan, error)
	GetByIDCabang(IDCabang int) (*models.Karyawan, error)
	Update(ID int, input *models.CreateKaryawanDTO) (*models.Karyawan, error)
	Delete(ID int) (*models.Karyawan, error)
}

type serviceKaryawan struct {
	repositoryKaryawan repository.RepositoryKaryawan
}

func NewKaryawanService(repositoryKaryawan repository.RepositoryKaryawan) *serviceKaryawan {
	return &serviceKaryawan{repositoryKaryawan}
}

func (s *serviceKaryawan) Create(karyawan *models.CreateKaryawanDTO) (*models.Karyawan, error) {
	val := &models.Karyawan{
		NamaKaryawan: karyawan.NamaKaryawan,
		IDCabang:     karyawan.IDCabang,
		NomorTelepon: karyawan.NomorTelepon,
		Alamat:       karyawan.Alamat,
		Komisi:       karyawan.Komisi,
		KomisiHarian: karyawan.KomisiHarian,
	}
	newVal, err := s.repositoryKaryawan.Create(val)
	if err != nil {
		return newVal, err
	}

	return newVal, nil

}

func (s *serviceKaryawan) GetByID(ID int) (*models.Karyawan, error) {
	cek, err := s.repositoryKaryawan.GetByID(ID)
	if err != nil {
		return cek, err
	}
	if cek == nil {
		return nil, err
	}
	return cek, nil
}

func (s *serviceKaryawan) GetByIDCabang(IDCabang int) (*models.Karyawan, error) {
	get, err := s.repositoryKaryawan.GetByIDCabang(IDCabang)
	if err != nil {
		return get, err
	}
	if get == nil {
		return nil, err
	}
	return get, nil
}

func (s *serviceKaryawan) Update(ID int, input *models.CreateKaryawanDTO) (*models.Karyawan, error) {
	getID, err := s.repositoryKaryawan.GetByID(ID)
	if err != nil {
		return getID, err
	}

	if getID == nil {
		return nil, err
	}

	getID.NamaKaryawan = input.NamaKaryawan
	getID.IDCabang = input.IDCabang
	getID.NomorTelepon = input.NomorTelepon
	getID.Alamat = input.Alamat
	getID.Komisi = input.Komisi
	getID.KomisiHarian = input.KomisiHarian

	update, err := s.repositoryKaryawan.Update(getID)
	if err != nil {
		return update, err
	}
	return update, nil
}

func (s *serviceKaryawan) Delete(ID int) (*models.Karyawan, error) {
	getID, err := s.repositoryKaryawan.GetByID(ID)
	if err != nil {
		return getID, err
	}

	if getID == nil {
		return nil, err
	}

	del, err := s.repositoryKaryawan.Delete(ID)
	if err != nil {
		return del, err
	}
	return del, nil
}

func (s *serviceKaryawan) StartCommissionScheduler() {
	location := time.FixedZone("WIB", 7*3600)
	c := cron.New(cron.WithLocation(location))

	_, err := c.AddFunc("0 0 * * *", func() {
		s.resetDailyCommission()
	})
	if err != nil {
		log.Println("Gagal menjadwalkan reset komisi harian:", err)
	}

	_, err = c.AddFunc("0 0 1 * *", func() {
		s.resetMonthlyCommission()
	})
	if err != nil {
		log.Println("Gagal menjadwalkan reset komisi bulanan:", err)
	}

	c.Start()
	log.Println("Commission scheduler berjalan (harian & bulanan).")
}

func (s *serviceKaryawan) resetDailyCommission() {
	err := s.repositoryKaryawan.ResetDailyCommission()
	if err != nil {
		log.Println("Gagal reset komisi harian:", err)
	} else {
		log.Println("Berhasil reset komisi harian.")
	}
}

func (s *serviceKaryawan) resetMonthlyCommission() {
	err := s.repositoryKaryawan.ResetMonthlyCommission()
	if err != nil {
		log.Println("Gagal reset komisi bulanan:", err)
	} else {
		log.Println("Berhasil reset komisi bulanan.")
	}
}
