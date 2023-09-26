package model

import (
	"worker-validation/internal/abstraction"
	"worker-validation/pkg/util/date"

	"gorm.io/gorm"
)

type MutasiRuaEntity struct {
	Period    string `json:"period" validate:"required"`
	Versions  int    `json:"versions" validate:"required"`
	CompanyID int    `json:"company_id" validate:"required"`
	// FormatterID int    `json:"formatter_id" validate:"required"`
	Status int `json:"status" validate:"required"`
}

type MutasiRuaFilter struct {
	Period   *string `query:"period"`
	Versions *int    `query:"versions"`
	// FormatterID *int    `query:"formatter_id"`
	Status *int `query:"status"`
}

type MutasiRuaEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	MutasiRuaEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	Company CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`
	// Formatter       FormatterEntityModel         `json:"formatter" gorm:"foreignKey:FormatterID"`
	MutasiRuaCostDetail []MutasiRuaDetailFmtEntityModel `json:"mutasi_rua_cost_detail" gorm:"-"`
	MutasiRuaADDetail   []MutasiRuaDetailFmtEntityModel `json:"mutasi_rua_accumulated_depreciation_detail" gorm:"-"`
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type MutasiRuaFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	MutasiRuaFilter
	CompanyCustomFilter
}

type MutasiRuaVersionModel struct {
	Versions []map[int]string `json:"versions"`
}

func (MutasiRuaEntityModel) TableName() string {
	return "mutasi_rua"
}

func (m *MutasiRuaEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *MutasiRuaEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
