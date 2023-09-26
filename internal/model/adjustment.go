package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type AdjustmentEntity struct {
	TrxNumber string `json:"trx_number"`
	Note      string `json:"note" `
	CompanyID int    `json:"company_id"  `
	Period    string `json:"period" `
	TbID      int    `json:"tb_id"  `
	Status    int    `json:"status"`
}

type AdjustmentFilter struct {
	TrxNumber   *string `query:"trx_number" filter:"ILIKE"`
	Period      *string `query:"period" filter:"DATESTRING"`
	TbID        *int    `query:"tb_id"`
	// CompanyID   *int    `query:"company_id"`
	Start       *string `query:"start"`
	End         *string `query:"end"`
	ArrVersions *[]int  `filter:"CUSTOM" example:"1"`
	Status      *int    `query:"status" example:"1"`
	Search      *string `query:"s" filter:"CUSTOM"`
}

type AdjustmentEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	AdjustmentEntity

	// relations
	TrialBalance     TrialBalanceEntityModel       `json:"trial_balance" gorm:"foreignKey:TbID"`
	Company          CompanyEntityModel            `json:"company" gorm:"foreignKey:CompanyID"`
	AdjustmentDetail []AdjustmentDetailEntityModel `json:"adjustment_detail" gorm:"foreignKey:AdjustmentID"`
	UserRelationModel
	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type AdjustmentFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	AdjustmentFilter
	CompanyCustomFilter
}

func (AdjustmentEntityModel) TableName() string {
	return "adjustment"
}

func (m *AdjustmentEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *AdjustmentEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
