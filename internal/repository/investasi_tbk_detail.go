package repository

import (
	"fmt"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type InvestasiTbkDetail interface {
	Find(ctx *abstraction.Context, m *model.InvestasiTbkDetailFilterModel) (*[]model.InvestasiTbkDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.InvestasiTbkDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.InvestasiTbkDetailEntityModel) (*model.InvestasiTbkDetailEntityModel, error)
}

type investasitbkdetail struct {
	abstraction.Repository
}

func NewInvestasiTbkDetail(db *gorm.DB) *investasitbkdetail {
	return &investasitbkdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *investasitbkdetail) Find(ctx *abstraction.Context, m *model.InvestasiTbkDetailFilterModel) (*[]model.InvestasiTbkDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.InvestasiTbkDetailEntityModel

	query := conn.Model(&model.InvestasiTbkDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *investasitbkdetail) Create(ctx *abstraction.Context, e *model.InvestasiTbkDetailEntityModel) (*model.InvestasiTbkDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *investasitbkdetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.InvestasiTbkDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.InvestasiTbkDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code ILIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
