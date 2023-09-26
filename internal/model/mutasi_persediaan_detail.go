package model

import (
	"worker/internal/abstraction"
)

type MutasiPersediaanDetailEntity struct {
	FormatterBridgesID int      `json:"formatter_bridges_id" validate:"required"`
	Code               string   `json:"code" validate:"required"`
	Description        string   `json:"description" validate:"required"`
	Amount             *float64 `json:"amount" validate:"required"`
	SortID             int      `json:"sort_id" validate:"required"`
}

type MutasiPersediaanDetailFilter struct {
	FormatterBridgesID *int     `query:"formatter_bridges_id" validate:"required"`
	Code               *string  `query:"code" filter:"ILIKE"`
	Description        *string  `query:"description" filter:"ILIKE"`
	Amount             *float64 `query:"amount"`
	SortID             *int     `query:"sort_id"`
}

type MutasiPersediaanDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	MutasiPersediaanDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	FormatterBridges FormatterBridgesEntityModel `json:"formatter_bridges" gorm:"foreignKey:FormatterBridgesID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type MutasiPersediaanDetailFilterModel struct {

	// filter
	MutasiPersediaanDetailFilter
}

func (MutasiPersediaanDetailEntityModel) TableName() string {
	return "mutasi_persediaan_detail"
}
