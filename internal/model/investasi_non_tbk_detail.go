package model

import (
	"worker-validation/internal/abstraction"
)

type InvestasiNonTbkDetailEntity struct {
	InvestasiNonTbkID   int      `json:"investasi_non_tbk_id" validate:"required"`
	Code                string   `json:"code" validate:"required"`
	Description         string   `json:"description" validate:"required"`
	LbrSahamOwnership   *float64 `json:"lbr_saham_ownership" validate:"required"`
	TotalLbrSaham       *float64 `json:"total_lbr_saham" validate:"required"`
	PercentageOwnership *float64 `json:"percentage_ownership" validate:"required"`
	HargaPar            *float64 `json:"harga_par" validate:"required"`
	TotalHargaPar       *float64 `json:"total_harga_par" validate:"required"`
	HargaBeli           *float64 `json:"harga_beli" validate:"required"`
	TotalHargaBeli      *float64 `json:"total_harga_beli" validate:"required"`
	SortId              int      `json:"sort_id" validate:"required"`
}

type InvestasiNonTbkDetailFilter struct {
	InvestasiNonTbkID   *int     `query:"investasi_non_tbk_id" validate:"required"`
	Code                *string  `query:"code"`
	Description         *string  `query:"description"`
	LbrSahamOwnership   *float64 `query:"lbr_saham_ownership"`
	Oci                 *float64 `query:"oci"`
	TotalLbrSaham       *float64 `query:"total_aham"`
	PercentageOwnership *float64 `query:"percentage_ow"`
	HargaPar            *float64 `query:"harga_par"`
	TotalHargaPar       *float64 `query:"total_harga_par"`
	HargaBeli           *float64 `query:"harga_beli"`
	TotalHargaBeli      *float64 `query:"total_harga_beli"`
	SortId              *int     `query:"sort_id"`
}

type InvestasiNonTbkDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	InvestasiNonTbkDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	InvestasiNonTbk InvestasiNonTbkEntityModel `json:"investasi_non_tbk" gorm:"foreignKey:InvestasiNonTbkID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type InvestasiNonTbkDetailFilterModel struct {

	// filter
	InvestasiNonTbkDetailFilter
}

func (InvestasiNonTbkDetailEntityModel) TableName() string {
	return "investasi_non_tbk_detail"
}
