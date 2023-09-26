package model

import (
	"worker/internal/abstraction"
	"worker/pkg/util/date"

	"gorm.io/gorm"
)

type AgingUtangPiutangEntity struct {
	Period    string `json:"period" validate:"required" example:"2022-01-01"`
	Versions  int    `json:"versions" validate:"required" example:"1"`
	CompanyID int    `json:"company_id" validate:"required" example:"1"`
	// FormatterID int    `json:"formatter_id" validate:"required" example:"1"`
	Status      *int    `json:"status" validate:"required"`
}

type AgingUtangPiutangFilter struct {
	Period      *string `query:"period" example:"2022-01-01" filter:"DATESTRING"`
	Versions    *int    `query:"versions" example:"1"`
	ArrVersions *[]int  `filter:"CUSTOM" example:"1"`
	// FormatterID *int    `query:"formatter_id" example:"1"`
	Status *int `query:"status" example:"1"`
}

type AgingUtangPiutangEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	AgingUtangPiutangEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	Company CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`
	// Formatter               FormatterEntityModel                 `json:"formatter" gorm:"foreignKey:FormatterID"`
	AgingUtangPiutangDetail []AgingUtangPiutangDetailEntityModel `json:"aging_utang_piutang_detail" gorm:"-"`
	AgingUtangPiutangMEcl   []AgingUtangPiutangDetailEntityModel `json:"aging_utang_piutang_mutasi_ecl" gorm:"-"`
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type AgingUtangPiutangFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	AgingUtangPiutangFilter
	CompanyCustomFilter
}

type AgingUtangPiutangVersionModel struct {
	Version []map[int]string `json:"versions"`
}

func (AgingUtangPiutangEntityModel) TableName() string {
	return "aging_utang_piutang"
}

func (m *AgingUtangPiutangEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *AgingUtangPiutangEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
