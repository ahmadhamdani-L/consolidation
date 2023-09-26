package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type ParameterEntity struct {
	Code       string `json:"code" validate:"required" example:"AGING-UTANG-PIUTANG-IMPORT-ROW-START"`
	DataType   string `json:"data_type" validate:"required" example:"numeric"`
	Value      string `json:"value" validate:"required" example:"5"`
	IsEditable *bool  `json:"is_editable" validate:"required" example:"false"`
}

type ParameterFilter struct {
	Code       *string `query:"code" filter:"ILIKE" example:"AGING-UTANG-PIUTANG-IMPORT-ROW-START"`
	DataType   *string `json:"data_type" example:"numeric"`
	Value      *string `query:"value" filter:"ILIKE" example:"5"`
	IsEditable *bool   `json:"is_editable" example:"false"`
	Search    *string `query:"s" example:"Lutfi Ramadhan" filter:"CUSTOM"`
}

type ParameterEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	ParameterEntity
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type ParameterFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	ParameterFilter
}

func (ParameterEntityModel) TableName() string {
	return "parameters"
}

func (m *ParameterEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *ParameterEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
