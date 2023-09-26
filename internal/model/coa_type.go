package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type CoaTypeEntity struct {
	Code       string `json:"code" validate:"required" example:"1000"`
	Name       string `json:"name" validate:"required" example:"Test Group"`
	CoaGroupID int    `json:"coa_group_id" validate:"required" example:"1"`
}

type CoaTypeFilter struct {
	Code          *string `query:"code" filter:"LIKE" example:"1000"`
	Name          *string `query:"name" filter:"ILIKE" example:"Test Group"`
	CoaGroupID    *int    `query:"coa_group_id" example:"1"`
	ArrCoaGroupID *[]int  `filter:"CUSTOM"`
}

type CoaTypeEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	CoaTypeEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	CoaGroup CoaGroupEntityModel `json:"coa_group" gorm:"foreignKey:CoaGroupID;"`
	Coa      []CoaEntityModel    `json:"coa" gorm:"foreignKey:CoaTypeID;"`
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type CoaTypeFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	CoaTypeFilter
}

func (CoaTypeEntityModel) TableName() string {
	return "m_coa_type"
}

func (m *CoaTypeEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *CoaTypeEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
