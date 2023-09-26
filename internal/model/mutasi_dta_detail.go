package model

import (
	"worker/internal/abstraction"
)

type MutasiDtaDetailEntity struct {
	// MutasiDtaID         int      `json:"mutasi_dta_id" validate:"required"`
	FormatterBridgesID  int      `json:"formatter_bridges_id" validate:"required"`
	Code                string   `json:"code" validate:"required"`
	Description         string   `json:"description" validate:"required"`
	SaldoAwal           *float64 `json:"saldo_awal" validate:"required"`
	ManfaatBebanPajak   *float64 `json:"manfaat_beban_pajak" validate:"required"`
	Oci                 *float64 `json:"oci" validate:"required"`
	AkuisisiEntitasAnak *float64 `json:"akuisisi_entitas_anak" validate:"required"`
	DibebankanKeLr      *float64 `json:"dibebankan_ke_lr" validate:"required"`
	DibebankanKeOci     *float64 `json:"dibebankan_ke_oci" validate:"required"`
	SaldoAkhir          *float64 `json:"saldo_akhir" validate:"required"`
	SortId              int      `json:"sort_id" validate:"required"`
}

type MutasiDtaDetailFilter struct {
	MutasiDtaID         *int     `query:"mutasi_dta_id" validate:"required" filter:"CUSTOM"`
	FormatterBridgesID  *int     `query:"formatter_bridges_id" validate:"required"`
	Code                *string  `query:"code"`
	Description         *string  `query:"description"`
	SaldoAwal           *float64 `query:"saldo_awal"`
	ManfaatBebanPajak   *float64 `query:"manfaat_beban_pajak"`
	Oci                 *float64 `query:"oci"`
	AkuisisiEntitasAnak *float64 `query:"akuisisi_entitas_anak"`
	DibebankanKeLr      *float64 `query:"dibebankan_ke_lr"`
	DibebankanKeOci     *float64 `query:"dibebankan_ke_oci"`
	SaldoAkhir          *float64 `query:"saldo_akhir"`
	SortId              *int     `query:"sort_id"`
}

type MutasiDtaDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	MutasiDtaDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
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
