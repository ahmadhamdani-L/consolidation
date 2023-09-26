package dto

import (
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

type ExportRequest struct {
	ImportID int    `query:"import_id" validate:"required,numeric"`
	Request  string `query:"file_request"`
}
type ExportResponse struct {
	Message string `json:"message"`
}
type ExportResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data []model.FormatterEntityModel `json:"data"`
	} `json:"body"`
}

type GetExportRequest struct {
	NotificationID int `param:"notification_id" validate:"required,numeric"`
}

type ExportConsolRequest struct {
	ConsolidationID int `query:"consolidation_id" validate:"required"`
}
