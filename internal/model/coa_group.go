package model

import (
	"worker-consol/internal/abstraction"
	"worker-consol/pkg/util/date"

	"gorm.io/gorm"
)

type CoaGroupEntity struct {
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type CoaGroupFilter struct {
	Code *string `query:"code" filter:"LIKE"`
	Name *string `query:"name" filter:"ILIKE"`
}

type CoaGroupEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	CoaGroupEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	Coa []CoaEntityModel `json:"coa_group" gorm:"foreignKey:CoaGroupId;"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type CoaGroupFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	CoaGroupFilter
}

func (CoaGroupEntityModel) TableName() string {
	return "m_coa_group"
}

func (m *CoaGroupEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *CoaGroupEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
