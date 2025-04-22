package service

import (
	"annisa-api/models"
	"annisa-api/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type ServiceTransaksi interface {
	CreateTransaksi(tx *sql.Tx, req interface{}, status int) (*models.Transaksi, error)
	UpdateKomisiKaryawan(tx *sql.Tx, items []models.ItemTransaksi, waktuTransaksi time.Time, idCabang *int) error
	GetTransaksiByID(id int) (*models.Transaksi, error)
	GetTransaksiByDateAndCabang(date string, idCabang int) ([]*models.Transaksi, error)
	GetMonthlyTransaksiByCabang(month, year, idCabang int) ([]*models.Transaksi, error)
	GetDraftTransaksiByCabang(idCabang int) ([]*models.Transaksi, error)
	DeleteTransaksi(ctx context.Context, idTransaksi int) error
	GetTotalMoneyByDateAndCabang(date string, idCabang int) (*models.TotalMoneyResult, error)
	GetTotalMoneyByMonthAndYear(month, year, idCabang int) (*models.TotalMoneyResult, error)
}

type serviceTransaksi struct {
	db                      *sql.DB
	repositoryTransaksi     repository.RepositoryTranskasi
	repositoryCabang        repository.RepositoryCabang
	repositoryItemTransaksi repository.RepositoryItemTransaksi
	repositoryLayanan       repository.RepositoryLayanan
	repositoryKaryawan      repository.RepositoryKaryawan
}

func NewTransaksiService(db *sql.DB, repositoryTransaksi repository.RepositoryTranskasi, repositoryCabang repository.RepositoryCabang, repositoryItemTransaksi repository.RepositoryItemTransaksi, repositoryLayanan repository.RepositoryLayanan, repositoryKaryawan repository.RepositoryKaryawan) ServiceTransaksi {
	return &serviceTransaksi{db, repositoryTransaksi, repositoryCabang, repositoryItemTransaksi, repositoryLayanan, repositoryKaryawan}
}

func (s *serviceTransaksi) GetTotalMoneyByDateAndCabang(date string, idCabang int) (*models.TotalMoneyResult, error) {
	return s.repositoryTransaksi.GetTotalMoneyByDateAndCabang(date, idCabang)
}

func (s *serviceTransaksi) GetTotalMoneyByMonthAndYear(month, year, idCabang int) (*models.TotalMoneyResult, error) {
	return s.repositoryTransaksi.GetTotalMoneyByMonthAndYear(month, year, idCabang)
}

// Get transaksi by ID
func (s *serviceTransaksi) GetTransaksiByID(id int) (*models.Transaksi, error) {
	get, err := s.repositoryTransaksi.Get(id)
	if err != nil {
		return nil, err
	}
	if get == nil {
		return nil, errors.New("transaksi not found")
	}
	return get, nil
}

// Get by date & cabang
func (s *serviceTransaksi) GetTransaksiByDateAndCabang(date string, idCabang int) ([]*models.Transaksi, error) {
	get, err := s.repositoryTransaksi.GetByDateAndCabang(date, idCabang)
	if err != nil {
		return nil, err
	}
	if len(get) == 0 {
		return nil, errors.New("no transaksi found for this date and cabang")
	}
	return get, nil
}

// Get monthly by cabang
func (s *serviceTransaksi) GetMonthlyTransaksiByCabang(month, year, idCabang int) ([]*models.Transaksi, error) {
	get, err := s.repositoryTransaksi.GetMonthlyByCabang(month, year, idCabang)
	if err != nil {
		return nil, err
	}
	if len(get) == 0 {
		return nil, errors.New("no transaksi found for this month, year, and cabang")
	}
	return get, nil
}

// Get draft by cabang
func (s *serviceTransaksi) GetDraftTransaksiByCabang(idCabang int) ([]*models.Transaksi, error) {
	get, err := s.repositoryTransaksi.GetDraftByCabang(idCabang)
	if err != nil {
		return nil, err
	}
	if len(get) == 0 {
		return nil, errors.New("no draft transaksi found for this cabang")
	}
	return get, nil
}

func (s *serviceTransaksi) DeleteTransaksi(ctx context.Context, idTransaksi int) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	trx, err := s.repositoryTransaksi.GetTx(tx, idTransaksi)
	if err == sql.ErrNoRows {
		return fmt.Errorf("transaction not found")
	} else if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	// Cek null dan dereference trx.Status
	if trx.Status == nil {
		return fmt.Errorf("status is nil")
	}
	if *trx.Status != 0 {
		_ = s.repositoryItemTransaksi.DeleteByTransaksiIDTx(tx, idTransaksi)
		return s.repositoryTransaksi.DeleteTx(tx, idTransaksi)
	}

	if trx.IDCabang == nil {
		return fmt.Errorf("id_cabang is nil")
	}
	jamBuka, jamTutup, err := s.repositoryCabang.GetJamOperasional(tx, *trx.IDCabang)
	if err != nil {
		return fmt.Errorf("failed to get branch working hours: %w", err)
	}

	today := time.Now().Truncate(24 * time.Hour)
	transactionDate := trx.CreatedAt.Truncate(24 * time.Hour)
	isToday := today.Equal(transactionDate)

	layout := "15:04:05"
	openingTime, _ := time.Parse(layout, jamBuka)
	closingTime, _ := time.Parse(layout, jamTutup)
	tTime := trx.CreatedAt
	isOutside := tTime.Hour() < openingTime.Hour() || tTime.Hour() > closingTime.Hour()

	items, err := s.repositoryItemTransaksi.GetByTransaksiIDTx(tx, idTransaksi)
	if err != nil {
		return fmt.Errorf("failed to get transaction items: %w", err)
	}

	for _, it := range items {
		if it.IDLayanan == nil || it.IDKaryawan == nil {
			continue
		}

		pKomisi, pLuarJam, err := s.repositoryLayanan.GetPersentaseKomisiTx(tx, *it.IDLayanan)
		if err != nil {
			continue
		}

		persen := map[bool]float64{true: pLuarJam, false: pKomisi}[isOutside]
		komisi := it.Harga * (persen / 100)

		err = s.repositoryKaryawan.UpdateKomisiTx(tx, *it.IDKaryawan, komisi, isToday)
		if err != nil {
			return fmt.Errorf("failed to update komisi: %w", err)
		}
	}

	err = s.repositoryItemTransaksi.DeleteByTransaksiIDTx(tx, idTransaksi)
	if err != nil {
		return fmt.Errorf("failed to delete item_transaksi: %w", err)
	}

	return s.repositoryTransaksi.DeleteTx(tx, idTransaksi)
}

// func deleteTransaksiItemsAndTransaksi(tx *sql.Tx, idTransaksi int) error {
// 	_, err := tx.Exec(`DELETE FROM item_transaksi WHERE id_transaksi = ?`, idTransaksi)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete item_transaksi: %w", err)
// 	}

// 	res, err := tx.Exec(`DELETE FROM transaksi WHERE id_transaksi = ?`, idTransaksi)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete transaksi: %w", err)
// 	}

// 	affected, err := res.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("failed to get affected rows: %w", err)
// 	}

// 	if affected == 0 {
// 		return fmt.Errorf("transaction not found")
// 	}

// 	return nil
// }

func (s *serviceTransaksi) CreateTransaksi(tx *sql.Tx, req interface{}, status int) (*models.Transaksi, error) {
	transaksiReq := req.(models.TransaksiRequest)

	_, _, err := s.repositoryCabang.GetJamOperasional(tx, *transaksiReq.IDCabang)
	if err != nil {
		return nil, errors.New("cabang not found or fail geting operating hours")
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	transaksi := &models.Transaksi{
		NamaPelanggan:    transaksiReq.NamaPelanggan,
		NomorTelepon:     transaksiReq.NomorTelepon,
		TotalHarga:       transaksiReq.TotalHarga,
		MetodePembayaran: transaksiReq.MetodePembayaran,
		IDMember:         transaksiReq.IDMember,
		IDCabang:         transaksiReq.IDCabang,
		Status:           &status,
		CreatedAt:        now,
	}

	idTransaksi, err := s.repositoryTransaksi.CreateTx(tx, transaksi)
	if err != nil {
		return nil, err
	}

	var items []models.ItemTransaksi
	for _, item := range transaksiReq.Items {
		items = append(items, models.ItemTransaksi{
			IDTransaksi: &[]int{int(idTransaksi)}[0],
			IDLayanan:   item.IDLayanan,
			Catatan:     item.Catatan,
			Harga:       item.Harga,
			IDKaryawan:  item.IDKaryawan,
			CreatedAt:   now,
		})
	}

	err = s.repositoryItemTransaksi.CreateBulkTx(tx, items)
	if err != nil {
		return nil, err
	}

	transaksi.IDTransaksi = int(idTransaksi)
	return transaksi, nil
}

func (s *serviceTransaksi) UpdateKomisiKaryawan(tx *sql.Tx, items []models.ItemTransaksi, waktuTransaksi time.Time, idCabang *int) error {
	jamBuka, jamTutup, err := s.repositoryCabang.GetJamOperasional(tx, *idCabang)
	if err != nil {
		return err
	}

	jam := waktuTransaksi.Format("15:04:05")

	for _, item := range items {
		persen, luarJam, err := s.repositoryLayanan.GetPersentaseKomisi(*item.IDLayanan)
		if err != nil {
			return err
		}

		isLuarJam := jam < jamBuka || jam > jamTutup
		rate := persen
		if isLuarJam {
			rate = luarJam
		}
		komisi := item.Harga * (rate / 100)

		err = s.repositoryKaryawan.UpdateKomisi(tx, *item.IDKaryawan, komisi)
		if err != nil {
			return err
		}
	}
	return nil
}
