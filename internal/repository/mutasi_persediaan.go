package repository

import (
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type MutasiPersediaan interface {
	Find(ctx *abstraction.Context, m *model.MutasiPersediaanFilterModel) (*[]model.MutasiPersediaanEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.MutasiPersediaanFilterModel) (*int64, error)
	Create(ctx *abstraction.Context, e *model.MutasiPersediaanEntityModel) (*model.MutasiPersediaanEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.MutasiPersediaanEntityModel) (*model.MutasiPersediaanEntityModel, error)
	FindByID(ctx *abstraction.Context, version *int, company *int, period *string) (*model.MutasiPersediaanEntityModel, error)
}

type mutasipersediaan struct {
	abstraction.Repository
}

func NewMutasiPersediaan(db *gorm.DB) *mutasipersediaan {
	return &mutasipersediaan{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *mutasipersediaan) FindByID(ctx *abstraction.Context, version *int, company *int, period *string) (*model.MutasiPersediaanEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.MutasiPersediaanEntityModel
	err := conn.Where("versions = ? AND company_id = ? AND period = ?", version, company, period).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *mutasipersediaan) Update(ctx *abstraction.Context, id *int, e *model.MutasiPersediaanEntityModel) (*model.MutasiPersediaanEntityModel, error) {
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

func (r *mutasipersediaan) Create(ctx *abstraction.Context, e *model.MutasiPersediaanEntityModel) (*model.MutasiPersediaanEntityModel, error) {
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

func (r *mutasipersediaan) Find(ctx *abstraction.Context, m *model.MutasiPersediaanFilterModel) (*[]model.MutasiPersediaanEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiPersediaanEntityModel

	query := conn.Model(&model.MutasiPersediaanEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *mutasipersediaan) GetCount(ctx *abstraction.Context, m *model.MutasiPersediaanFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.MutasiPersediaanEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}
