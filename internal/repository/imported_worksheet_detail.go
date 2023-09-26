package repository

import (
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type ImportedWorksheetDetail interface {
	Create(ctx *abstraction.Context, e *model.ImportedWorksheetDetailEntityModel) (*model.ImportedWorksheetDetailEntityModel, error)
	GetCountStatus(ctx *abstraction.Context, e *model.ImportedWorksheetDetailEntityModel) (*[]model.ImportedWorksheetDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.ImportedWorksheetDetailEntityModel) (*model.ImportedWorksheetDetailEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.ImportedWorksheetDetailEntityModel, error)
}

type importedWorksheetdetail struct {
	abstraction.Repository
}

func NewImportedWorksheetDetail(db *gorm.DB) *importedWorksheetdetail {
	return &importedWorksheetdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *importedWorksheetdetail) Create(ctx *abstraction.Context, e *model.ImportedWorksheetDetailEntityModel) (*model.ImportedWorksheetDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *importedWorksheetdetail) GetCountStatus (ctx *abstraction.Context, e *model.ImportedWorksheetDetailEntityModel) (*[]model.ImportedWorksheetDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.ImportedWorksheetDetailEntityModel

	if err := conn.Where("imported_worksheet_id =? AND status =? ", e.ImportedWorksheetID, e.Status).Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *importedWorksheetdetail) Update(ctx *abstraction.Context, id *int, e *model.ImportedWorksheetDetailEntityModel) (*model.ImportedWorksheetDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *importedWorksheetdetail) FindByID(ctx *abstraction.Context, id *int) (*model.ImportedWorksheetDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ImportedWorksheetDetailEntityModel
	err := conn.Where("id = ?", id).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}
