package model

import (
	"worker-consol/internal/abstraction"
	"worker-consol/pkg/util/date"

	"gorm.io/gorm"
)

type MutasiDtaEntity struct {
	Period    string `json:"period" validate:"required"`
	Versions  int    `json:"versions" validate:"required"`
	CompanyID int    `json:"company_id" validate:"required"`
	// FormatterID int    `json:"formatter_id" validate:"required"`
	Status int `json:"status" validate:"required"`
}

type MutasiDtaFilter struct {
	Period      *string `query:"period"`
	Versions    *int    `query:"versions"`
	ArrVersions *[]int  `filter:"CUSTOM" example:"1"`
	// FormatterID *int    `query:"formatter_id"`
	Status *int `query:"status"`
}

type MutasiDtaEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	MutasiDtaEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	Company CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`
	// Formatter       FormatterEntityModel         `json:"formatter" gorm:"foreignKey:FormatterID"`
	MutasiDtaDetail []MutasiDtaDetailFmtEntityModel `json:"mutasi_dta_detail" gorm:"-"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type MutasiDtaFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	MutasiDtaFilter
	CompanyCustomFilter
}

type MutasiDtaVersionModel struct {
	Version []map[int]string `json:"versions"`
}

func (MutasiDtaEntityModel) TableName() string {
	return "mutasi_dta"
}

func (m *MutasiDtaEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *MutasiDtaEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
