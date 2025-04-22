package repository

import (
	"annisa-api/models"
	"database/sql"
	"time"
)

type RepositoryMember interface {
	Create(member *models.Member) (*models.Member, error)
	Get(ID int) (*models.Member, error)
	GetAll() ([]*models.Member, error)
	Update(member *models.Member) (*models.Member, error)
	Delete(ID int) (*models.Member, error)
	GetMemberByIDCabang(IDCabang int) ([]*models.Member, error)
}

type repositoryMember struct {
	db *sql.DB
}

func NewMemberRepository(db *sql.DB) RepositoryMember {
	return &repositoryMember{db}
}

func (r *repositoryMember) Create(m *models.Member) (*models.Member, error) {
	query := `
		INSERT INTO member (nomor_pelanggan, nama_member, nomor_telepon, alamat, tanggal_lahir, tanggal_daftar, id_cabang, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.Exec(query,
		m.NomorPelanggan, m.NamaMember, m.NomorTelepon, m.Alamat,
		m.TanggalLahir, m.TanggalDaftar, m.IDCabang, now, now,
	)
	if err != nil {
		return nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	m.IDMember = int(lastID)
	m.CreatedAt = now
	m.UpdatedAt = now
	return m, nil
}

func (r *repositoryMember) Get(ID int) (*models.Member, error) {
	query := `
		SELECT id_member, nomor_pelanggan, nama_member, nomor_telepon, alamat, tanggal_lahir, tanggal_daftar, id_cabang, created_at, updated_at
		FROM member WHERE id_member = ?
	`

	m := &models.Member{}
	err := r.db.QueryRow(query, ID).Scan(
		&m.IDMember, &m.NomorPelanggan, &m.NamaMember, &m.NomorTelepon,
		&m.Alamat, &m.TanggalLahir, &m.TanggalDaftar, &m.IDCabang,
		&m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return m, nil
}

func (r *repositoryMember) GetAll() ([]*models.Member, error) {
	query := `
		SELECT id_member, nomor_pelanggan, nama_member, nomor_telepon, alamat, tanggal_lahir, tanggal_daftar, id_cabang, created_at, updated_at
		FROM member ORDER BY nama_member ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.Member
	for rows.Next() {
		m := &models.Member{}
		err := rows.Scan(
			&m.IDMember, &m.NomorPelanggan, &m.NamaMember, &m.NomorTelepon,
			&m.Alamat, &m.TanggalLahir, &m.TanggalDaftar, &m.IDCabang,
			&m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, m)
	}

	return members, nil
}

func (r *repositoryMember) Update(m *models.Member) (*models.Member, error) {
	query := `
		UPDATE member SET
			nomor_pelanggan = ?, nama_member = ?, nomor_telepon = ?, alamat = ?,
			tanggal_lahir = ?, tanggal_daftar = ?, id_cabang = ?, updated_at = ?
		WHERE id_member = ?
	`

	now := time.Now()
	_, err := r.db.Exec(query,
		m.NomorPelanggan, m.NamaMember, m.NomorTelepon, m.Alamat,
		m.TanggalLahir, m.TanggalDaftar, m.IDCabang, now, m.IDMember,
	)
	if err != nil {
		return nil, err
	}

	m.UpdatedAt = now
	return m, nil
}

func (r *repositoryMember) Delete(ID int) (*models.Member, error) {
	member, err := r.Get(ID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, nil
	}

	query := `DELETE FROM member WHERE id_member = ?`
	_, err = r.db.Exec(query, ID)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (r *repositoryMember) GetMemberByIDCabang(IDCabang int) ([]*models.Member, error) {
	query := `
		SELECT id_member, nomor_pelanggan, nama_member, nomor_telepon, alamat, tanggal_lahir, tanggal_daftar, id_cabang, created_at, updated_at
		FROM member WHERE id_cabang = ? ORDER BY nama_member ASC
	`

	rows, err := r.db.Query(query, IDCabang)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.Member
	for rows.Next() {
		m := &models.Member{}
		err := rows.Scan(
			&m.IDMember, &m.NomorPelanggan, &m.NamaMember, &m.NomorTelepon,
			&m.Alamat, &m.TanggalLahir, &m.TanggalDaftar, &m.IDCabang,
			&m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, m)
	}

	return members, nil
}
