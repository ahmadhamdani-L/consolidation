package model

import (
	"worker/internal/abstraction"
	"worker/pkg/util/date"

	"gorm.io/gorm"
)

type FormatterDetailEntity struct {
	Code           string `json:"code" validate:"required"`
	Description    string `json:"name" validate:"required"`
	SortId         int    `json:"company_id" validate:"required"`
	IsCoa          *bool  `json:"is_coa" validate:"required"`
	AutoSummary    *bool  `json:"auto_summary" validate:"required"`
	FxSummary      string `json:"fx_summary" validate:"required"`
	IsTotal        *bool  `json:"is_total" validate:"required"`
	IsControl      *bool  `json:"is_control" validate:"required"`
	ControlFormula string `json:"control_formula" validate:"required"`
	FormatterID    int    `json:"formatter_id" validate:"required"`
}

type FormatterDetailFilter struct {
	Code           *string `query:"code" filter:"ILIKE"`
	Description    *string `query:"name" filter:"ILIKE"`
	SortId         *int    `query:"sort_id"`
	IsCoa          *bool   `query:"is_coa"`
	AutoSummary    *bool   `query:"auto_summary"`
	FxSummary      *string `query:"fx_summary" filter:"ILIKE"`
	IsTotal        *bool   `query:"is_total"`
	IsControl      *bool   `query:"is_control"`
	ControlFormula *string `query:"control_formula" filter:"ILIKE"`
	FormatterID    *int    `query:"formatter_id"`
}

type FormatterDetailEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	FormatterDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	Formatter FormatterEntityModel `json:"formatter" gorm:"foreignKey:FormatterID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type FormatterDetailFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	FormatterDetailFilter
}

func (FormatterDetailEntityModel) TableName() string {
	return "m_formatter_detail"
}

func (m *FormatterDetailEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *FormatterDetailEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
