package model

import (
	"worker/internal/abstraction"
)

type InvestasiTbkDetailEntity struct {
	// InvestasiTbkID     int      `json:"investasi_tbk_id" validate:"required"`
	FormatterBridgesID int      `json:"formatter_bridges_id" validate:"required"`
	Stock              string   `json:"stock" validate:"required"`
	EndingShares       *float64 `json:"ending_shares" validate:"required"`
	AvgPrice           *float64 `json:"avg_price" validate:"required"`
	AmountCost         *float64 `json:"amount_cost" validate:"required"`
	ClosingPrice       *float64 `json:"closing_price" validate:"required"`
	AmountFv           *float64 `json:"amount_fv" validate:"required"`
	UnrealizedGain     *float64 `json:"unrealized_gain" validate:"required"`
	RealizedGain       *float64 `json:"realized_gain" validate:"required"`
	Fee                *float64 `json:"fee" validate:"required"`
	SortId             int      `json:"sort_id" validate:"required"`
}

type InvestasiTbkDetailFilter struct {
	FormatterBridgesID *int     `query:"formatter_bridges_id"`
	InvestasiTbkID     *int     `query:"investasi_tbk_id" validate:"required" filter:"CUSTOM"`
	Stock              *string  `query:"stock"`
	EndingShares       *float64 `query:"ending_shares"`
	AvgPrice           *float64 `query:"avg_price"`
	Oci                *float64 `query:"oci"`
	AmountCost         *float64 `query:"amount_cost"`
	ClosingPrice       *float64 `query:"closing_price"`
	AmountFv           *float64 `query:"amount_fv"`
	UnrealizedGain     *float64 `query:"unrealized_gain"`
	RealizedGain       *float64 `query:"realized_gain"`
	Fee                *float64 `query:"fee"`
	SortId             *int     `query:"sort_id"`
}

type InvestasiTbkDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	InvestasiTbkDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
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
}

type InvestasiTbkDetailFilterModel struct {

	// filter
	InvestasiTbkDetailFilter
}

func (InvestasiTbkDetailEntityModel) TableName() string {
	return "investasi_tbk_detail"
}
