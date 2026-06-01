package repository

import (
	"annisa-api/models"
	"database/sql"
	"errors"
	"time"
)

// ErrInsufficientStock is returned when an adjustment would drive stok below 0.
var ErrInsufficientStock = errors.New("insufficient stock")

type RepositoryInventory interface {
	GetByCabang(idCabang int) ([]*models.InventoryStokView, error)
	GetStok(idItem, idCabang int) (*models.Inventory, error)
	AdjustStok(idItem, idCabang, delta int, now time.Time) (*models.Inventory, error)
	SetStok(idItem, idCabang, stok int, now time.Time) (*models.Inventory, error)
}

type repositoryInventory struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) RepositoryInventory {
	return &repositoryInventory{db}
}

// GetByCabang returns the full item catalog joined with this branch's stock.
// Items with no stock row for the branch report stok 0.
func (r *repositoryInventory) GetByCabang(idCabang int) ([]*models.InventoryStokView, error) {
	query := `
		SELECT i.id_item, i.nama_item, i.satuan, i.batas_bawah, i.batas_atas, COALESCE(inv.stok, 0)
		FROM item i
		LEFT JOIN inventory inv ON inv.id_item = i.id_item AND inv.id_cabang = ?
		ORDER BY i.nama_item ASC
	`
	rows, err := r.db.Query(query, idCabang)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.InventoryStokView
	for rows.Next() {
		v := &models.InventoryStokView{IDCabang: idCabang}
		if err := rows.Scan(&v.IDItem, &v.NamaItem, &v.Satuan, &v.BatasBawah, &v.BatasAtas, &v.Stok); err != nil {
			return nil, err
		}
		v.IsLow = v.BatasBawah > 0 && v.Stok <= v.BatasBawah
		list = append(list, v)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repositoryInventory) GetStok(idItem, idCabang int) (*models.Inventory, error) {
	inv := &models.Inventory{}
	err := r.db.QueryRow(
		`SELECT id_inventory, id_item, id_cabang, stok, created_at, updated_at
		 FROM inventory WHERE id_item = ? AND id_cabang = ?`, idItem, idCabang).
		Scan(&inv.IDInventory, &inv.IDItem, &inv.IDCabang, &inv.Stok, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return inv, nil
}

// AdjustStok applies a relative change to a branch's stock inside a transaction,
// creating the stock row on first touch (upsert) and refusing to go below zero.
func (r *repositoryInventory) AdjustStok(idItem, idCabang, delta int, now time.Time) (*models.Inventory, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	var stok int
	exists := true
	err = tx.QueryRow(
		`SELECT stok FROM inventory WHERE id_item = ? AND id_cabang = ? FOR UPDATE`, idItem, idCabang).
		Scan(&stok)
	if err == sql.ErrNoRows {
		exists = false
		stok = 0
	} else if err != nil {
		return nil, err
	}

	newStok := stok + delta
	if newStok < 0 {
		return nil, ErrInsufficientStock
	}

	if exists {
		_, err = tx.Exec(
			`UPDATE inventory SET stok = ?, updated_at = ? WHERE id_item = ? AND id_cabang = ?`,
			newStok, now, idItem, idCabang)
	} else {
		_, err = tx.Exec(
			`INSERT INTO inventory (id_item, id_cabang, stok, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
			idItem, idCabang, newStok, now, now)
	}
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	committed = true

	return r.GetStok(idItem, idCabang)
}

// SetStok sets a branch's stock to an absolute value (upsert).
func (r *repositoryInventory) SetStok(idItem, idCabang, stok int, now time.Time) (*models.Inventory, error) {
	_, err := r.db.Exec(`
		INSERT INTO inventory (id_item, id_cabang, stok, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE stok = VALUES(stok), updated_at = VALUES(updated_at)`,
		idItem, idCabang, stok, now, now)
	if err != nil {
		return nil, err
	}
	return r.GetStok(idItem, idCabang)
}
