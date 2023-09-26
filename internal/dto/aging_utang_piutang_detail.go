package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type AgingUtangPiutangDetailGetRequest struct {
	abstraction.Pagination
	model.AgingUtangPiutangDetailFilterModel
}
type AgingUtangPiutangDetailGetResponse struct {
	Datas model.AgingUtangPiutangEntityModel
}
type AgingUtangPiutangDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                           `json:"meta"`
		Data AgingUtangPiutangDetailGetResponse `json:"data"`
	} `json:"body"`
}

// GetByID
type AgingUtangPiutangDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type AgingUtangPiutangDetailGetByIDResponse struct {
	model.AgingUtangPiutangDetailEntityModel
}
type AgingUtangPiutangDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                               `json:"meta"`
		Data AgingUtangPiutangDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type AgingUtangPiutangDetailCreateRequest struct {
	model.AgingUtangPiutangDetailEntity
}
type AgingUtangPiutangDetailCreateResponse struct {
	model.AgingUtangPiutangDetailEntityModel
}
type AgingUtangPiutangDetailCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                              `json:"meta"`
		Data AgingUtangPiutangDetailCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type AgingUtangPiutangDetailUpdateRequest struct {
	ID                           int      `param:"id" validate:"required,numeric"`
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
}
type AgingUtangPiutangDetailUpdateResponse struct {
	model.AgingUtangPiutangDetailEntityModel
}
type AgingUtangPiutangDetailUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                              `json:"meta"`
		Data AgingUtangPiutangDetailUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type AgingUtangPiutangDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type AgingUtangPiutangDetailDeleteResponse struct {
	// model.AgingUtangPiutangDetailEntityModel
}
type AgingUtangPiutangDetailDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                              `json:"meta"`
		Data AgingUtangPiutangDetailDeleteResponse `json:"data"`
	} `json:"body"`
}
