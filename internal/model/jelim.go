package model

import (
	"worker-consol/internal/abstraction"
	"worker-consol/pkg/util/date"

	"gorm.io/gorm"
)

type JelimEntity struct {
	TrxNumber       string `json:"trx_number"`
	Note            string `json:"note" `
	CompanyID       int    `json:"company_id"  `
	Period          string `json:"period" `
	ConsolidationID int    `json:"consolidation_id"  `
	Status          int    `json:"status"`
}

type JelimFilter struct {
	TrxNumber       *string `query:"trx_number" filter:"ILIKE"`
	Period          *string `query:"period" filter:"DATESTRING"`
	ConsolidationID *int    `query:"consolidation_id"`
	Start           *string `query:"start"`
	End             *string `query:"end"`
	ArrVersions     *[]int  `filter:"CUSTOM" example:"1"`
	Status          *int    `query:"status" example:"1"`
	Search          *string `query:"s" filter:"CUSTOM"`
}

type JelimEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	JelimEntity

	// relations
	Company     CompanyEntityModel       `json:"company" gorm:"foreignKey:CompanyID"`
	JelimDetail []JelimDetailEntityModel `json:"jelim_detail" gorm:"foreignKey:JelimID"`
	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type JelimFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	JelimFilter
	CompanyCustomFilter
}

func (JelimEntityModel) TableName() string {
	return "jelim"
}

func (m *JelimEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *JelimEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
