package model

import (
	"worker-validation/internal/abstraction"
	"worker-validation/pkg/util/date"

	"gorm.io/gorm"
)

type ParameterEntity struct {
	Code       string `json:"code" validate:"required"`
	DataType   string `json:"data_type" validate:"required"`
	Value      string `json:"value" validate:"required"`
	IsEditable *bool  `json:"is_editable" validate:"required"`
}

type ParameterFilter struct {
	Code       *string `json:"code" filter:"ILIKE"`
	DataType   *string `json:"data_type"`
	Value      *string `json:"value"`
	IsEditable *bool   `json:"is_editable"`
}

type ParameterEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	ParameterEntity

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
