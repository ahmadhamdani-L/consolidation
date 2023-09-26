package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type CompanyEntity struct {
	Code            string `json:"code" validate:"required" example:"abc"`
	Name            string `json:"name" validate:"required" example:"PT ABC"`
	Pic             string `json:"pic" example:"lutfi ramadhan"`
	ParentCompanyID *int   `json:"parent_company_id" example:"1"`
	IsActive        *bool  `json:"is_active" example:"true"`
}

type CompanyFilter struct {
	ID				*int	`query:"id" example:"1"`
	Code            *string `query:"code" filter:"LIKE" example:"abc"`
	Name            *string `query:"name" filter:"ILIKE" example:"PT ABC"`
	Pic             *string `query:"pic" filter:"ILIKE" example:"lutfi"`
	ParentCompanyID *int    `query:"parent_company_id" example:"1"`
	IsActive        *bool   `query:"is_active" example:"false"`
	WithChild       *bool   `query:"with_child" example:"true" filter:"CUSTOM"`
}

type CompanyEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	CompanyEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	ParentCompany *CompanyEntityModel  `json:"parent_company" gorm:"foreignKey:ParentCompanyID"`
	ChildCompany  []CompanyEntityModel `json:"child_company" gorm:"foreignKey:ParentCompanyID"`
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type CompanyFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	CompanyFilter
}

func (CompanyEntityModel) TableName() string {
	return "m_company"
}

func (m *CompanyEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *CompanyEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
