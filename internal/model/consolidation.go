package model

import (
	"worker/internal/abstraction"
	"worker/pkg/util/date"

	"gorm.io/gorm"
)

type ConsolidationEntity struct {
	Period                string `json:"period" validate:"required"`
	Versions              int    `json:"versions" validate:"required"`
	CompanyID             int    `json:"company_id" validate:"required"`
	ConsolidationVersions int    `json:"consolidation_versions" validate:"required"`
	Status                int    `json:"status"`
}

type ConsolidationFilter struct {
	Period                *string `query:"period"`
	Versions              *int    `query:"versions"`
	ConsolidationVersions *int    `query:"consolidation_versions"`
	// FormatterID *int    `query:"formatter_id"`
	Status *int `query:"status"`
}

type ConsolidationEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	ConsolidationEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	// Formatter               FormatterEntityModel                 `json:"formatter" gorm:"foreignKey:FormatterID"`
	Company             CompanyEntityModel               `json:"company" gorm:"foreignKey:CompanyID"`
	ConsolidationBridge []ConsolidationBridgeEntityModel `json:"consolidation_bridge" gorm:"foreignKey:ConsolidationID"`
	ConsolidationDetail []ConsolidationDetailEntityModel `json:"consolidation_detail" gorm:"foreignKey:ConsolidationID"`

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
	return
}

func (m *ConsolidationEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	return
}
