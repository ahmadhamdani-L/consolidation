package model

import (
	"mcash-finance-console-core/internal/abstraction"
)

type ConsolidationBridgeDetailEntity struct {
	ConsolidationBridgeID int     `json:"consolidation_bridge_id" validate:"required"`
	Code                  string  `json:"code" validate:"required"`
	Amount                float64 `json:"amount" validate:"required"`
}

type ConsolidationBridgeDetailFilter struct {
	Code                   *string `query:"code" filter:"ILIKE" example:"110101001"`
	ConsolidationBridgesID *int    `query:"consolidation_bridge_id" validate:"required" example:"1" filter:"CUSTOM"`
}

type ConsolidationBridgeDetailEntityModel struct {
	// abstraction
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	ConsolidationBridgeDetailEntity

	// relations
	ConsolidationBridges ConsolidationBridgeEntityModel `json:"-" gorm:"foreignKey:ConsolidationBridgeID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type ConsolidationBridgeDetailFilterModel struct {

	// filter
	ConsolidationBridgeDetailFilter
}

func (ConsolidationBridgeDetailEntityModel) TableName() string {
	return "consolidation_bridge_detail"
}
