package models

import "time"

// Inventory is a per-branch stock row for a catalog item.
// Unique per (id_item, id_cabang).
type Inventory struct {
	IDInventory int       `json:"id_inventory"`
	IDItem      int       `json:"id_item"`
	IDCabang    int       `json:"id_cabang"`
	Stok        int       `json:"stok"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// InventoryStokView is one catalog item joined with a single branch's stock
// (0 when the branch has no stock row yet). IsLow flags stok at/below the
// item's reorder threshold.
type InventoryStokView struct {
	IDItem     int    `json:"id_item"`
	NamaItem   string `json:"nama_item"`
	Satuan     string `json:"satuan"`
	BatasBawah int    `json:"batas_bawah"`
	BatasAtas  int    `json:"batas_atas"`
	IDCabang   int    `json:"id_cabang"`
	Stok       int    `json:"stok"`
	IsLow      bool   `json:"is_low"`
}

// AdjustStokDTO changes a branch's stok by a relative amount:
// positive = stock in, negative = stock out. Zero is rejected.
type AdjustStokDTO struct {
	Delta int `json:"delta" binding:"required"`
}

// SetStokDTO sets a branch's stok to an absolute value. Pointer so that an
// explicit 0 is accepted while a missing field is rejected.
type SetStokDTO struct {
	Stok *int `json:"stok" binding:"required"`
}
