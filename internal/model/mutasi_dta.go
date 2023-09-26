package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type MutasiDtaEntity struct {
	Period    string `json:"period" validate:"required" example:"2022-01-01"`
	Versions  int    `json:"versions" validate:"required" example:"1"`
	CompanyID int    `json:"company_id" validate:"required" example:"1"`
	// FormatterID int    `json:"formatter_id" validate:"required" example:"1"`
	Status int `json:"status" validate:"required" example:"1"`
}

type MutasiDtaFilter struct {
	Period      *string `query:"period" example:"2022-01-01" filter:"DATESTRING"`
	Versions    *int    `query:"versions" example:"1"`
	ArrVersions *[]int  `filter:"CUSTOM" example:"1"`
	// FormatterID *int    `query:"formatter_id" example:"1"`
	Status *int `query:"status" example:"1"`
}

type MutasiDtaEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	MutasiDtaEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	Company CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`
	// Formatter       FormatterEntityModel         `json:"formatter" gorm:"foreignKey:FormatterID"`
	MutasiDtaDetail []MutasiDtaDetailFmtEntityModel `json:"mutasi_dta_detail" gorm:"-"`
	ControlDTA		MutasiDtaDetailEntityModel `json:"control_dta" gorm:"-"`
	UserRelationModel

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
