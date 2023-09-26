package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type JcteEntity struct {
	TrxNumber string `json:"trx_number"`
	Note      string `json:"note" `
	CompanyID int    `json:"company_id"  `
	Period    string `json:"period" `
	ConsolidationID      int    `json:"consolidation_id"  `
	Status    int    `json:"status"`
}

type JcteFilter struct {
	TrxNumber   			*string `query:"trx_number" filter:"ILIKE"`
	Period      			*string `query:"period" filter:"DATESTRING"`
	ConsolidationID        	*int    `query:"consolidation_id"`
	// CompanyID   *int    `query:"company_id"`
	Start       			*string `query:"start"`
	End         			*string `query:"end"`
	ArrVersions 			*[]int  `filter:"CUSTOM" example:"1"`
	Status      			*int    `query:"status" example:"1"`
	Search      			*string `query:"s" filter:"CUSTOM"`
}

type JcteEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	JcteEntity

	// relations
	Consolidation ConsolidationEntityModel `json:"consolidation" gorm:"foreignKey:ConsolidationID"`
	Company      CompanyEntityModel      `json:"company" gorm:"foreignKey:CompanyID"`
	JcteDetail   []JcteDetailEntityModel `json:"jcte_detail" gorm:"foreignKey:JcteID"`
	// context
	Context *abstraction.Context `json:"-" gorm:"-"`

	UserRelationModel
}

type JcteFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	JcteFilter

	CompanyCustomFilter
}

func (JcteEntityModel) TableName() string {
	return "jcte"
}

func (m *JcteEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *JcteEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
