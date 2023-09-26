package model

import (
	"worker/internal/abstraction"
)

type MutasiDtaDetailEntity struct {
	FormatterBridgesID         int      `json:"formatter_bridges" validate:"required"`
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
	FormatterBridgesIDtaID         *int     `query:"formatter_bridges" validate:"required"`
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
	FormatterBridges FormatterBridgesEntityModel `json:"formatter_bridges" gorm:"foreignKey:FormatterBridgesID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type MutasiDtaDetailFilterModel struct {

	// filter
	MutasiDtaDetailFilter
}

func (MutasiDtaDetailEntityModel) TableName() string {
	return "mutasi_dta_detail"
}
