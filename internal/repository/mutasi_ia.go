package repository

import (
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type MutasiIa interface {
	Find(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*[]model.MutasiIaEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*int64, error)
	Create(ctx *abstraction.Context, e *model.MutasiIaEntityModel) (*model.MutasiIaEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.MutasiIaEntityModel) (*model.MutasiIaEntityModel, error)
	FindByID(ctx *abstraction.Context, version *int, company *int, period *string) (*model.MutasiIaEntityModel, error)
}

type mutasiia struct {
	abstraction.Repository
}

func NewMutasiIa(db *gorm.DB) *mutasiia {
	return &mutasiia{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *mutasiia) FindByID(ctx *abstraction.Context, version *int, company *int, period *string) (*model.MutasiIaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.MutasiIaEntityModel
	err := conn.Where("versions = ? AND company_id = ? AND period = ?", version, company, period).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *mutasiia) Update(ctx *abstraction.Context, id *int, e *model.MutasiIaEntityModel) (*model.MutasiIaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).Preload("Company").Preload("UserCreated").Preload("UserModified").First(e).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil

}

func (r *mutasiia) Create(ctx *abstraction.Context, e *model.MutasiIaEntityModel) (*model.MutasiIaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Preload("UserCreated").Preload("UserModified").First(e).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name

	return e, nil
}

func (r *mutasiia) Find(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*[]model.MutasiIaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiIaEntityModel

	query := conn.Model(&model.MutasiIaEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *mutasiia) GetCount(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.MutasiIaEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}
