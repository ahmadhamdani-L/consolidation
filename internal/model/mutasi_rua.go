package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type MutasiRuaEntity struct {
	Period    string `json:"period" validate:"required" example:"2022-01-01"`
	Versions  int    `json:"versions" validate:"required" example:"1"`
	CompanyID int    `json:"company_id" validate:"required" example:"1"`
	// FormatterID int    `json:"formatter_id" validate:"required" example:"1"`
	Status int `json:"status" validate:"required" example:"1"`
}

type MutasiRuaFilter struct {
	Period      *string `query:"period" example:"2022-01-01" filter:"DATESTRING"`
	Versions    *int    `query:"versions" example:"1"`
	ArrVersions *[]int  `filter:"CUSTOM" example:"1"`
	// FormatterID *int    `query:"formatter_id" example:"1"`
	Status *int `query:"status" example:"1"`
}

type MutasiRuaEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	MutasiRuaEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	Company CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`
	// Formatter       FormatterEntityModel         `json:"formatter" gorm:"foreignKey:FormatterID"`
	MutasiRuaCostDetail []MutasiRuaDetailFmtEntityModel `json:"mutasi_rua_cost_detail" gorm:"-"`
	MutasiRuaADDetail   []MutasiRuaDetailFmtEntityModel `json:"mutasi_rua_accumulated_depreciation_detail" gorm:"-"`
	ControlMIACost 			[]MutasiRuaDetailEntityModel `json:"control_rua_cost" gorm:"-"`
	ControlMIAD 			[]MutasiRuaDetailEntityModel `json:"control_rua_d" gorm:"-"`
	MutasiDetailPengurangan []MutasiRuaDetailFmtEntityModel `json:"mutasi_detail_pengurangan" gorm:"-"`
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
	Version []map[int]string `json:"versions"`
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
