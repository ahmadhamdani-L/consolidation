package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type ImportedWorksheetEntity struct {
	Versions  int    `json:"versions"`
	CompanyID int    `json:"company_id" form:"company_id" validate:"required"`
	Period    string `json:"period" form:"period" validate:"required"`
	Note      string `json:"note"`
	Status    int    `json:"status"`
}
type ImportedWorksheetFilter struct {
	Versions    *int    `query:"versions"`
	// CompanyID   *int    `query:"company_id"`
	Period      *string `query:"period"`
	Note        *string `query:"note"`
	Status      *int    `query:"status"`
	Succes      *int    `query:"succes"`
	Failed      *int    `query:"failed"`
	ArrVersions *[]int  `filter:"CUSTOM" example:"1"`
	Search      *string `query:"s" filter:"CUSTOM"`
}

type ImportedWorksheetEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	ImportedWorksheetEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	Company CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`
	// Formatter       FormatterEntityModel         `json:"formatter" gorm:"foreignKey:FormatterID"`
	TrialBalance TrialBalanceEntityModel `json:"trial_balance" gorm:"-"`
	ImportedWorksheetDetail []ImportedWorksheetDetailEntityModel `json:"imported_worksheet_detail" gorm:"foreignKey:ImportedWorksheetID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
	UserRelationModel
}

type ImportedWorksheetFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	ImportedWorksheetFilter

	CompanyCustomFilter
}

func (ImportedWorksheetEntityModel) TableName() string {
	return "imported_worksheet"
}

func (m *ImportedWorksheetEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *ImportedWorksheetEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
