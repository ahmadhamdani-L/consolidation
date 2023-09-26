package model

import (
	"worker-consol/internal/abstraction"
)

type ConsolidationBridgeEntity struct {
	ConsolidationID       int    `json:"consolidation_id" validate:"required"`
	CompanyID             int    `json:"company_id" validate:"required"`
	Versions              int    `json:"versions" validate:"required"`
	ConsolidationVersions int    `json:"consolidation_versions" validate:"required"`
	Period                string `json:"period" validate:"required"`
}

type ConsolidationBridgeFilter struct {
	ConsolidationID       *int    `query:"consolidation_id" validate:"required"`
	CompanyID             *int    `query:"company_id"`
	Versions              *int    `query:"versions"`
	ConsolidationVersions *int    `query:"consolidation_versions"`
	Period                *string `query:"period"`
}

type ConsolidationBridgeEntityModel struct {
	// abstraction
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	ConsolidationBridgeEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	// Formatter               FormatterEntityModel                 `json:"formatter" gorm:"foreignKey:FormatterID"`
	Company                   CompanyEntityModel                     `json:"company" gorm:"foreignKey:CompanyID"`
	Consolidation             ConsolidationEntityModel               `json:"consolidation" gorm:"foreignKey:ConsolidationID"`
	ConsolidationBridgeDetail []ConsolidationBridgeDetailEntityModel `json:"consolidation_detail" gorm:"foreignKey:ConsolidationBridgeID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type ConsolidationBridgeFilterModel struct {
	// abstraction

	// filter
	ConsolidationBridgeFilter
	CompanyCustomFilter
}

func (ConsolidationBridgeEntityModel) TableName() string {
	return "consolidation_bridge"
}
