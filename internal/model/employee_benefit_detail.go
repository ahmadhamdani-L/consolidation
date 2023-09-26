package model

import (
	"worker-validation/internal/abstraction"
)

type EmployeeBenefitDetailEntity struct {
	FormatterBridgesID int      `json:"formatter_bridges_id" validate:"required" example:"1"`
	Code               string   `json:"code" validate:"required" example:"BELUM_JATUH_TEMPO"`
	Description        string   `json:"description" validate:"required" example:"Belum jatuh tempo"`
	Amount             *float64 `json:"amount" validate:"required" example:"10000.00"`
	SortID             int      `json:"sort_id" validate:"required" example:"1"`
	Value              string   `json:"value" validate:"required"`
	IsValue            *bool    `json:"is_value" validate:"required"`
}

type EmployeeBenefitDetailFilter struct {
	EmployeeBenefitID  *int     `query:"employee_benefit_id" validate:"required" example:"1" filter:"CUSTOM"`
	FormatterBridgesID *int     `query:"formatter_bridges_id" example:"1" filter:"CUSTOM"`
	Code               *string  `query:"code" example:"BELUM_JATUH_TEMPO"`
	Description        *string  `query:"description" filter:"ILIKE" example:"Belum jatuh tempo"`
	Amount             *float64 `query:"amount" example:"10000.00"`
	SortID             *int     `query:"sort_id" example:"1"`
	Value              *string  `query:"value" filter:"ILIKE" example:"Projected Unit Credit"`
	IsValue            *bool    `query:"is_value"`
}

type EmployeeBenefitDetailEntityModel struct {
	// abstraction
	// abstraction.Entity
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	EmployeeBenefitDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	// EmployeeBenefit EmployeeBenefitEntityModel `json:"employee_benefit" gorm:"foreignKey:FormatterBridgesID"`
	FormatterBridges FormatterBridgesEntityModel `json:"-" gorm:"foreignKey:FormatterBridgesID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type EmployeeBenefitDetailFmtEntityModel struct {
	EmployeeBenefitDetailEntityModel
	AutoSummary    *bool   `json:"auto_summary"`
	IsTotal        *bool   `json:"is_total"`
	IsControl      *bool   `json:"is_control"`
	IsLabel        *bool   `json:"is_label"`
	ControlFormula *string `json:"control_formula"`
}

type EmployeeBenefitDetailFilterModel struct {
	// abstraction
	// abstraction.Filter

	// filter
	EmployeeBenefitDetailFilter
}

func (EmployeeBenefitDetailEntityModel) TableName() string {
	return "employee_benefit_detail"
}
