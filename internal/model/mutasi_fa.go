package model

import (
	"worker-validation/internal/abstraction"
	"worker-validation/pkg/util/date"

	"gorm.io/gorm"
)

type MutasiFaEntity struct {
	Period    string `json:"period" validate:"required"`
	Versions  int    `json:"versions" validate:"required"`
	CompanyID int    `json:"company_id" validate:"required"`
	// FormatterID int    `json:"formatter_id" validate:"required"`
	Status int `json:"status" validate:"required"`
}

type MutasiFaFilter struct {
	Period   *string `query:"period"`
	Versions *int    `query:"versions"`
	// FormatterID *int    `query:"formatter_id"`
	Status *int `query:"status"`
}

type MutasiFaEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	MutasiFaEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	Company CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`
	// Formatter      FormatterEntityModel        `json:"formatter" gorm:"foreignKey:FormatterID"`
	MutasiFaCostDetail []MutasiFaDetailFmtEntityModel `json:"mutasi_fa_cost_detail" gorm:"-"`
	MutasiFaADDetail   []MutasiFaDetailFmtEntityModel `json:"mutasi_fa_accumulated_depreciation_detail" gorm:"-"`
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type MutasiFaFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	MutasiFaFilter
	CompanyCustomFilter
}

type MutasiFaVersionModel struct {
	Versions []map[int]string `json:"versions"`
}

func (MutasiFaEntityModel) TableName() string {
	return "mutasi_fa"
}

func (m *MutasiFaEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *MutasiFaEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
