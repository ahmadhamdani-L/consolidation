package model

import (
	"worker/internal/abstraction"
	"worker/pkg/util/date"

	"gorm.io/gorm"
)

type JpmEntity struct {
	TrxNumber       string `json:"trx_number"`
	Note            string `json:"note"`
	CompanyID       int    `json:"company_id"`
	Period          string `json:"period"`
	ConsolidationID string `json:"consolidation_id"`
	Status          int    `json:"status"`
}

type JpmFilter struct {
	TrxNumber       *string `query:"trx_number" filter:"ILIKE"`
	Period          *string `query:"period" filter:"DATESTRING"`
	ConsolidationID *int    `query:"consolidation_id"`
	Start           *string `query:"start"`
	End             *string `query:"end"`
	ArrVersions     *[]int  `filter:"CUSTOM" example:"1"`
	Status          *int    `query:"status" example:"1"`
	Search          *string `query:"s" filter:"CUSTOM"`
}

type JpmEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	JpmEntity

	// relations
	Company   CompanyEntityModel     `json:"company" gorm:"foreignKey:CompanyID"`
	JpmDetail []JpmDetailEntityModel `json:"jpm_detail" gorm:"foreignKey:JpmID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type JpmFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	JpmFilter
	CompanyCustomFilter
}

func (JpmEntityModel) TableName() string {
	return "jpm"
}

func (m *JpmEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *JpmEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
