package model

import (
	"mcash-finance-console-core/internal/abstraction"
)

type FormatterDetailEntity struct {
	Code           string  `json:"code" example:"pihak-berelasi"`
	Description    string  `json:"name"  example:"pihak berelasi"`
	SortID         float64 `json:"company_id" `
	IsCoa          *bool   `json:"is_coa" `
	AutoSummary    *bool   `json:"auto_summary" `
	FxSummary      string  `json:"fx_summary" `
	IsTotal        *bool   `json:"is_total" `
	IsControl      *bool   `json:"is_control" `
	IsLabel        *bool   `json:"is_label" `
	ControlFormula string  `json:"control_formula" `
	FormatterID    int     `json:"formatter_id" `
	ShowGroupCoa    *bool   `json:"show_group_coa" `
	ParentID       *int    `json:"parent_id"`
	SummaryCoa     *string `json:"summary_coa_id" `
	IsParent       *bool   `json:"is_parent" `
	IsShowView     *bool   `json:"is_show_view" example:"false"`
	IsShowExport   *bool   `json:"is_show_export" example:"false"`
	IsRecalculate  *bool   `json:"is_recalculate" example:"false"`
	CoaGroupID		*int   `json:"coa_group_id" `
	Level 			*int   `json:"level" `
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
}

type FormatterDetailEntityModel struct {
	// abstraction
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	FormatterDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	Formatter FormatterEntityModel `json:"formatter" gorm:"foreignKey:FormatterID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type FormatterDetailFilterModel struct {
	// abstraction

	// filter
	FormatterDetailFilter
}
type FormatterDetailFmtEntityModel struct {
	FormatterDetailEntityModel
	FormatterDetailID int                                `json:"formatter_detail_id"`
	ParentID          int                                `json:"parent_id"`
	AutoSummary       bool                               `json:"auto_summary"`
	IsTotal           bool                               `json:"is_total"`
	IsControl         bool                               `json:"is_control"`
	IsLabel           bool                               `json:"is_label"`
	ControlFormula    string                             `json:"control_formula"`
	SortID    		  float64                            `json:"sort_id"`
	CoaGroupID		  *int 								 `json:"coa_group_id" `
	Children          []FormatterDetailFmtEntityModel 	 `json:"children" gorm:"-"`
	ShowGroupCoa      *bool   `json:"show_group_coa"`
}
func (FormatterDetailEntityModel) TableName() string {
	return "m_formatter_detail"
}
