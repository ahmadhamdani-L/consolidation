package repository

import (
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type ImportedWorksheet interface {
	Create(ctx *abstraction.Context, e *model.ImportedWorksheetEntityModel) (*model.ImportedWorksheetEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.ImportedWorksheetEntityModel) (*model.ImportedWorksheetEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.ImportedWorksheetEntityModel, error)
	CreateNotifikasi(ctx *abstraction.Context, e *model.NotificationEntityModel) (*model.NotificationEntityModel, error)
}

type importedWorksheet struct {
	abstraction.Repository
}

func NewImportedWorksheet(db *gorm.DB) *importedWorksheet {
	return &importedWorksheet{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *importedWorksheet) Create(ctx *abstraction.Context, e *model.ImportedWorksheetEntityModel) (*model.ImportedWorksheetEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *importedWorksheet) Update(ctx *abstraction.Context, id *int, e *model.ImportedWorksheetEntityModel) (*model.ImportedWorksheetEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *importedWorksheet) FindByID(ctx *abstraction.Context, id *int) (*model.ImportedWorksheetEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ImportedWorksheetEntityModel
	err := conn.Where("id = ?", id).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *importedWorksheet) CreateNotifikasi(ctx *abstraction.Context, e *model.NotificationEntityModel) (*model.NotificationEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}
