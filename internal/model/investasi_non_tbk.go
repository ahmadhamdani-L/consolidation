package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type InvestasiNonTbkEntity struct {
	Period    string `json:"period" validate:"required" example:"2022-01-01"`
	Versions  int    `json:"versions" validate:"required" example:"1"`
	CompanyID int    `json:"company_id" validate:"required" example:"1"`
	Status    int    `json:"status" validate:"required" example:"1"`
}

type InvestasiNonTbkFilter struct {
	Period      *string `query:"period" example:"2022-01-01" filter:"DATESTRING"`
	Versions    *int    `query:"versions" example:"1"`
	ArrVersions *[]int  `filter:"CUSTOM" example:"1"`
	Status      *int    `query:"status" example:"1"`
}

type InvestasiNonTbkEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	InvestasiNonTbkEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	Company               CompanyEntityModel                 `json:"company" gorm:"foreignKey:CompanyID"`
	InvestasiNonTbkDetail []InvestasiNonTbkDetailEntityModel `json:"investasi_non_tbk" gorm:"foreignKey:InvestasiNonTbkID"`
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type InvestasiNonTbkFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	InvestasiNonTbkFilter
	CompanyCustomFilter
}
type InvestasiNonTbkVersionModel struct {
	Version []map[int]string `json:"versions"`
}

func (InvestasiNonTbkEntityModel) TableName() string {
	return "investasi_non_tbk"
}

func (m *InvestasiNonTbkEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *InvestasiNonTbkEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
