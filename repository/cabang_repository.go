package repository

import (
	"annisa-api/models"
	"database/sql"
	"time"
)

type RepositoryCabang interface {
	Create(cabang *models.Cabang) (*models.Cabang, error)
	GetByID(ID int) (*models.Cabang, error)
	GetAll() ([]*models.Cabang, error)
	Update(id int, cabang *models.CabangDTO) (*models.Cabang, error)
	Delete(id int) error
	GetJamOperasional(tx *sql.Tx, idCabang int) (string, string, error)
}

type cabangRepository struct {
	db *sql.DB
}

func NewCabangRepository(db *sql.DB) *cabangRepository {
	return &cabangRepository{db}
}

func (r *cabangRepository) GetAll() ([]*models.Cabang, error) {
	query := `SELECT id_cabang, nama_cabang, kode_cabang, jam_buka, jam_tutup, created_at, updated_at FROM cabang`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cabangs []*models.Cabang
	for rows.Next() {
		var jamBukaStr, jamTutupStr string
		c := &models.Cabang{}
		err := rows.Scan(
			&c.IDCabang,
			&c.NamaCabang,
			&c.KodeCabang,
			&jamBukaStr,
			&jamTutupStr,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		layout := "15:04:05"
		c.JamBuka, _ = time.Parse(layout, jamBukaStr)
		c.JamTutup, _ = time.Parse(layout, jamTutupStr)
		cabangs = append(cabangs, c)
	}
	return cabangs, nil
}

func (r *cabangRepository) Update(id int, input *models.CabangDTO) (*models.Cabang, error) {
	query := `UPDATE cabang SET nama_cabang=?, kode_cabang=?, jam_buka=?, jam_tutup=?, updated_at=? WHERE id_cabang=?`

	_, err := r.db.Exec(query,
		input.NamaCabang,
		input.KodeCabang,
		input.JamBuka,
		input.JamTutup,
		time.Now(),
		id,
	)
	if err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

func (r *cabangRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM cabang WHERE id_cabang=?`, id)
	return err
}

func (r *cabangRepository) Create(cabang *models.Cabang) (*models.Cabang, error) {
	now := time.Now()
	query := `INSERT INTO cabang (nama_cabang, kode_cabang, jam_buka, jam_tutup, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?)`

	// ✅ Format jam_buka dan jam_tutup ke string yang cocok dengan kolom TIME di MySQL
	formattedJamBuka := cabang.JamBuka.Format("15:04:05")
	formattedJamTutup := cabang.JamTutup.Format("15:04:05")

	result, err := r.db.Exec(
		query,
		cabang.NamaCabang,
		cabang.KodeCabang,
		formattedJamBuka,
		formattedJamTutup,
		now,
		now,
	)
	if err != nil {
		return cabang, err
	}
	lastID, _ := result.LastInsertId()
	cabang.IDCabang = int(lastID)
	return cabang, nil
}

func (r *cabangRepository) GetByID(ID int) (*models.Cabang, error) {
	query := `SELECT id_cabang, nama_cabang, kode_cabang, jam_buka, jam_tutup, created_at, updated_at FROM cabang WHERE id_cabang = ?`

	var jamBukaStr, jamTutupStr string
	l := &models.Cabang{}

	err := r.db.QueryRow(query, ID).Scan(
		&l.IDCabang,
		&l.NamaCabang,
		&l.KodeCabang,
		&jamBukaStr,
		&jamTutupStr,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	layout := "15:04:05"
	l.JamBuka, _ = time.Parse(layout, jamBukaStr)
	l.JamTutup, _ = time.Parse(layout, jamTutupStr)

	return l, nil
}

func (r *cabangRepository) GetJamOperasional(tx *sql.Tx, idCabang int) (string, string, error) {
	var jamBuka, jamTutup string
	err := tx.QueryRow(`SELECT jam_buka, jam_tutup FROM cabang WHERE id_cabang = ? FOR UPDATE`, idCabang).
		Scan(&jamBuka, &jamTutup)
	return jamBuka, jamTutup, err
}
