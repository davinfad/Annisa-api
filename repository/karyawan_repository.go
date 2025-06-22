package repository

import (
	"annisa-api/models"
	"database/sql"
	"fmt"
	"time"
)

type RepositoryKaryawan interface {
	Create(karyawan *models.Karyawan) (*models.Karyawan, error)
	GetByID(ID int) (*models.Karyawan, error)
	GetByIDCabang(IDCabang int) (*models.Karyawan, error)
	Update(karyawan *models.Karyawan) (*models.Karyawan, error)
	Delete(ID int) (*models.Karyawan, error)
	ResetDailyCommission() error
	ResetMonthlyCommission() error
	UpdateKomisi(tx *sql.Tx, idKaryawan int, komisi float64) error
	UpdateKomisiTx(tx *sql.Tx, idKaryawan int, komisi float64, isToday bool) error
}

type repositoryKaryawan struct {
	db *sql.DB
}

func NewKaryawanRepository(db *sql.DB) RepositoryKaryawan {
	return &repositoryKaryawan{db}
}

func (r *repositoryKaryawan) Create(karyawan *models.Karyawan) (*models.Karyawan, error) {
	query := `
		INSERT INTO karyawan 
		(nama_karyawan, id_cabang, nomor_telepon, alamat, komisi, komisi_harian, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()

	result, err := r.db.Exec(query,
		karyawan.NamaKaryawan,
		karyawan.IDCabang,
		karyawan.NomorTelepon,
		karyawan.Alamat,
		karyawan.Komisi,
		karyawan.KomisiHarian,
		now,
		now,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	karyawan.IDKaryawan = int(id)

	return karyawan, nil
}

func (r *repositoryKaryawan) GetByID(ID int) (*models.Karyawan, error) {
	query := `SELECT id_karyawan, nama_karyawan, id_cabang, nomor_telepon, alamat, komisi, komisi_harian, created_at, updated_at FROM karyawan WHERE id_karyawan = ?`

	row := r.db.QueryRow(query, ID)

	var karyawan models.Karyawan
	err := row.Scan(
		&karyawan.IDKaryawan,
		&karyawan.NamaKaryawan,
		&karyawan.IDCabang,
		&karyawan.NomorTelepon,
		&karyawan.Alamat,
		&karyawan.Komisi,
		&karyawan.KomisiHarian,
		&karyawan.CreatedAt,
		&karyawan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &karyawan, nil
}

func (r *repositoryKaryawan) GetByIDCabang(IDCabang int) (*models.Karyawan, error) {
	query := `SELECT id_karyawan, nama_karyawan, id_cabang, nomor_telepon, alamat, komisi, komisi_harian, created_at, updated_at FROM karyawan WHERE id_cabang = ? ORDER BY nama_karyawan ASC LIMIT 1`

	row := r.db.QueryRow(query, IDCabang)

	var karyawan models.Karyawan
	err := row.Scan(
		&karyawan.IDKaryawan,
		&karyawan.NamaKaryawan,
		&karyawan.IDCabang,
		&karyawan.NomorTelepon,
		&karyawan.Alamat,
		&karyawan.Komisi,
		&karyawan.KomisiHarian,
		&karyawan.CreatedAt,
		&karyawan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &karyawan, nil
}

func (r *repositoryKaryawan) Update(karyawan *models.Karyawan) (*models.Karyawan, error) {
	query := `
		UPDATE karyawan 
		SET nama_karyawan = ?, id_cabang = ?, nomor_telepon = ?, alamat = ?, updated_at = ? 
		WHERE id_karyawan = ?
	`

	now := time.Now()

	_, err := r.db.Exec(query,
		karyawan.NamaKaryawan,
		karyawan.IDCabang,
		karyawan.NomorTelepon,
		karyawan.Alamat,
		// karyawan.Komisi,
		// karyawan.KomisiHarian,
		now,
		karyawan.IDKaryawan,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(karyawan.IDKaryawan)
}

func (r *repositoryKaryawan) Delete(ID int) (*models.Karyawan, error) {
	karyawan, err := r.GetByID(ID)
	if err != nil || karyawan == nil {
		return nil, fmt.Errorf("karyawan not found")
	}

	query := `DELETE FROM karyawan WHERE id_karyawan = ?`
	_, err = r.db.Exec(query, ID)
	if err != nil {
		return nil, err
	}

	return karyawan, nil
}

func (r *repositoryKaryawan) ResetDailyCommission() error {
	_, err := r.db.Exec("UPDATE karyawan SET komisi_harian = 0")
	return err
}

func (r *repositoryKaryawan) ResetMonthlyCommission() error {
	_, err := r.db.Exec("UPDATE karyawan SET komisi = 0")
	return err
}

func (r *repositoryKaryawan) UpdateKomisi(tx *sql.Tx, idKaryawan int, komisi float64) error {
	_, err := tx.Exec(`
		UPDATE karyawan
		SET komisi_harian = komisi_harian + ?, komisi = komisi + ?
		WHERE id_karyawan = ?`, komisi, komisi, idKaryawan)
	return err
}

func (r *repositoryKaryawan) UpdateKomisiTx(tx *sql.Tx, idKaryawan int, komisi float64, isToday bool) error {
	if isToday {
		_, err := tx.Exec(`UPDATE karyawan SET komisi_harian = komisi_harian - ?, komisi = komisi - ? WHERE id_karyawan = ?`, komisi, komisi, idKaryawan)
		return err
	} else {
		_, err := tx.Exec(`UPDATE karyawan SET komisi = komisi - ? WHERE id_karyawan = ?`, komisi, idKaryawan)
		return err
	}
}
