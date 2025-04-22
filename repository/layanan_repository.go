package repository

import (
	"annisa-api/models"
	"database/sql"
	"time"
)

type RepositoryLayanan interface {
	Create(l *models.Layanan) (*models.Layanan, error)
	GetByID(ID int) (*models.Layanan, error)
	GetAll() ([]*models.Layanan, error)
	Update(l *models.Layanan) (*models.Layanan, error)
	Delete(ID int) (*models.Layanan, error)
	GetPersentaseKomisi(idLayanan int) (float64, float64, error)
	GetPersentaseKomisiTx(tx *sql.Tx, idLayanan int) (float64, float64, error)
}

type repositoryLayanan struct {
	db *sql.DB
}

func NewLayananRepository(db *sql.DB) RepositoryLayanan {
	return &repositoryLayanan{db}
}

func (r *repositoryLayanan) GetPersentaseKomisiTx(tx *sql.Tx, idLayanan int) (float64, float64, error) {
	var persen, persenLuar float64
	err := tx.QueryRow(`SELECT persen_komisi, persen_komisi_luarjam FROM layanan WHERE id_layanan = ?`, idLayanan).
		Scan(&persen, &persenLuar)
	return persen, persenLuar, err
}

func (r *repositoryLayanan) Create(layanan *models.Layanan) (*models.Layanan, error) {
	query := `
		INSERT INTO layanan (nama_layanan, persen_komisi, persen_komisi_luarjam, kategori, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		layanan.NamaLayanan,
		layanan.PersenKomisi,
		layanan.PersenKomisiLuarJam,
		layanan.Kategori,
		layanan.CreatedAt,
		layanan.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	layanan.IDLayanan = int(id)
	return layanan, nil
}

func (r *repositoryLayanan) GetByID(ID int) (*models.Layanan, error) {
	query := `
		SELECT id_layanan, nama_layanan, persen_komisi, persen_komisi_luarjam, kategori, created_at, updated_at
		FROM layanan WHERE id_layanan = ?
	`

	l := &models.Layanan{}
	err := r.db.QueryRow(query, ID).Scan(
		&l.IDLayanan, &l.NamaLayanan, &l.PersenKomisi, &l.PersenKomisiLuarJam,
		&l.Kategori, &l.CreatedAt, &l.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}

	return l, nil
}

func (r *repositoryLayanan) GetPersentaseKomisi(idLayanan int) (float64, float64, error) {
	var persen, luarJam float64
	err := r.db.QueryRow(`SELECT persen_komisi, persen_komisi_luarjam FROM layanan WHERE id_layanan = ?`, idLayanan).Scan(&persen, &luarJam)
	return persen, luarJam, err
}

func (r *repositoryLayanan) GetAll() ([]*models.Layanan, error) {
	query := `
		SELECT id_layanan, nama_layanan, persen_komisi, persen_komisi_luarjam, kategori, created_at, updated_at
		FROM layanan ORDER BY id_layanan ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Layanan
	for rows.Next() {
		l := &models.Layanan{}
		err := rows.Scan(
			&l.IDLayanan, &l.NamaLayanan, &l.PersenKomisi, &l.PersenKomisiLuarJam,
			&l.Kategori, &l.CreatedAt, &l.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, l)
	}

	return list, nil
}

func (r *repositoryLayanan) Update(l *models.Layanan) (*models.Layanan, error) {
	query := `
		UPDATE layanan SET 
			nama_layanan = ?, 
			persen_komisi = ?, 
			persen_komisi_luarjam = ?, 
			kategori = ?, 
			updated_at = ?
		WHERE id_layanan = ?
	`

	now := time.Now()
	_, err := r.db.Exec(query,
		l.NamaLayanan, l.PersenKomisi, l.PersenKomisiLuarJam,
		l.Kategori, now, l.IDLayanan,
	)
	if err != nil {
		return nil, err
	}

	l.UpdatedAt = now
	return l, nil
}

func (r *repositoryLayanan) Delete(ID int) (*models.Layanan, error) {
	// Get data before delete (for return)
	l, err := r.GetByID(ID)
	if err != nil {
		return nil, err
	}
	if l == nil {
		return nil, nil // Not found
	}

	query := `DELETE FROM layanan WHERE id_layanan = ?`
	_, err = r.db.Exec(query, ID)
	if err != nil {
		return nil, err
	}

	return l, nil
}
