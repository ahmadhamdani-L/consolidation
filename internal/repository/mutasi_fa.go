package repository

import (
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type MutasiFa interface {
	Find(ctx *abstraction.Context, m *model.MutasiFaFilterModel) (*[]model.MutasiFaEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.MutasiFaFilterModel) (*int64, error)
	Update(ctx *abstraction.Context, id *int, e *model.MutasiFaEntityModel) (*model.MutasiFaEntityModel, error)
}

type mutasifa struct {
	abstraction.Repository
}

func NewMutasiFa(db *gorm.DB) *mutasifa {
	return &mutasifa{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *mutasifa) Find(ctx *abstraction.Context, m *model.MutasiFaFilterModel) (*[]model.MutasiFaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiFaEntityModel

	query := conn.Model(&model.MutasiFaEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *mutasifa) GetCount(ctx *abstraction.Context, m *model.MutasiFaFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.MutasiFaEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}

func (r *mutasifa) Update(ctx *abstraction.Context, id *int, e *model.MutasiFaEntityModel) (*model.MutasiFaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}
