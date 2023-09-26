package model

import (
	"worker/internal/abstraction"
	"worker/pkg/util/date"

	"gorm.io/gorm"
)

type MutasiPersediaanEntity struct {
	Period    string `json:"period" validate:"required" example:"2022-01-31"`
	Versions  int    `json:"versions" validate:"required" example:"1"`
	CompanyID int    `json:"company_id" validate:"required" example:"1"`
	// FormatterID int    `json:"formatter_id" validate:"required" example:"1"`
	Status      *int    `json:"status" validate:"required"`
}

type MutasiPersediaanFilter struct {
	Period      *string `query:"period" example:"2022-01-31"`
	Versions    *int    `query:"versions" example:"1"`
	ArrVersions *[]int  `filter:"CUSTOM" example:"1"`
	// FormatterID *int    `query:"formatter_id" example:"1"`
	Status *int `query:"status" example:"1"`
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
	MutasiPersediaanDetail              []MutasiPersediaanDetailEntityModel `json:"mutasi_persediaan_detail" gorm:"-"`
	MutasiCadanganPenghapusanpersediaan []MutasiPersediaanDetailEntityModel `json:"mutasi_cadangan_penghapusan_persediaan" gorm:"-"`
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
