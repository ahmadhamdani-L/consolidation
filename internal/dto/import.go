package dto

import (
	"mcash-finance-console-core/internal/model"
)

type ImportedWorksheetRequest struct {
	model.ImportedWorksheetDetailEntity
}

type ImportReUploadRequest struct {
	ID   int `param:"id" validate:"required,numeric"`
}
type ImportReUploadResponse struct {
	model.ImportedWorksheetEntityModel
}

type ImportReUploadDetailRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type ImportReUploadDetailResponse struct {
	Data []model.ImportedWorksheetDetailEntityModel
}

type ImportReUploadTemplateRequest struct {
	Template string `query:"template" `
}
type AjeTemplateRequest struct {
	Jurnal string `param:"jurnal"`
}