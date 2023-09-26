package model

import (
	"worker/internal/abstraction"
)

type AgingUtangPiutangDetailEntity struct {
	// AgingUtangPiutangID          int      `json:"aging_utang_piutang_id" validate:"required"`
	FormatterBridgesID           int      `json:"formatter_bridges_id" validate:"required" example:"1"`
	Code                         string   `json:"code" validate:"required"`
	Description                  string   `json:"description" validate:"required"`
	Piutangusaha3rdparty         *float64 `json:"piutang_usaha_3rdparty" validate:"required" gorm:"column:piutangusaha_3rdparty"`
	PiutangusahaBerelasi         *float64 `json:"piutang_usaha_berelasi" validate:"required"`
	Piutanglainshortterm3rdparty *float64 `json:"piutang_lain_shortterm_3rdparty" validate:"required" gorm:"column:piutanglainshortterm_3rdparty"`
	PiutanglainshorttermBerelasi *float64 `json:"piutang_lain_shortterm_berelasi" validate:"required"`
	Piutangberelasishortterm     *float64 `json:"piutang_berelasi_shortterm" validate:"required"`
	Piutanglainlongterm3rdparty  *float64 `json:"piutang_lain_longterm_3rdparty" validate:"required" gorm:"column:piutanglainlongterm_3rdparty"`
	PiutanglainlongtermBerelasi  *float64 `json:"piutang_lainlongterm_berelasi" validate:"required"`
	Piutangberelasilongterm      *float64 `json:"piutang_berelasi_longterm" validate:"required"`
	Utangusaha3rdparty           *float64 `json:"utang_usaha_3rdparty" validate:"required" gorm:"column:utangusaha_3rdparty"`
	UtangusahaBerelasi           *float64 `json:"utang_usaha_berelasi" validate:"required"`
	Utanglainshortterm3rdparty   *float64 `json:"utang_lain_shortterm_3rdparty" validate:"required" gorm:"column:utanglainshortterm_3rdparty"`
	UtanglainshorttermBerelasi   *float64 `json:"utang_lain_shortterm_berelasi" validate:"required"`
	Utangberelasishortterm       *float64 `json:"utang_berelasi_short_term" validate:"required"`
	Utanglainlongterm3rdparty    *float64 `json:"utang_lain_longterm_3rdparty" validate:"required" gorm:"column:utanglainlongterm_3rdparty"`
	UtanglainlongtermBerelasi    *float64 `json:"utang_lain_longterm_berelasi" validate:"required"`
	Utangberelasilongterm        *float64 `json:"utang_berelasi_longterm" validate:"required"`
	SortID                       int      `json:"sort_id" validate:"required"`
}

type AgingUtangPiutangDetailFilter struct {
	AgingUtangPiutangID          *int     `query:"aging_utang_piutang_id" filter:"CUSTOM"`
	FormatterBridgesID           *int     `query:"formatter_bridges_id" example:"1"`
	Code                         *string  `query:"code" filter:"ILIKE"`
	Description                  *string  `query:"description" filter:"ILIKE"`
	Piutangusaha3rdparty         *float64 `query:"piutang_usaha_3rdparty"`
	PiutangusahaBerelasi         *float64 `query:"piutang_usaha_berelasi"`
	Piutanglainshortterm3rdparty *float64 `query:"piutang_lainshortterm_3rdparty"`
	PiutanglainshorttermBerelasi *float64 `query:"piutang_lainshortterm_berelasi"`
	Piutangberelasishortterm     *float64 `query:"piutang_berelasi_shortterm"`
	Piutanglainlongterm3rdparty  *float64 `query:"piutang_lain_longterm_3rdparty"`
	PiutanglainlongtermBerelasi  *float64 `query:"piutang_lainlongterm_berelasi"`
	Piutangberelasilongterm      *float64 `query:"piutang_berelasi_longterm"`
	Utangusaha3rdparty           *float64 `query:"utang_usaha_3rdparty"`
	UtangusahaBerelasi           *float64 `query:"utang_usaha_berelasi"`
	Utanglainshortterm3rdparty   *float64 `query:"utang_lain_shortterm_3rdparty"`
	UtanglainshorttermBerelasi   *float64 `query:"utang_lain_shortterm_berelasi"`
	Utangberelasishortterm       *float64 `query:"utang_berelasi_short_term"`
	Utanglainlongterm3rdparty    *float64 `query:"utang_lain_longterm_3rdparty"`
	UtanglainlongtermBerelasi    *float64 `query:"utang_lain_longterm_berelasi"`
	Utangberelasilongterm        *float64 `query:"utang_berelasi_longterm"`
	SortID                       *int     `query:"sort_id"`
}

type AgingUtangPiutangDetailEntityModel struct {
	// abstraction
	// abstraction.Entity
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	AgingUtangPiutangDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	// AgingUtangPiutang AgingUtangPiutangEntityModel `json:"aging_utang_piutang" gorm:"foreignKey:AgingUtangPiutangID"`
	FormatterBridges FormatterBridgesEntityModel `json:"-" gorm:"foreignKey:FormatterBridgesID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type AgingUtangPiutangDetailFmtEntityModel struct {
	AgingUtangPiutangDetailEntityModel
	AutoSummary    *bool   `json:"auto_summary"`
	IsTotal        *bool   `json:"is_total"`
	IsControl      *bool   `json:"is_control"`
	IsLabel        *bool   `json:"is_label"`
	ControlFormula *string `json:"control_formula"`
}

type AgingUtangPiutangDetailFilterModel struct {
	// abstraction
	// abstraction.Filter

	// filter
	AgingUtangPiutangDetailFilter
}

func (AgingUtangPiutangDetailEntityModel) TableName() string {
	return "aging_utang_piutang_detail"
}

// func (m *AgingUtangPiutangDetailEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
// 	m.CreatedAt = *date.DateTodayLocal()
// 	m.CreatedBy = m.Context.Auth.ID
// 	return
// }

// func (m *AgingUtangPiutangDetailEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
// 	m.ModifiedAt = date.DateTodayLocal()
// 	m.ModifiedBy = &m.Context.Auth.ID
// 	return
// }
