package model

import (
	"worker/internal/abstraction"
	"worker/pkg/util/date"

	"gorm.io/gorm"
)

type ImportedWorksheetEntity struct {
	Versions  int    `json:"versions" validate:"required"`
	CompanyID int    `json:"company_id" validate:"required"`
	Period    string `json:"period" validate:"required"`
	Note      string `json:"note" validate:"required"`
	Status      int    `json:"status" validate:"required"`
}

type ImportedWorksheetFilter struct {
	Versions  *int    `query:"versions"`
	CompanyID *int    `query:"company_id"`
	Period    *string `query:"period"`
	Note      *string `query:"note"`
	Status    *int    `query:"status"`
}

type ImportedWorksheetEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	ImportedWorksheetEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	Company CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`
	// Formatter       FormatterEntityModel         `json:"formatter" gorm:"foreignKey:FormatterID"`
	ImportedWorksheetDetail []ImportedWorksheetDetailEntityModel `json:"imported_worksheet_detail" gorm:"foreignKey:ImportedWorksheetID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
	UserRelationModel
}

type ImportedWorksheetFilterModel struct {
	// abstraction
	abstraction.Entity

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
