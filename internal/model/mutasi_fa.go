package model

import (
	"worker/internal/abstraction"
	"worker/pkg/util/date"

	"gorm.io/gorm"
)

type MutasiFaEntity struct {
	Period    string `json:"period" validate:"required" example:"2022-01-01"`
	Versions  int    `json:"versions" validate:"required" example:"1"`
	CompanyID int    `json:"company_id" validate:"required" example:"1"`
	// FormatterID int    `json:"formatter_id" validate:"required" example:"2"`
	Status      *int    `json:"status" validate:"required"`
}

type MutasiFaFilter struct {
	Period      *string `query:"period" example:"2022-01-01" filter:"DATESTRING"`
	Versions    *int    `query:"versions" example:"1"`
	ArrVersions *[]int  `filter:"CUSTOM" example:"1"`
	// FormatterID *int    `query:"formatter_id" example:"1"`
	Status *int `query:"status" example:"1"`
}

type MutasiFaEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	MutasiFaEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	Company CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`
	// Formatter FormatterEntityModel `json:"formatter" gorm:"foreignKey:FormatterID"`
	MutasiFaDetail []MutasiFaDetailEntityModel `json:"mutasi_fa_detail" gorm:"-"`
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type MutasiFaFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	MutasiFaFilter
	CompanyCustomFilter
}

type MutasiFaVersionModel struct {
	Version []map[int]string `json:"versions"`
}

func (MutasiFaEntityModel) TableName() string {
	return "mutasi_fa"
}

func (m *MutasiFaEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *MutasiFaEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
