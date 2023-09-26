package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type MutasiPersediaanEntity struct {
	Period    string `json:"period" validate:"required" example:"2022-01-31"`
	Versions  int    `json:"versions" validate:"required" example:"1"`
	CompanyID int    `json:"company_id" validate:"required" example:"1"`
	// FormatterID int    `json:"formatter_id" validate:"required" example:"1"`
	Status int `json:"status" example:"1"`
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
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	Company CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`
	// Formatter              FormatterEntityModel                `json:"formatter" gorm:"foreignKey:FormatterID"`
	MutasiPersediaanDetail              []MutasiPersediaanDetailFmtEntityModel `json:"mutasi_persediaan_detail" gorm:"-"`
	MutasiCadanganPenghapusanpersediaan []MutasiPersediaanDetailFmtEntityModel `json:"mutasi_cadangan_penghapusan_persediaan" gorm:"-"`
	ControlPersediaan 			[]MutasiPersediaanDetailEntityModel `json:"control_persediaan" gorm:"-"`
	ControlPersediaanPenghapusan			[]MutasiPersediaanDetailEntityModel `json:"control_persediaan_penghapusan" gorm:"-"`
	ControlMIAD 			[]MutasiIaDetailEntityModel `json:"control_mia_d" gorm:"-"`
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
