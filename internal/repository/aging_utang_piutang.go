package repository

import (
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type AgingUtangPiutang interface {
	Find(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel) (*[]model.AgingUtangPiutangEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel) (*int64, error)
	Create(ctx *abstraction.Context, e *model.AgingUtangPiutangEntityModel) (*model.AgingUtangPiutangEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.AgingUtangPiutangEntityModel) (*model.AgingUtangPiutangEntityModel, error)
	FindByID(ctx *abstraction.Context, version *int, company *int, period *string) (*model.AgingUtangPiutangEntityModel, error)
}

type agingutangpiutang struct {
	abstraction.Repository
}

func NewAgingUtangPiutang(db *gorm.DB) *agingutangpiutang {
	return &agingutangpiutang{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *agingutangpiutang) FindByID(ctx *abstraction.Context, version *int, company *int, period *string) (*model.AgingUtangPiutangEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.AgingUtangPiutangEntityModel
	err := conn.Where("versions = ? AND company_id = ? AND period = ?", version, company, period).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *agingutangpiutang) Update(ctx *abstraction.Context, id *int, e *model.AgingUtangPiutangEntityModel) (*model.AgingUtangPiutangEntityModel, error) {
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

func (r *agingutangpiutang) GetCount(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.AgingUtangPiutangEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}

func (r *agingutangpiutang) Create(ctx *abstraction.Context, e *model.AgingUtangPiutangEntityModel) (*model.AgingUtangPiutangEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return e, err
	}
	if err := conn.Model(e).Preload("UserCreated").Preload("UserModified").First(e).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name

	return e, nil
}

func (r *agingutangpiutang) Find(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel) (*[]model.AgingUtangPiutangEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.AgingUtangPiutangEntityModel

	query := conn.Model(&model.AgingUtangPiutangEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}
	return &datas, nil
}
