package model

import (
	"worker/internal/abstraction"
)

type ConsolidationBridgeDetailEntity struct {
	ConsolidationBridgeID int     `json:"consolidation_bridge_id" validate:"required"`
	Code                  string  `json:"code" validate:"required"`
	Amount                float64 `json:"amount" validate:"required"`
}

type ConsolidationBridgeDetailFilter struct {
	ConsolidationBridgeID *int     `query:"consolidation_bridge_id" validate:"required"`
	Code                  *string  `query:"code"`
	Amount                *float64 `query:"amount"`
}

type ConsolidationBridgeDetailEntityModel struct {
	// abstraction
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	ConsolidationBridgeDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	// Formatter               FormatterEntityModel                 `json:"formatter" gorm:"foreignKey:FormatterID"`
	ConsolidationBridge ConsolidationBridgeEntityModel `json:"consolidation_bridge" gorm:"foreignKey:ConsolidationBridgeID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type ConsolidationBridgeDetailFilterModel struct {
	// abstraction

	// filter
	ConsolidationBridgeDetailFilter
	CompanyCustomFilter
}

func (ConsolidationBridgeDetailEntityModel) TableName() string {
	return "consolidation_bridge_detail"
}
