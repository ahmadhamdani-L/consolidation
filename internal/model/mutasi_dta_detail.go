package model

import (
	"mcash-finance-console-core/internal/abstraction"
)

type MutasiDtaDetailEntity struct {
	// MutasiDtaID         int      `json:"mutasi_dta_id" validate:"required" example:"1"`
	FormatterBridgesID  int      `json:"formatter_bridges_id" validate:"required" example:"1"`
	Code                string   `json:"code" validate:"required" example:"Lain-lain"`
	Description         string   `json:"description" validate:"required" example:"Lain-lain"`
	SaldoAwal           *float64 `json:"saldo_awal" validate:"required" example:"10000.00"`
	ManfaatBebanPajak   *float64 `json:"manfaat_beban_pajak" validate:"required" example:"10000.00"`
	Oci                 *float64 `json:"oci" validate:"required" example:"1000.00"`
	AkuisisiEntitasAnak *float64 `json:"akuisisi_entitas_anak" validate:"required" example:"10000.00"`
	DibebankanKeLr      *float64 `json:"dibebankan_ke_lr" validate:"required" example:"10000.00"`
	DibebankanKeOci     *float64 `json:"dibebankan_ke_oci" validate:"required" example:"10000.00"`
	SaldoAkhir          *float64 `json:"saldo_akhir" validate:"required" example:"10000.00"`
	SortID              int      `json:"sort_id" validate:"required" example:"1"`
}

type MutasiDtaDetailFilter struct {
	FormatterBridgesID  *int     `query:"formatter_bridges_id" example:"1"`
	MutasiDtaID         *int     `query:"mutasi_dta_id" validate:"required" example:"1" filter:"CUSTOM"`
	Code                *string  `query:"code" example:"Lain-lain"`
	Description         *string  `query:"description" example:"Lain-lain"`
	SaldoAwal           *float64 `query:"saldo_awal" example:"10000.00"`
	ManfaatBebanPajak   *float64 `query:"manfaat_beban_pajak" example:"10000.00"`
	Oci                 *float64 `query:"oci" example:"1000.00"`
	AkuisisiEntitasAnak *float64 `query:"akuisisi_entitas_anak" example:"10000.00"`
	DibebankanKeLr      *float64 `query:"dibebankan_ke_lr" example:"10000.00"`
	DibebankanKeOci     *float64 `query:"dibebankan_ke_oci" example:"10000.00"`
	SaldoAkhir          *float64 `query:"saldo_akhir" example:"10000.00"`
	SortID              *int     `query:"sort_id" example:"1"`
}

type MutasiDtaDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	MutasiDtaDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	// MutasiDta MutasiDtaEntityModel `json:"mutasi_dta" gorm:"foreignKey:MutasiDtaID"`
	FormatterBridges FormatterBridgesEntityModel `json:"-" gorm:"foreignKey:FormatterBridgesID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type MutasiDtaDetailFmtEntityModel struct {
	MutasiDtaDetailEntityModel
	AutoSummary    *bool   `json:"auto_summary"`
	IsTotal        *bool   `json:"is_total"`
	IsControl      *bool   `json:"is_control"`
	IsLabel        *bool   `json:"is_label"`
	ControlFormula *string `json:"control_formula"`
}

type MutasiDtaDetailFilterModel struct {

	// filter
	MutasiDtaDetailFilter
}

func (MutasiDtaDetailEntityModel) TableName() string {
	return "mutasi_dta_detail"
}
