package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type CoaGroupEntity struct {
	Code string `json:"code" validate:"required" example:"1000"`
	Name string `json:"name" validate:"required" example:"Test Group"`
}

type CoaGroupFilter struct {
	Code *string `query:"code" filter:"LIKE" example:"1000"`
	Name *string `query:"name" filter:"ILIKE" example:"Test Group"`
}

type CoaGroupEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	CoaGroupEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	CoaType []CoaTypeEntityModel `json:"coa_type" gorm:"foreignKey:CoaGroupID;"`
	UserRelationModel

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
