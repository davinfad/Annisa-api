package repository

import (
	"annisa-api/models"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type RepositoryTranskasi interface {
	CreateTx(tx *sql.Tx, t *models.Transaksi) (int64, error)
	Get(ID int) (*models.Transaksi, error)
	GetAll() ([]*models.Transaksi, error)
	GetTx(tx *sql.Tx, id int) (*models.Transaksi, error)
	DeleteTx(tx *sql.Tx, id int) error
	GetByDateAndCabang(date string, idCabang int) ([]*models.Transaksi, error)
	GetMonthlyByCabang(month, year int, idCabang int) ([]*models.Transaksi, error)
	GetDraftByCabang(idCabang int) ([]*models.Transaksi, error)
	GetTotalMoneyByDateAndCabang(date string, idCabang int) (*models.TotalMoneyResult, error)
	GetTotalMoneyByMonthAndYear(month, year, idCabang int) (*models.TotalMoneyResult, error)
}

type repositoryTransaksi struct {
	db *sql.DB
}

func NewTransaksiRepository(db *sql.DB) RepositoryTranskasi {
	return &repositoryTransaksi{db}
}

func (r *repositoryTransaksi) GetTotalMoneyByDateAndCabang(date string, idCabang int) (*models.TotalMoneyResult, error) {
	query := `
		SELECT 
			COALESCE(SUM(total_harga), 0) AS total_money,
			COALESCE(SUM(CASE WHEN metode_pembayaran = 'cash' THEN total_harga ELSE 0 END), 0) AS total_cash,
			COALESCE(SUM(CASE WHEN metode_pembayaran = 'transfer' THEN total_harga ELSE 0 END), 0) AS total_transfer,
			COALESCE(COUNT(CASE WHEN metode_pembayaran = 'cash' THEN 1 ELSE NULL END), 0) AS count_cash,
			COALESCE(COUNT(CASE WHEN metode_pembayaran = 'transfer' THEN 1 ELSE NULL END), 0) AS count_transfer
		FROM transaksi 
		WHERE DATE(created_at) = ? AND id_cabang = ? AND status = 0
	`

	var result models.TotalMoneyResult
	err := r.db.QueryRow(query, date, idCabang).Scan(
		&result.TotalMoney,
		&result.TotalCash,
		&result.TotalTransfer,
		&result.CountCash,
		&result.CountTransfer,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return all fields as 0 if no data found
			return &models.TotalMoneyResult{}, nil
		}
		return nil, fmt.Errorf("failed to get total money by month and year: %w", err)
	}

	return &result, nil
}

func (r *repositoryTransaksi) GetTotalMoneyByMonthAndYear(month, year, idCabang int) (*models.TotalMoneyResult, error) {
	query := `
		SELECT 
			COALESCE(SUM(total_harga), 0) AS total_money,
			COALESCE(SUM(CASE WHEN metode_pembayaran = 'cash' THEN total_harga ELSE 0 END), 0) AS total_cash,
			COALESCE(SUM(CASE WHEN metode_pembayaran = 'transfer' THEN total_harga ELSE 0 END), 0) AS total_transfer,
			COALESCE(COUNT(CASE WHEN metode_pembayaran = 'cash' THEN 1 ELSE NULL END), 0) AS count_cash,
			COALESCE(COUNT(CASE WHEN metode_pembayaran = 'transfer' THEN 1 ELSE NULL END), 0) AS count_transfer
		FROM transaksi 
		WHERE MONTH(created_at) = ? AND YEAR(created_at) = ? AND id_cabang = ? AND status = 0
	`

	var result models.TotalMoneyResult
	err := r.db.QueryRow(query, month, year, idCabang).Scan(
		&result.TotalMoney,
		&result.TotalCash,
		&result.TotalTransfer,
		&result.CountCash,
		&result.CountTransfer,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return all fields as 0 if no data found
			return &models.TotalMoneyResult{}, nil
		}
		return nil, fmt.Errorf("failed to get total money by month and year: %w", err)
	}

	return &result, nil
}

func (r *repositoryTransaksi) GetByDateAndCabang(date string, idCabang int) ([]*models.Transaksi, error) {
	query := `
		SELECT id_transaksi, id_cabang, id_member, nama_pelanggan, nomor_telepon, total_harga, metode_pembayaran, diskon, status, created_at
		FROM transaksi
		WHERE DATE(created_at) = ? AND id_cabang = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, date, idCabang)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Transaksi
	for rows.Next() {
		var t models.Transaksi
		var idCabang sql.NullInt64
		var idMember sql.NullInt64
		var status sql.NullInt64

		err := rows.Scan(
			&t.IDTransaksi,
			&idCabang,
			&idMember,
			&t.NamaPelanggan,
			&t.NomorTelepon,
			&t.TotalHarga,
			&t.MetodePembayaran,
			&t.Diskon,
			&status,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if idCabang.Valid {
			val := int(idCabang.Int64)
			t.IDCabang = &val
		}
		if idMember.Valid {
			val := int(idMember.Int64)
			t.IDMember = &val
		}
		if status.Valid {
			val := int(status.Int64)
			t.Status = &val
		}

		result = append(result, &t)
	}

	return result, nil
}

func (r *repositoryTransaksi) GetMonthlyByCabang(month, year, idCabang int) ([]*models.Transaksi, error) {
	query := `
		SELECT id_transaksi, id_cabang, id_member, nama_pelanggan, nomor_telepon, total_harga, metode_pembayaran, status, diskon, created_at
		FROM transaksi
		WHERE MONTH(created_at) = ? AND YEAR(created_at) = ? AND id_cabang = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, month, year, idCabang)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Transaksi
	for rows.Next() {
		var t models.Transaksi
		var idCabang sql.NullInt64
		var idMember sql.NullInt64
		var status sql.NullInt64
		err := rows.Scan(
			&t.IDTransaksi, &t.IDCabang, &t.IDMember, &t.NamaPelanggan, &t.NomorTelepon,
			&t.TotalHarga, &t.MetodePembayaran, &t.Status, &t.Diskon, &t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if idCabang.Valid {
			val := int(idCabang.Int64)
			t.IDCabang = &val
		}
		if idMember.Valid {
			val := int(idMember.Int64)
			t.IDMember = &val
		}
		if status.Valid {
			val := int(status.Int64)
			t.Status = &val
		}

		result = append(result, &t)
	}

	return result, nil
}

func (r *repositoryTransaksi) GetDraftByCabang(idCabang int) ([]*models.Transaksi, error) {
	query := `
		SELECT id_transaksi, id_cabang, id_member, nama_pelanggan, nomor_telepon, total_harga, metode_pembayaran, diskon, status, created_at
		FROM transaksi
		WHERE id_cabang = ? AND status = 1 ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, idCabang)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Transaksi
	for rows.Next() {
		var t models.Transaksi
		var idCabang sql.NullInt64
		var idMember sql.NullInt64
		var status sql.NullInt64
		err := rows.Scan(
			&t.IDTransaksi, &t.IDCabang, &t.IDMember, &t.NamaPelanggan, &t.NomorTelepon,
			&t.TotalHarga, &t.MetodePembayaran, &t.Diskon, &t.Status, &t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if idCabang.Valid {
			val := int(idCabang.Int64)
			t.IDCabang = &val
		}
		if idMember.Valid {
			val := int(idMember.Int64)
			t.IDMember = &val
		}
		if status.Valid {
			val := int(status.Int64)
			t.Status = &val
		}
		result = append(result, &t)
	}

	return result, nil
}

func (r *repositoryTransaksi) CreateTx(tx *sql.Tx, t *models.Transaksi) (int64, error) {
	query := `
		INSERT INTO transaksi (nama_pelanggan, nomor_telepon, total_harga, metode_pembayaran, id_member, id_cabang, status, diskon, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := tx.Exec(query, t.NamaPelanggan, t.NomorTelepon, t.TotalHarga, t.MetodePembayaran, t.IDMember, t.IDCabang, t.Status, t.Diskon, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *repositoryTransaksi) Get(ID int) (*models.Transaksi, error) {
	var status sql.NullInt64
	query := `
		SELECT id_transaksi, id_cabang, id_member, nama_pelanggan, nomor_telepon, total_harga, metode_pembayaran, diskon, status, created_at
		FROM transaksi WHERE id_transaksi = ?
	`

	row := r.db.QueryRow(query, ID)

	t := &models.Transaksi{}
	err := row.Scan(
		&t.IDTransaksi,
		&t.IDCabang,
		&t.IDMember,
		&t.NamaPelanggan,
		&t.NomorTelepon,
		&t.TotalHarga,
		&t.MetodePembayaran,
		&t.Diskon,
		&t.Status,
		&t.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if status.Valid {
		val := int(status.Int64)
		t.Status = &val
	}

	return t, nil
}

func (r *repositoryTransaksi) GetAll() ([]*models.Transaksi, error) {
	query := `
		SELECT id_transaksi, id_cabang, id_member, nama_pelanggan, nomor_telepon, total_harga, metode_pembayaran, status, diskon, created_at
		FROM transaksi ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transaksis []*models.Transaksi
	for rows.Next() {
		t := &models.Transaksi{}
		err := rows.Scan(
			&t.IDTransaksi,
			&t.IDCabang,
			&t.IDMember,
			&t.NamaPelanggan,
			&t.NomorTelepon,
			&t.TotalHarga,
			&t.MetodePembayaran,
			&t.Status,
			&t.Diskon,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transaksis = append(transaksis, t)
	}

	return transaksis, nil
}

func (r *repositoryTransaksi) GetTx(tx *sql.Tx, id int) (*models.Transaksi, error) {
	var t models.Transaksi
	err := tx.QueryRow(`SELECT id_cabang, status, created_at FROM transaksi WHERE id_transaksi = ?`, id).
		Scan(&t.IDCabang, &t.Status, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	t.IDTransaksi = id
	return &t, nil
}

func (r *repositoryTransaksi) DeleteTx(tx *sql.Tx, id int) error {
	res, err := tx.Exec(`DELETE FROM transaksi WHERE id_transaksi = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
