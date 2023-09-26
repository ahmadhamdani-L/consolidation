package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type InvestasiTbkDetailGetRequest struct {
	abstraction.Pagination
	model.InvestasiTbkDetailFilterModel
}
type InvestasiTbkDetailGetResponse struct {
	Datas struct {
		TotalAmountCost     float64                       `json:"total_amount_cost"`
		TotalAmountFv       float64                       `json:"total_amount_fv"`
		TotalUnrealizedGain float64                       `json:"total_unrealized_gain"`
		TotalRealizedGain   float64                       `json:"total_realized_gain"`
		TotalFee            float64                       `json:"total_fee"`
		Data                model.InvestasiTbkEntityModel `json:"investasi_tbk"`
	}
	PaginationInfo abstraction.PaginationInfo
}
type InvestasiTbkDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data InvestasiTbkDetailGetResponse `json:"data"`
	} `json:"body"`
}

// GetByID
type InvestasiTbkDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type InvestasiTbkDetailGetByIDResponse struct {
	model.InvestasiTbkDetailEntityModel
}
type InvestasiTbkDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                          `json:"meta"`
		Data InvestasiTbkDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type InvestasiTbkDetailCreateRequest struct {
	model.InvestasiTbkDetailEntity
}
type InvestasiTbkDetailCreateResponse struct {
	model.InvestasiTbkDetailEntityModel
}
type InvestasiTbkDetailCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                         `json:"meta"`
		Data InvestasiTbkDetailCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type InvestasiTbkDetailUpdateRequest struct {
	ID           int      `param:"id" validate:"required,numeric"`
	EndingShares *float64 `json:"ending_shares" validate:"required" example:"10000.00"`
	AvgPrice     *float64 `json:"avg_price" validate:"required" example:"1000.00"`
	ClosingPrice *float64 `json:"closing_price" validate:"required" example:"10000.00"`
	RealizedGain *float64 `json:"realized_gain" validate:"required" example:"10000.00"`
	Fee          *float64 `json:"fee" validate:"required" example:"10000.00"`
}
type InvestasiTbkDetailUpdateResponse struct {
	model.InvestasiTbkDetailEntityModel
}
type InvestasiTbkDetailUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                         `json:"meta"`
		Data InvestasiTbkDetailUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type InvestasiTbkDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type InvestasiTbkDetailDeleteResponse struct {
	model.InvestasiTbkDetailEntityModel
}
type InvestasiTbkDetailDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                         `json:"meta"`
		Data InvestasiTbkDetailDeleteResponse `json:"data"`
	} `json:"body"`
}
