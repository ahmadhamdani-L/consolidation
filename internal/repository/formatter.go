package repository

import (
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"

	"gorm.io/gorm"
)

type Formatter interface {
	Find(ctx *abstraction.Context, m *model.FormatterFilterModel) (*[]model.FormatterEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.FormatterEntityModel, error)
	Create(ctx *abstraction.Context, e *model.FormatterEntityModel) (*model.FormatterEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.FormatterEntityModel) (*model.FormatterEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.FormatterEntityModel) (*model.FormatterEntityModel, error)
	FindWithDetail(ctx *abstraction.Context, m *model.FormatterFilterModel) (*model.FormatterEntityModel, error)
}

type formatter struct {
	abstraction.Repository
}

func NewFormatter(db *gorm.DB) *formatter {
	return &formatter{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *formatter) Find(ctx *abstraction.Context, m *model.FormatterFilterModel) (*[]model.FormatterEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.FormatterEntityModel

	query := conn.Model(&model.FormatterEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *formatter) FindWithDetail(ctx *abstraction.Context, m *model.FormatterFilterModel) (*model.FormatterEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.FormatterEntityModel

	if err := conn.Where("formatter_for", &m.FormatterFor).Preload("FormatterDetail", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_id asc")
	}).First(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *formatter) FindByID(ctx *abstraction.Context, id *int) (*model.FormatterEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.FormatterEntityModel
	if err := conn.Where("id = ?", &id).Preload("FormatterDetail").First(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *formatter) Create(ctx *abstraction.Context, e *model.FormatterEntityModel) (*model.FormatterEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *formatter) Update(ctx *abstraction.Context, id *int, e *model.FormatterEntityModel) (*model.FormatterEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", id).Updates(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", id).First(e).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *formatter) Delete(ctx *abstraction.Context, id *int, e *model.FormatterEntityModel) (*model.FormatterEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Where("id =?", id).Delete(e).Error; err != nil {
		return nil, err
	}
	return e, nil
}
