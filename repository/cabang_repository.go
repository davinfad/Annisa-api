package repository

import (
	"annisa-api/models"
	"database/sql"
)

type RepositoryCabang interface {
	Create(cabang *models.Cabang) (*models.Cabang, error)
	GetByID(ID int) (*models.Cabang, error)
	GetJamOperasional(tx *sql.Tx, idCabang int) (string, string, error)
}

type cabangRepository struct {
	db *sql.DB
}

func NewCabangRepository(db *sql.DB) *cabangRepository {
	return &cabangRepository{db}
}

func (r *cabangRepository) Create(cabang *models.Cabang) (*models.Cabang, error) {
	query := "INSERT INTO cabang (id_cabang, nama_cabang, kode_cabang, jam_buka, jam_tutup, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := r.db.Exec(query, cabang.IDCabang, cabang.NamaCabang, cabang.KodeCabang, cabang.JamBuka, cabang.JamTutup, cabang.CreatedAt, cabang.UpdatedAt)
	if err != nil {
		return cabang, err
	}

	cabangID, _ := result.LastInsertId()
	cabang.IDCabang = int(cabangID)

	return cabang, nil
}

func (r *cabangRepository) GetByID(ID int) (*models.Cabang, error) {
	query := `SELECT id_cabang, nama_cabang, kode_cabang, jam_buka, jam_tutup, created_at, updated_at FROM cabang WHERE id_cabang = ?`

	l := &models.Cabang{}
	err := r.db.QueryRow(query, ID).Scan(
		&l.IDCabang, &l.NamaCabang, &l.KodeCabang, &l.JamBuka,
		&l.JamTutup, &l.CreatedAt, &l.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}

	return l, nil
}

func (r *cabangRepository) GetJamOperasional(tx *sql.Tx, idCabang int) (string, string, error) {
	var jamBuka, jamTutup string
	err := tx.QueryRow(`SELECT jam_buka, jam_tutup FROM cabang WHERE id_cabang = ? FOR UPDATE`, idCabang).
		Scan(&jamBuka, &jamTutup)
	return jamBuka, jamTutup, err
}
