package model

import (
	"worker/internal/abstraction"
	"worker/pkg/util/date"

	"gorm.io/gorm"
)

type CoaEntity struct {
	Code       string `json:"code" validate:"required"`
	Name       string `json:"name" validate:"required"`
	CoaGroupId int    `json:"coa_group_id" validate:"required"`
}

type CoaFilter struct {
	Code       *string `query:"code" filter:"LIKE"`
	Name       *string `query:"name" filter:"ILIKE"`
	CoaGroupId *int    `query:"coa_group_id"`
}

type CoaEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	CoaEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	CoaGroup CoaGroupEntityModel `json:"coa_group" gorm:"foreignKey:CoaGroupId"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type CoaFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	CoaFilter
}

func (CoaEntityModel) TableName() string {
	return "m_coa"
}

func (m *CoaEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *CoaEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
