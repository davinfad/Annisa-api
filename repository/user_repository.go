package repository

import (
	"annisa-api/models"
	"database/sql"
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
		INSERT INTO users (username, password, id_cabang, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		user.Username,
		user.Password,
		user.IDCabang,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return user, err
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	query := `
		SELECT username, password, id_cabang, created_at, updated_at
		FROM users
		WHERE username = ?
	`

	row := r.db.QueryRow(query, username)

	var user models.User
	err := row.Scan(
		&user.Username,
		&user.Password,
		&user.IDCabang,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
