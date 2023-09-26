package model

import (
	"worker-consol/internal/abstraction"
	"worker-consol/pkg/util/date"

	"gorm.io/gorm"
)

type InvestasiTbkEntity struct {
	Period    string `json:"period" validate:"required"`
	Versions  int    `json:"versions" validate:"required"`
	CompanyID int    `json:"company_id" validate:"required"`
	Status    int    `json:"status" validate:"required"`
	// FormatterID int    `json:"formatter_id" validate:"required"`
}

type InvestasiTbkFilter struct {
	Period   *string `query:"period"`
	Versions *int    `query:"versions"`
	Status   *int    `query:"status"`
	// FormatterID *int    `query:"formatter_id"`
}

type InvestasiTbkEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	InvestasiTbkEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	Company CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`
	// Formatter          FormatterEntityModel            `json:"formatter" gorm:"foreignKey:FormatterID"`
	InvestasiTbkDetail []InvestasiTbkDetailFmtEntityModel `json:"investasi_tbk_detail" gorm:"-"`

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
