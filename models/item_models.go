package models

import "time"

// Item is the shared catalog entry, managed by the owner. Stock is tracked
// per branch in the inventory table; thresholds are global to the item.
type Item struct {
	IDItem     int       `json:"id_item"`
	NamaItem   string    `json:"nama_item"`
	Satuan     string    `json:"satuan"`
	BatasBawah int       `json:"batas_bawah"`
	BatasAtas  int       `json:"batas_atas"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateItemDTO struct {
	NamaItem   string `json:"nama_item" binding:"required"`
	Satuan     string `json:"satuan" binding:"required"`
	BatasBawah int    `json:"batas_bawah"`
	BatasAtas  int    `json:"batas_atas"`
}
