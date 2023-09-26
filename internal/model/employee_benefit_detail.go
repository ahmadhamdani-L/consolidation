package model

import (
	"worker/internal/abstraction"
)

type EmployeeBenefitDetailEntity struct {
	FormatterBridgesID int      `json:"formatter_bridges_id" validate:"required"`
	Code               string   `json:"code" validate:"required"`
	Description        string   `json:"description" validate:"required"`
	Amount             *float64 `json:"amount" validate:"required"`
	SortID             int      `json:"sort_id" validate:"required"`
	Value				string	`json:"value" validate:"required"`
	IsValue			   bool		`json:"is_value" validate:"required"`

}

type EmployeeBenefitDetailFilter struct {
	FormatterBridgesID *int     `query:"formatter_bridges_id" validate:"required"`
	Code               *string  `query:"code" filter:"ILIKE"`
	Description        *string  `query:"description" filter:"ILIKE"`
	Amount             *float64 `query:"amount"`
	SortID             *int     `query:"sort_id"`
	Value				string	`query:"value"`
	isValue				bool	`query:"is_value"`

}

type EmployeeBenefitDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	EmployeeBenefitDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	FormatterBridges FormatterBridgesEntityModel `json:"formatter_bridges" gorm:"foreignKey:FormatterBridgesID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type EmployeeBenefitDetailFilterModel struct {

	// filter
	EmployeeBenefitDetailFilter
}

func (EmployeeBenefitDetailEntityModel) TableName() string {
	return "employee_benefit_detail"
}
