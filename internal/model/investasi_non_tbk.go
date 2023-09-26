package model

import (
	"worker/internal/abstraction"
	"worker/pkg/util/date"

	"gorm.io/gorm"
)

type InvestasiNonTbkEntity struct {
	Period    string `json:"period" validate:"required"`
	Versions  int    `json:"versions" validate:"required"`
	CompanyID int    `json:"company_id" validate:"required"`
	Status    int    `json:"status" validate:"required"`
}

type InvestasiNonTbkFilter struct {
	Period   *string `query:"period"`
	Versions *int    `query:"versions"`
	Status   *int    `query:"status"`
}

type InvestasiNonTbkEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	InvestasiNonTbkEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	Company               CompanyEntityModel                 `json:"company" gorm:"foreignKey:CompanyID"`
	InvestasiNonTbkDetail []InvestasiNonTbkDetailEntityModel `json:"investasi_non_tbk" gorm:"foreignKey:InvestasiNonTbkID"`

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
