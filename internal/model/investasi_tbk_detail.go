package model

import (
	"mcash-finance-console-core/internal/abstraction"
)

type InvestasiTbkDetailEntity struct {
	FormatterBridgesID int `json:"formatter_bridges_id" validate:"required" example:"1"`
	// InvestasiTbkID int      `json:"investasi_tbk_id" validate:"required" example:"1"`
	Stock          string   `json:"stock" validate:"required" example:"BBCA"`
	EndingShares   *float64 `json:"ending_shares" validate:"required" example:"10000.00"`
	AvgPrice       *float64 `json:"avg_price" validate:"required" example:"1000.00"`
	AmountCost     *float64 `json:"amount_cost" validate:"required" example:"10000.00"`
	ClosingPrice   *float64 `json:"closing_price" validate:"required" example:"10000.00"`
	AmountFv       *float64 `json:"amount_fv" validate:"required" example:"10000.00"`
	UnrealizedGain *float64 `json:"unrealized_gain" validate:"required" example:"10000.00"`
	RealizedGain   *float64 `json:"realized_gain" validate:"required" example:"10000.00"`
	Fee            *float64 `json:"fee" validate:"required" example:"10000.00"`
	SortID         int      `json:"sort_id" validate:"required" example:"1"`
}

type InvestasiTbkDetailFilter struct {
	InvestasiTbkID     *int     `query:"investasi_tbk_id" validate:"required" example:"1" filter:"CUSTOM"`
	FormatterBridgesID *int     `query:"formatter_bridges_id" example:"1"`
	Stock              *string  `query:"stock" example:"BBCA"`
	EndingShares       *float64 `query:"ending_shares" example:"10000.00"`
	AvgPrice           *float64 `query:"avg_price" example:"1000.00"`
	AmountCost         *float64 `query:"amount_cost" example:"10000.00"`
	ClosingPrice       *float64 `query:"closing_price" example:"10000.00"`
	AmountFv           *float64 `query:"amount_fv" example:"10000.00"`
	UnrealizedGain     *float64 `query:"unrealized_gain" example:"10000.00"`
	RealizedGain       *float64 `query:"realized_gain" example:"10000.00"`
	Fee                *float64 `query:"fee" example:"10000.00"`
	SortID             *int     `query:"sort_id" example:"1"`
}

type InvestasiTbkDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	InvestasiTbkDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	// InvestasiTbk InvestasiTbkEntityModel `json:"investasi_tbk" gorm:"foreignKey:InvestasiTbkID"`
	FormatterBridges FormatterBridgesEntityModel `json:"-" gorm:"foreignKey:FormatterBridgesID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}
type InvestasiTbkDetailFmtEntityModel struct {
	InvestasiTbkDetailEntityModel
	AutoSummary    *bool   `json:"auto_summary"`
	IsTotal        *bool   `json:"is_total"`
	IsControl      *bool   `json:"is_control"`
	IsLabel        *bool   `json:"is_label"`
	ControlFormula *string `json:"control_formula"`
	FxSummary      *string `json:"-"`
}

type InvestasiTbkDetailFilterModel struct {

	// filter
	InvestasiTbkDetailFilter
}

func (InvestasiTbkDetailEntityModel) TableName() string {
	return "investasi_tbk_detail"
}
