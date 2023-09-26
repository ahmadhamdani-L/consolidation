package model

import (
	"mcash-finance-console-core/internal/abstraction"
)

type ControllerEntity struct {
	Code              string `json:"code" validate:"required"`
	Name              string `json:"name" validate:"required"`
	ControllerCommand string `json:"controller_command" validate:"required"`
	ControllerType    int    `json:"controller_type" validate:"required"`
	Coa1              string `json:"coa1" validate:"required" gorm:"column:coa_1"`
	Coa2              string `json:"coa2" validate:"required" gorm:"column:coa_2"`
	FormatterID       int    `json:"formatter_id" validate:"required"`
}

type ControllerFilter struct {
	Code              *string `query:"code" filter:"LIKE"`
	Name              *string `query:"name" filter:"ILIKE"`
	ControllerCommand *string `query:"controller_command" filter:"ILIKE"`
	ControllerType    *int    `query:"controller_type"`
	Coa1              *string `query:"coa1"`
	Coa2              *string `query:"coa2"`
	FormatterID       *int    `query:"formatter_id"`
}

type ControllerEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	ControllerEntity

	// relations

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type ControllerFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	ControllerFilter
}

func (ControllerEntityModel) TableName() string {
	return "m_controller"
}