package model

import (
	"mcash-finance-console-core/internal/abstraction"
)

type InvestasiNonTbkDetailEntity struct {
	// FormatterBridgesID int `json:"formatter_bridges_id" validate:"required" example:"1"`
	InvestasiNonTbkID   int      `json:"investasi_non_tbk_id" validate:"required" example:"1"`
	Code                string   `json:"code" validate:"required" example:"A"`
	Description         string   `json:"description" validate:"required" example:"PT A"`
	LbrSahamOwnership   *float64 `json:"lbr_saham_ownership" validate:"required" example:"10"`
	TotalLbrSaham       *float64 `json:"total_lbr_saham" validate:"required" example:"10000.00"`
	PercentageOwnership *float64 `json:"percentage_ownership" validate:"required" example:"10"`
	HargaPar            *float64 `json:"harga_par" validate:"required" example:"10000.00"`
	TotalHargaPar       *float64 `json:"total_harga_par" validate:"required" example:"10000.00"`
	HargaBeli           *float64 `json:"harga_beli" validate:"required" example:"10000.00"`
	TotalHargaBeli      *float64 `json:"total_harga_beli" validate:"required" example:"10000.00"`
	SortID              int      `json:"sort_id" validate:"required" example:"1"`
}

type InvestasiNonTbkDetailFilter struct {
	InvestasiNonTbkID   *int     `query:"investasi_non_tbk_id" validate:"required" example:"1"`
	Code                *string  `query:"code" example:"A"`
	Description         *string  `query:"description" example:"PT A"`
	LbrSahamOwnership   *float64 `query:"lbr_saham_ownership" example:"10"`
	TotalLbrSaham       *float64 `query:"total_saham" example:"100"`
	PercentageOwnership *float64 `query:"percentage_ownership" example:"100"`
	HargaPar            *float64 `query:"harga_par" example:"10000.00"`
	TotalHargaPar       *float64 `query:"total_harga_par" example:"10000.00"`
	HargaBeli           *float64 `query:"harga_beli" example:"1000.00"`
	TotalHargaBeli      *float64 `query:"total_harga_beli" example:"1000.00"`
	SortID              *int     `query:"sort_id" example:"1"`
}

type InvestasiNonTbkDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	InvestasiNonTbkDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	InvestasiNonTbk InvestasiNonTbkEntityModel `json:"investasi_non_tbk" gorm:"foreignKey:InvestasiNonTbkID"`
	Company         CompanyEntityModel         `json:"company" gorm:"foreignKey:Code;references:Code"`

	// FormatterBridges FormatterBridgesEntityModel `json:"-" gorm:"foreignKey:FormatterBridgesID"`

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
