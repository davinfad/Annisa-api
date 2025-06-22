package repository

import (
	"annisa-api/models"
	"database/sql"
	"time"
)

type RepositoryUser interface {
	Create(user *models.User) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (username, password, access_code, id_cabang, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()

	_, err := r.db.Exec(
		query,
		user.Username,
		user.Password,
		user.AccessCode,
		user.IDCabang,
		now,
		now,
	)

	return user, err
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	query := `
		SELECT 
			u.username, u.password, u.access_code, u.id_cabang, u.created_at, u.updated_at,
			c.id_cabang, c.nama_cabang, c.kode_cabang, c.jam_buka, c.jam_tutup, c.created_at, c.updated_at
		FROM users u
		LEFT JOIN cabang c ON u.id_cabang = c.id_cabang
		WHERE u.username = ?
		`

	row := r.db.QueryRow(query, username)

	var user models.User
	var jamBukaStr, jamTutupStr string

	err := row.Scan(
		&user.Username,
		&user.Password,
		&user.AccessCode,
		&user.IDCabang,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Cabangs.IDCabang,
		&user.Cabangs.NamaCabang,
		&user.Cabangs.KodeCabang,
		&jamBukaStr,
		&jamTutupStr,
		&user.Cabangs.CreatedAt,
		&user.Cabangs.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Konversi waktu string ke time.Time
	layout := "15:04:05"
	user.Cabangs.JamBuka, _ = time.Parse(layout, jamBukaStr)
	user.Cabangs.JamTutup, _ = time.Parse(layout, jamTutupStr)

	return &user, nil

}
