package model

import (
	"worker/internal/abstraction"
	"worker/pkg/util/date"

	"gorm.io/gorm"
)

type FormatterDetailEntity struct {
	Code           string  `json:"code" validate:"required" example:"pihak-berelasi"`
	Description    string  `json:"name" validate:"required" example:"pihak berelasi"`
	SortID         float64 `json:"company_id" validate:"required" example:"1"`
	IsCoa          *bool   `json:"is_coa" validate:"required" example:"false"`
	AutoSummary    *bool   `json:"auto_summary" validate:"required" example:"false"`
	FxSummary      string  `json:"fx_summary" validate:"required" example:"110501+110502"`
	IsTotal        *bool   `json:"is_total" validate:"required" example:"false"`
	IsControl      *bool   `json:"is_control" validate:"required" example:"false"`
	IsLabel        *bool   `json:"is_label" validate:"required" example:"false"`
	ControlFormula string  `json:"control_formula" validate:"required" example:"testing-control"`
	FormatterID    int     `json:"formatter_id" validate:"required" example:"1"`
	ShowGoupCoa    *bool   `json:"show_goup_coa" validate:"required" example:"false"`
	ParentID       *int    `json:"parent_id" validate:"required" example:"1"`
	SummaryCoa     *string `json:"summary_coa_id" validate:"required" example:"1"`
	IsParent       *bool   `json:"is_parent" validate:"required" example:"false"`
	IsShowView     *bool   `json:"is_show_view" example:"false"`
	IsShowExport   *bool   `json:"is_show_export" example:"false"`
	IsRecalculate  *bool   `json:"is_recalculate" example:"false"`
	Level 			*int   `json:"level" example:"1"`
}

type FormatterDetailFilter struct {
	Code           *string  `query:"code" filter:"ILIKE" example:"pihak-berelasi"`
	Description    *string  `query:"name" filter:"ILIKE" example:"pihak berelasi"`
	SortID         *float64 `query:"sort_id" example:"1"`
	IsCoa          *bool    `query:"is_coa" example:"false"`
	AutoSummary    *bool    `query:"auto_summary" example:"false"`
	FxSummary      *string  `query:"fx_summary" filter:"ILIKE" example:"110501+110502"`
	IsTotal        *bool    `query:"is_total" example:"false"`
	IsControl      *bool    `query:"is_control" example:"false"`
	IsLabel        *bool    `query:"is_label" example:"false"`
	ControlFormula *string  `query:"control_formula" filter:"ILIKE" example:"testing-control"`
	FormatterID    *int     `query:"formatter_id" example:"1"`
	ShowGoupCoa    *bool    `query:"show_goup_coa" example:"false"`
	ParentID       *int     `query:"parent_id" example:"1"`
	SummaryCoa     *string  `query:"summary_coa_id" example:"1"`
	IsParent       *bool    `query:"is_parent" example:"false"`
	IsShowView     *bool    `query:"is_show_view" example:"false"`
	IsShowExport   *bool    `query:"is_show_export" example:"false"`
	IsRecalculate  *bool    `query:"is_recalculate" example:"false"`
	Level          *int     `query:"level" example:"1"`
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
