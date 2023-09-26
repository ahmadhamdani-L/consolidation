package repository

import (
	"worker/internal/abstraction"
	"worker/internal/model"
	"gorm.io/gorm"
)

type FormatterBridges interface {
	Find(ctx *abstraction.Context, m *model.FormatterBridgesFilterModel) (*[]model.FormatterBridgesEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.FormatterBridgesEntityModel, error)
	FindByIDTrx(ctx *abstraction.Context, source *string , trx_ref_id *int) (*model.FormatterBridgesEntityModel, error)
	Create(ctx *abstraction.Context, e *model.FormatterBridgesEntityModel) (*model.FormatterBridgesEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.FormatterBridgesEntityModel) (*model.FormatterBridgesEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.FormatterBridgesEntityModel) (*model.FormatterBridgesEntityModel, error)
	// FindWithDetail(ctx *abstraction.Context, m *model.FormatterBridgesFilterModel) (*model.FormatterBridgesEntityModel, error)
}

type formatterbridges struct {
	abstraction.Repository
}

func NewFormatterBridges(db *gorm.DB) *formatterbridges {
	return &formatterbridges{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *formatterbridges) Find(ctx *abstraction.Context, m *model.FormatterBridgesFilterModel) (*[]model.FormatterBridgesEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.FormatterBridgesEntityModel

	query := conn.Model(&model.FormatterBridgesEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}
// func (r *formatterbridges) FindWithDetail(ctx *abstraction.Context, m *model.FormatterBridgesFilterModel) (*model.FormatterBridgesEntityModel, error) {
// 	conn := r.CheckTrx(ctx)
// 	var data model.FormatterBridgesEntityModel

// 	if err := conn.Where("formatterbridges_for", &m.FormatterBridgesFor).Preload("FormatterBridgesDetail", func(db *gorm.DB) *gorm.DB {
// 		return db.Order("sort_id asc")
// 	}).First(&data).Error; err != nil {
// 		return &data, err
// 	}
// 	return &data, nil
// }

func (r *formatterbridges) FindByID(ctx *abstraction.Context, id *int) (*model.FormatterBridgesEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.FormatterBridgesEntityModel
	if err := conn.Where("id = ?", &id).Preload("FormatterBridges").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *formatterbridges) FindByIDTrx(ctx *abstraction.Context, source *string , trxRefId *int) (*model.FormatterBridgesEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.FormatterBridgesEntityModel
	if err := conn.Where("source = ? AND trx_ref_id = ?", &source, &trxRefId ).Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *formatterbridges) Create(ctx *abstraction.Context, e *model.FormatterBridgesEntityModel) (*model.FormatterBridgesEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *formatterbridges) Update(ctx *abstraction.Context, id *int, e *model.FormatterBridgesEntityModel) (*model.FormatterBridgesEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *formatterbridges) Delete(ctx *abstraction.Context, id *int, e *model.FormatterBridgesEntityModel) (*model.FormatterBridgesEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Where("id =?", id).Delete(e).Error; err != nil {
		return nil, err
	}
	return e, nil
}
