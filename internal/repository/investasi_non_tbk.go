package repository

import (
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type InvestasiNonTbk interface {
	Find(ctx *abstraction.Context, m *model.InvestasiNonTbkFilterModel) (*[]model.InvestasiNonTbkEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.InvestasiNonTbkFilterModel) (*int64, error)
	Create(ctx *abstraction.Context, e *model.InvestasiNonTbkEntityModel) (*model.InvestasiNonTbkEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.InvestasiNonTbkEntityModel) (*model.InvestasiNonTbkEntityModel, error)
	FindByID(ctx *abstraction.Context, version *int, company *int, period *string) (*model.InvestasiNonTbkEntityModel, error)
}

type investasinontbk struct {
	abstraction.Repository
}

func NewInvestasiNonTbk(db *gorm.DB) *investasinontbk {
	return &investasinontbk{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *investasinontbk) FindByID(ctx *abstraction.Context, version *int, company *int, period *string) (*model.InvestasiNonTbkEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.InvestasiNonTbkEntityModel
	err := conn.Where("versions = ? AND company_id = ? AND period = ?", version, company, period).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *investasinontbk) Update(ctx *abstraction.Context, id *int, e *model.InvestasiNonTbkEntityModel) (*model.InvestasiNonTbkEntityModel, error) {
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

func (r *investasinontbk) Create(ctx *abstraction.Context, e *model.InvestasiNonTbkEntityModel) (*model.InvestasiNonTbkEntityModel, error) {
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

func (r *investasinontbk) Find(ctx *abstraction.Context, m *model.InvestasiNonTbkFilterModel) (*[]model.InvestasiNonTbkEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.InvestasiNonTbkEntityModel

	query := conn.Model(&model.InvestasiNonTbkEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *investasinontbk) GetCount(ctx *abstraction.Context, m *model.InvestasiNonTbkFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.InvestasiNonTbkEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}
