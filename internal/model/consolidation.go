package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type ConsolidationEntity struct {
	Period                string `json:"period" validate:"required" example:"2022-01-01"`
	Versions              int    `json:"versions" validate:"required" example:"1"`
	ConsolidationVersions int    `json:"consolidation_versions" validate:"required" example:"1"`
	CompanyID             int    `json:"company_id" validate:"required" example:"1"`
	// IsDuplicate           bool   `json:"is_duplicate"`
	Status                int    `json:"status" validate:"required" example:"1"`
}

type ConsolidationFilter struct {
	Period         *string `query:"period" example:"2022-01-01" filter:"DATESTRING"`
	Versions       *int    `query:"versions" example:"1"`
	ArrVersions    *[]int  `filter:"CUSTOM" example:"1"`
	Status         *int    `query:"status" example:"1"`
	Search         *string `query:"s" example:"Lutfi Ramadhan" filter:"CUSTOM"`
	ConsolidationID *int    `query:"id"`
	ArrStatus *[]int  `filter:"CUSTOM" example:"1"`
	ConsolidationVersions  *int `query:"consolidation_versions" json:"consolidation_versions" example:"1"`
}

type ConsolidationEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	ConsolidationEntity

	// relations
	Company CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`
	ConsolidationDetail  []ConsolidationDetailEntityModel `json:"consolidation_detail" gorm:"-"`
	ConsolidationBridge  []ConsolidationBridgeEntityModel `json:"consolidation_bridge" gorm:"-"`
	ConsolidationDetails []ConsolidationDetailFmtEntityModel `json:"consolidation_details" gorm:"-"`
	StatusString string `json:"status_string" gorm:"-"`
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type ConsolidationFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	ConsolidationFilter
	CompanyCustomFilter
}

func (ConsolidationEntityModel) TableName() string {
	return "consolidation"
}

func (m *ConsolidationEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *ConsolidationEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
