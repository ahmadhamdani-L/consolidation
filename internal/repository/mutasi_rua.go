package repository

import (
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type MutasiRua interface {
	Find(ctx *abstraction.Context, m *model.MutasiRuaFilterModel) (*[]model.MutasiRuaEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.MutasiRuaFilterModel) (*int64, error)
	Create(ctx *abstraction.Context, e *model.MutasiRuaEntityModel) (*model.MutasiRuaEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.MutasiRuaEntityModel) (*model.MutasiRuaEntityModel, error)
	FindByID(ctx *abstraction.Context, version *int, company *int, period *string) (*model.MutasiRuaEntityModel, error)
}

type mutasirua struct {
	abstraction.Repository
}

func NewMutasiRua(db *gorm.DB) *mutasirua {
	return &mutasirua{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *mutasirua) FindByID(ctx *abstraction.Context, version *int, company *int, period *string) (*model.MutasiRuaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.MutasiRuaEntityModel
	err := conn.Where("versions = ? AND company_id = ? AND period = ?", version, company, period).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *mutasirua) Update(ctx *abstraction.Context, id *int, e *model.MutasiRuaEntityModel) (*model.MutasiRuaEntityModel, error) {
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

func (r *mutasirua) Create(ctx *abstraction.Context, e *model.MutasiRuaEntityModel) (*model.MutasiRuaEntityModel, error) {
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

func (r *mutasirua) Find(ctx *abstraction.Context, m *model.MutasiRuaFilterModel) (*[]model.MutasiRuaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiRuaEntityModel

	query := conn.Model(&model.MutasiRuaEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *mutasirua) GetCount(ctx *abstraction.Context, m *model.MutasiRuaFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.MutasiRuaEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}
