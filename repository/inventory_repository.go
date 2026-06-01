package repository

import (
	"annisa-api/models"
	"database/sql"
	"time"
)

type RepositoryInventory interface {
	Create(inv *models.Inventory) (*models.Inventory, error)
	GetByID(ID int) (*models.Inventory, error)
	GetByCabang(idCabang int) ([]*models.Inventory, error)
	Update(inv *models.Inventory) (*models.Inventory, error)
	Delete(ID int) (*models.Inventory, error)
	AdjustStok(ID, delta int, updatedAt time.Time) (int64, error)
}

type repositoryInventory struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) RepositoryInventory {
	return &repositoryInventory{db}
}

func (r *repositoryInventory) Create(inv *models.Inventory) (*models.Inventory, error) {
	query := `
		INSERT INTO inventory (id_cabang, nama_item, batas_bawah, batas_atas, stok, satuan, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		inv.IDCabang,
		inv.NamaItem,
		inv.BatasBawah,
		inv.BatasAtas,
		inv.Stok,
		inv.Satuan,
		inv.CreatedAt,
		inv.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	inv.IDInventory = int(id)
	return inv, nil
}

func (r *repositoryInventory) GetByID(ID int) (*models.Inventory, error) {
	query := `
		SELECT id_inventory, id_cabang, nama_item, batas_bawah, batas_atas, stok, satuan, created_at, updated_at
		FROM inventory WHERE id_inventory = ?
	`

	inv := &models.Inventory{}
	err := r.db.QueryRow(query, ID).Scan(
		&inv.IDInventory, &inv.IDCabang, &inv.NamaItem, &inv.BatasBawah,
		&inv.BatasAtas, &inv.Stok, &inv.Satuan, &inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return inv, nil
}

func (r *repositoryInventory) GetByCabang(idCabang int) ([]*models.Inventory, error) {
	query := `
		SELECT id_inventory, id_cabang, nama_item, batas_bawah, batas_atas, stok, satuan, created_at, updated_at
		FROM inventory WHERE id_cabang = ? ORDER BY nama_item ASC
	`

	rows, err := r.db.Query(query, idCabang)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Inventory
	for rows.Next() {
		inv := &models.Inventory{}
		err := rows.Scan(
			&inv.IDInventory, &inv.IDCabang, &inv.NamaItem, &inv.BatasBawah,
			&inv.BatasAtas, &inv.Stok, &inv.Satuan, &inv.CreatedAt, &inv.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, inv)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *repositoryInventory) Update(inv *models.Inventory) (*models.Inventory, error) {
	query := `
		UPDATE inventory SET
			id_cabang = ?, nama_item = ?, batas_bawah = ?, batas_atas = ?, stok = ?, satuan = ?, updated_at = ?
		WHERE id_inventory = ?
	`

	_, err := r.db.Exec(query,
		inv.IDCabang, inv.NamaItem, inv.BatasBawah, inv.BatasAtas,
		inv.Stok, inv.Satuan, inv.UpdatedAt, inv.IDInventory,
	)
	if err != nil {
		return nil, err
	}

	return inv, nil
}

func (r *repositoryInventory) Delete(ID int) (*models.Inventory, error) {
	inv, err := r.GetByID(ID)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, nil
	}

	_, err = r.db.Exec(`DELETE FROM inventory WHERE id_inventory = ?`, ID)
	if err != nil {
		return nil, err
	}

	return inv, nil
}

// AdjustStok atomically applies a relative change to stok, refusing to let it go
// below zero. Returns the number of affected rows: 0 means either the item does
// not exist or the change would make stok negative.
func (r *repositoryInventory) AdjustStok(ID, delta int, updatedAt time.Time) (int64, error) {
	res, err := r.db.Exec(
		`UPDATE inventory SET stok = stok + ?, updated_at = ? WHERE id_inventory = ? AND stok + ? >= 0`,
		delta, updatedAt, ID, delta,
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
