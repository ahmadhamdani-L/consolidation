package dto

import (
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

type GetVersionRequest struct {
	Period    *string `query:"period"`
	Status    *int    `query:"status"`
	ArrStatus *[]int
	model.CompanyCustomFilter
}

type GetVersionResponse struct {
	Data model.GetVersionModel `json:"data"`
}

type GetVersionResponseDoc struct {
	Body struct {
		Meta res.Meta           `json:"meta"`
		Data GetVersionResponse `json:"data"`
	} `json:"body"`
}
