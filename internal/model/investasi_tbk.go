package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type InvestasiTbkEntity struct {
	Period    string `json:"period" validate:"required" example:"2022-01-01"`
	Versions  int    `json:"versions" validate:"required" example:"1"`
	CompanyID int    `json:"company_id" validate:"required" example:"1"`
	// FormatterID int    `json:"formatter_id" validate:"required" example:"2"`
	Status int `json:"status" validate:"required" example:"1"`
}

type InvestasiTbkFilter struct {
	Period      *string `query:"period" example:"2022-01-01" filter:"DATESTRING"`
	Versions    *int    `query:"versions" example:"1"`
	ArrVersions *[]int  `filter:"CUSTOM" example:"1"`
	// FormatterID *int    `query:"formatter_id" example:"1"`
	Status *int `query:"status" example:"1"`
}

type InvestasiTbkEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	InvestasiTbkEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	Company CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`
	// Formatter          FormatterEntityModel            `json:"formatter" gorm:"foreignKey:FormatterID"`
	InvestasiTbkDetail []InvestasiTbkDetailFmtEntityModel `json:"investasi_tbk_detail" gorm:"-"`
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type InvestasiTbkFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	InvestasiTbkFilter
	CompanyCustomFilter
}

type InvestasiTbkVersionModel struct {
	Version []map[int]string `json:"versions"`
}

func (InvestasiTbkEntityModel) TableName() string {
	return "investasi_tbk"
}

func (m *InvestasiTbkEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *InvestasiTbkEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
