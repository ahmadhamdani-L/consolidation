package model

import (
	"worker-validation/internal/abstraction"
	"worker-validation/pkg/util/date"

	"gorm.io/gorm"
)

type AgingUtangPiutangEntity struct {
	Period    string `json:"period" validate:"required"`
	Versions  int    `json:"versions" validate:"required"`
	CompanyID int    `json:"company_id" validate:"required"`
	// FormatterID int    `json:"formatter_id" validate:"required"`
	Status int `json:"status"`
}

type AgingUtangPiutangFilter struct {
	Period   *string `query:"period"`
	Versions *int    `query:"versions"`
	// FormatterID *int    `query:"formatter_id"`
	Status *int `query:"status"`
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
	AgingUtangPiutangDetail []AgingUtangPiutangDetailFmtEntityModel `json:"aging_utang_piutang_detail" gorm:"-"`
	AgingUtangPiutangMEcl   []AgingUtangPiutangDetailFmtEntityModel `json:"aging_utang_piutang_mutasi_ecl" gorm:"-"`

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
	Versions []map[int]string `json:"versions"`
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
