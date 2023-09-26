package model

import (
	"worker-validation/internal/abstraction"
	"worker-validation/pkg/util/date"

	"gorm.io/gorm"
)

type CompanyEntity struct {
	Code     string `json:"code" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Pic      string `json:"pic" validate:"required"`
	IsActive *bool  `json:"is_active"`
}

type CompanyFilter struct {
	Code     *string `query:"code" filter:"LIKE"`
	Name     *string `query:"name" filter:"ILIKE"`
	Pic      *string `query:"pic" filter:"ILIKE"`
	IsActive *bool   `query:"is_active"`
}

type CompanyEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	CompanyEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	ParentCompanyId *int                `json:"parent_company_id"`
	ParentCompany   *CompanyEntityModel `json:"parent_company" gorm:"foreign:ParentCompanyId"`

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
