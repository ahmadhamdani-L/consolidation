package model

import (
	"mcash-finance-console-core/internal/abstraction"
)

type AgingUtangPiutangDetailEntity struct {
	// AgingUtangPiutangID          int      `json:"aging_utang_piutang_id" validate:"required" example:"1"`
	FormatterBridgesID           int      `json:"formatter_bridges_id" validate:"required" example:"1"`
	Code                         string   `json:"code" validate:"required" example:"BELUM_JATUH_TEMPO"`
	Description                  string   `json:"description" validate:"required" example:"Belum jatuh tempo"`
	Piutangusaha3rdparty         *float64 `json:"piutang_usaha_3rdparty" validate:"required" gorm:"column:piutangusaha_3rdparty" example:"10000.00"`
	PiutangusahaBerelasi         *float64 `json:"piutang_usaha_berelasi" validate:"required" example:"10000.00"`
	Piutanglainshortterm3rdparty *float64 `json:"piutang_lain_shortterm_3rdparty" validate:"required" gorm:"column:piutanglainshortterm_3rdparty" example:"10000.00"`
	PiutanglainshorttermBerelasi *float64 `json:"piutang_lain_shortterm_berelasi" validate:"required" example:"10000.00"`
	Piutangberelasishortterm     *float64 `json:"piutang_berelasi_shortterm" validate:"required" example:"10000.00"`
	Piutanglainlongterm3rdparty  *float64 `json:"piutang_lain_longterm_3rdparty" validate:"required" gorm:"column:piutanglainlongterm_3rdparty" example:"10000.00"`
	PiutanglainlongtermBerelasi  *float64 `json:"piutang_lainlongterm_berelasi" validate:"required" example:"10000.00"`
	Piutangberelasilongterm      *float64 `json:"piutang_berelasi_longterm" validate:"required" example:"10000.00"`
	Utangusaha3rdparty           *float64 `json:"utang_usaha_3rdparty" validate:"required" gorm:"column:utangusaha_3rdparty" example:"10000.00"`
	UtangusahaBerelasi           *float64 `json:"utang_usaha_berelasi" validate:"required" example:"10000.00"`
	Utanglainshortterm3rdparty   *float64 `json:"utang_lain_shortterm_3rdparty" validate:"required" gorm:"column:utanglainshortterm_3rdparty" example:"10000.00"`
	UtanglainshorttermBerelasi   *float64 `json:"utang_lain_shortterm_berelasi" validate:"required" example:"10000.00"`
	Utangberelasishortterm       *float64 `json:"utang_berelasi_short_term" validate:"required" example:"10000.00"`
	Utanglainlongterm3rdparty    *float64 `json:"utang_lain_longterm_3rdparty" validate:"required" gorm:"column:utanglainlongterm_3rdparty" example:"10000.00"`
	UtanglainlongtermBerelasi    *float64 `json:"utang_lain_longterm_berelasi" validate:"required" example:"10000.00"`
	Utangberelasilongterm        *float64 `json:"utang_berelasi_longterm" validate:"required" example:"10000.00"`
	SortID                       int      `json:"sort_id" validate:"required" example:"1"`
}

type AgingUtangPiutangDetailFilter struct {
	AgingUtangPiutangID          *int     `query:"aging_utang_piutang_id" validate:"required" example:"1" filter:"CUSTOM"`
	FormatterBridgesID           *int     `query:"formatter_bridges_id" example:"1"`
	Code                         *string  `query:"code" filter:"ILIKE" example:"BELUM_JATUH_TEMPO"`
	Description                  *string  `query:"description" filter:"ILIKE" example:"Belum jatuh tempo"`
	Piutangusaha3rdparty         *float64 `query:"piutang_usaha_3rdparty" example:"10000.00"`
	PiutangusahaBerelasi         *float64 `query:"piutang_usaha_berelasi" example:"10000.00"`
	Piutanglainshortterm3rdparty *float64 `query:"piutang_lainshortterm_3rdparty" example:"10000.00"`
	PiutanglainshorttermBerelasi *float64 `query:"piutang_lainshortterm_berelasi" example:"10000.00"`
	Piutangberelasishortterm     *float64 `query:"piutang_berelasi_shortterm" example:"10000.00"`
	Piutanglainlongterm3rdparty  *float64 `query:"piutang_lain_longterm_3rdparty" example:"10000.00"`
	PiutanglainlongtermBerelasi  *float64 `query:"piutang_lainlongterm_berelasi" example:"10000.00"`
	Piutangberelasilongterm      *float64 `query:"piutang_berelasi_longterm" example:"10000.00"`
	Utangusaha3rdparty           *float64 `query:"utang_usaha_3rdparty" example:"10000.00"`
	UtangusahaBerelasi           *float64 `query:"utang_usaha_berelasi" example:"10000.00"`
	Utanglainshortterm3rdparty   *float64 `query:"utang_lain_shortterm_3rdparty" example:"10000.00"`
	UtanglainshorttermBerelasi   *float64 `query:"utang_lain_shortterm_berelasi" example:"10000.00"`
	Utangberelasishortterm       *float64 `query:"utang_berelasi_short_term" example:"10000.00"`
	Utanglainlongterm3rdparty    *float64 `query:"utang_lain_longterm_3rdparty" example:"10000.00"`
	UtanglainlongtermBerelasi    *float64 `query:"utang_lain_longterm_berelasi" example:"10000.00"`
	Utangberelasilongterm        *float64 `query:"utang_berelasi_longterm" example:"10000.00"`
	SortID                       *int     `query:"sort_id" example:"1"`
}

type AgingUtangPiutangDetailEntityModel struct {
	// abstraction
	// abstraction.Entity
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	AgingUtangPiutangDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
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
	FxSummary      *string `json:"-"`
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
