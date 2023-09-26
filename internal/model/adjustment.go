package model

import (
	"worker-validation/internal/abstraction"
	"worker-validation/pkg/util/date"

	"gorm.io/gorm"
)

type AdjustmentEntity struct {
	TrxNumber      string `json:"trx_number"`
	Note           string `json:"note" validate:"required"`
	CompanyID      int    `json:"company_id"`
	Period         string `json:"period"`
	TrialBalanceID int    `json:"trial_balance_id"`
	Status         int    `json:"status"`
}

type AdjustmentFilter struct {
	TrxNumber      *string `query:"trx_number" filter:"ILIKE"`
	Period         *string `query:"period"`
	TrialBalanceID *int    `query:"tb_id"`
}

type AdjustmentEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	AdjustmentEntity

	// relations
	Company          CompanyEntityModel            `json:"company" gorm:"foreignKey:CompanyID"`
	AdjustmentDetail []AdjustmentDetailEntityModel `json:"adjustment_detail" gorm:"foreignKey:AdjustmentID"`
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
