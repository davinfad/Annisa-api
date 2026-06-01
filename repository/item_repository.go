package repository

import (
	"annisa-api/models"
	"database/sql"
)

type RepositoryItem interface {
	Create(item *models.Item) (*models.Item, error)
	GetByID(ID int) (*models.Item, error)
	GetAll() ([]*models.Item, error)
	Update(item *models.Item) (*models.Item, error)
	Delete(ID int) (*models.Item, error)
}

type repositoryItem struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) RepositoryItem {
	return &repositoryItem{db}
}

func (r *repositoryItem) Create(item *models.Item) (*models.Item, error) {
	query := `
		INSERT INTO item (nama_item, satuan, batas_bawah, batas_atas, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		item.NamaItem, item.Satuan, item.BatasBawah, item.BatasAtas, item.CreatedAt, item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	item.IDItem = int(id)
	return item, nil
}

func (r *repositoryItem) GetByID(ID int) (*models.Item, error) {
	query := `
		SELECT id_item, nama_item, satuan, batas_bawah, batas_atas, created_at, updated_at
		FROM item WHERE id_item = ?
	`
	item := &models.Item{}
	err := r.db.QueryRow(query, ID).Scan(
		&item.IDItem, &item.NamaItem, &item.Satuan, &item.BatasBawah,
		&item.BatasAtas, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (r *repositoryItem) GetAll() ([]*models.Item, error) {
	query := `
		SELECT id_item, nama_item, satuan, batas_bawah, batas_atas, created_at, updated_at
		FROM item ORDER BY nama_item ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Item
	for rows.Next() {
		item := &models.Item{}
		err := rows.Scan(
			&item.IDItem, &item.NamaItem, &item.Satuan, &item.BatasBawah,
			&item.BatasAtas, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repositoryItem) Update(item *models.Item) (*models.Item, error) {
	query := `
		UPDATE item SET nama_item = ?, satuan = ?, batas_bawah = ?, batas_atas = ?, updated_at = ?
		WHERE id_item = ?
	`
	_, err := r.db.Exec(query,
		item.NamaItem, item.Satuan, item.BatasBawah, item.BatasAtas, item.UpdatedAt, item.IDItem,
	)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *repositoryItem) Delete(ID int) (*models.Item, error) {
	item, err := r.GetByID(ID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	if _, err := r.db.Exec(`DELETE FROM item WHERE id_item = ?`, ID); err != nil {
		return nil, err
	}
	return item, nil
}
