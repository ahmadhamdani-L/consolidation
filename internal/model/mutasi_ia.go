package model

import (
	"worker-validation/internal/abstraction"
	"worker-validation/pkg/util/date"

	"gorm.io/gorm"
)

type MutasiIaEntity struct {
	Period    string `json:"period" validate:"required"`
	Versions  int    `json:"versions" validate:"required"`
	CompanyID int    `json:"company_id" validate:"required"`
	// FormatterID int    `json:"formatter_id" validate:"required"`
	Status int `json:"status" validate:"required"`
}

type MutasiIaFilter struct {
	Period      *string `query:"period" example:"2022-01-01" filter:"DATESTRING"`
	Versions    *int    `query:"versions" example:"1"`
	ArrVersions *[]int  `filter:"CUSTOM" example:"1"`
	// FormatterID *int    `query:"formatter_id" example:"1"`
	Status *int `query:"status" example:"1"`
}

type MutasiIaEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	MutasiIaEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	Company CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`
	// Formatter      FormatterEntityModel        `json:"formatter" gorm:"foreignKey:FormatterID"`
	MutasiIaCostDetail []MutasiIaDetailFmtEntityModel `json:"mutasi_ia_cost_detail" gorm:"-"`
	MutasiIaADDetail   []MutasiIaDetailFmtEntityModel `json:"mutasi_ia_accumulated_depreciation_detail" gorm:"-"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type MutasiIaFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	MutasiIaFilter
	CompanyCustomFilter
}

type MutasiIaVersionModel struct {
	Versions []map[int]string `json:"versions"`
}

func (MutasiIaEntityModel) TableName() string {
	return "mutasi_ia"
}

func (m *MutasiIaEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *MutasiIaEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
