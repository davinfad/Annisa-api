package repository

import (
	"annisa-api/models"
	"database/sql"
	"time"
)

type RepositoryItemTransaksi interface {
	CreateBulkTx(tx *sql.Tx, items []models.ItemTransaksi) error
	GetByTransaksiIDTx(tx *sql.Tx, idTransaksi int) ([]*models.ItemTransaksi, error)
	DeleteByTransaksiIDTx(tx *sql.Tx, idTransaksi int) error
}

type repositoryItemTransaksi struct {
	db *sql.DB
}

func NewItemTransaksiRepository(db *sql.DB) RepositoryItemTransaksi {
	return &repositoryItemTransaksi{db}
}

func (r *repositoryItemTransaksi) GetByTransaksiIDTx(tx *sql.Tx, idTransaksi int) ([]*models.ItemTransaksi, error) {
	rows, err := tx.Query(`SELECT id_layanan, harga, id_karyawan FROM item_transaksi WHERE id_transaksi = ?`, idTransaksi)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.ItemTransaksi
	for rows.Next() {
		var it models.ItemTransaksi
		if err := rows.Scan(&it.IDLayanan, &it.Harga, &it.IDKaryawan); err != nil {
			return nil, err
		}
		items = append(items, &it)
	}
	return items, nil
}

func (r *repositoryItemTransaksi) DeleteByTransaksiIDTx(tx *sql.Tx, idTransaksi int) error {
	_, err := tx.Exec(`DELETE FROM item_transaksi WHERE id_transaksi = ?`, idTransaksi)
	return err
}

func (r *repositoryItemTransaksi) CreateBulkTx(tx *sql.Tx, items []models.ItemTransaksi) error {
	query := `INSERT INTO item_transaksi (id_transaksi, id_layanan, catatan, harga, id_karyawan, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now()

	for _, item := range items {
		_, err := stmt.Exec(item.IDTransaksi, item.IDLayanan, item.Catatan, item.Harga, item.IDKaryawan, now)
		if err != nil {
			return err
		}
	}
	return nil
}
