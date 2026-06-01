package models

import "time"

type Inventory struct {
	IDInventory int       `json:"id_inventory"`
	IDCabang    int       `json:"id_cabang"`
	NamaItem    string    `json:"nama_item"`
	BatasBawah  int       `json:"batas_bawah"`
	BatasAtas   int       `json:"batas_atas"`
	Stok        int       `json:"stok"`
	Satuan      string    `json:"satuan"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateInventoryDTO struct {
	IDCabang   int    `json:"id_cabang" binding:"required"`
	NamaItem   string `json:"nama_item" binding:"required"`
	BatasBawah int    `json:"batas_bawah"`
	BatasAtas  int    `json:"batas_atas"`
	Stok       int    `json:"stok"`
	Satuan     string `json:"satuan" binding:"required"`
}

// AdjustStokDTO changes stok by a relative amount: positive = stock in,
// negative = stock out. Zero is rejected (nothing to adjust).
type AdjustStokDTO struct {
	Delta int `json:"delta" binding:"required"`
}
