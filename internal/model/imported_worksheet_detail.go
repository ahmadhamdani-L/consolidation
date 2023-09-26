package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type ImportedWorksheetDetailEntity struct {
	ImportedWorksheetID int    `json:"imported_worksheet_id" validate:"required"`
	Code                string `json:"code" validate:"required"`
	Name                string `json:"name" validate:"required"`
	Status              int    `json:"status" validate:"required"`
	Note				string `json:"note"`
	FileName			string `json:"file_name" `
	ErrMessages         string `json:"err_messages" gorm:"column:err_,messages"`
}

type ImportedWorksheetDetailFilter struct {
	ImportedWorksheetID *int `query:"imported_worksheet_id" validate:"required"`
	Status              *int `query:"status" validate:"required"`
}

type ImportedWorksheetDetailEntityModel struct {
	abstraction.EntityImportedWorksheetDetail

	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	ImportedWorksheetDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	ImportedWorksheet ImportedWorksheetEntityModel `json:"imported_worksheet_id" gorm:"foreignKey:ImportedWorksheetID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type ImportedWorksheetDetailFilterModel struct {
	// filter
	ImportedWorksheetDetailFilter
}

func (ImportedWorksheetDetailEntityModel) TableName() string {
	return "imported_worksheet_detail"
}

func (m *ImportedWorksheetDetailEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	return
}