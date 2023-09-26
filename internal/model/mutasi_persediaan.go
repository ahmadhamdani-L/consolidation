package model

import (
	"worker-consol/internal/abstraction"
	"worker-consol/pkg/util/date"

	"gorm.io/gorm"
)

type MutasiPersediaanEntity struct {
	Period    string `json:"period" validate:"required"`
	Versions  int    `json:"versions" validate:"required"`
	CompanyID int    `json:"company_id" validate:"required"`
	// FormatterID int    `json:"formatter_id" validate:"required"`
	Status int `json:"status" validate:"required"`
}

type MutasiPersediaanFilter struct {
	Period      *string `query:"period"`
	Versions    *int    `query:"versions"`
	ArrVersions *[]int  `filter:"CUSTOM" example:"1"`
	// FormatterID *int    `query:"formatter_id"`
	Status *int `query:"status"`
}

type MutasiPersediaanEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	MutasiPersediaanEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	Company CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`
	// Formatter              FormatterEntityModel                `json:"formatter" gorm:"foreignKey:FormatterID"`
	MutasiPersediaanDetail              []MutasiPersediaanDetailFmtEntityModel `json:"mutasi_persediaan_detail" gorm:"-"`
	MutasiCadanganPenghapusanpersediaan []MutasiPersediaanDetailFmtEntityModel `json:"mutasi_cadangan_penghapusan_persediaan" gorm:"-"`
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type MutasiPersediaanFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	MutasiPersediaanFilter
	CompanyCustomFilter
}

type MutasiPersediaanVersionModel struct {
	Version []map[int]string `json:"versions"`
}

func (MutasiPersediaanEntityModel) TableName() string {
	return "mutasi_persediaan"
}

func (m *MutasiPersediaanEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *MutasiPersediaanEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
