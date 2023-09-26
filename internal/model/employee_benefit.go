package model

import (
	"worker/internal/abstraction"
	"worker/pkg/util/date"

	"gorm.io/gorm"
)

type EmployeeBenefitEntity struct {
	Period      string `json:"period" validate:"required"`
	Versions    int    `json:"versions" validate:"required"`
	CompanyID   int    `json:"company_id" validate:"required"`
	Status      *int    `json:"status" validate:"required"`
}

type EmployeeBenefitFilter struct {
	Period      *string `query:"period"`
	Versions    *int    `query:"versions"`
	Status      *int    `query:"status"`
}

type EmployeeBenefitEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	EmployeeBenefitEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	Company                CompanyEntityModel                  `json:"company" gorm:"foreignKey:CompanyID"`
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type EmployeeBenefitFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	EmployeeBenefitFilter
	CompanyCustomFilter
}

type EmployeeBenefitVersionModel struct {
	Version []map[int]string `json:"versions"`
}

func (EmployeeBenefitEntityModel) TableName() string {
	return "employee_benefit"
}

func (m *EmployeeBenefitEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *EmployeeBenefitEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
