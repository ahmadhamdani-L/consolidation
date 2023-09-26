package repository

import (
	"fmt"
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"

	"gorm.io/gorm"
)

type InvestasiNonTbkDetail interface {
	Find(ctx *abstraction.Context, m *model.InvestasiNonTbkDetailFilterModel) (*[]model.InvestasiNonTbkDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.InvestasiNonTbkDetailEntityModel, error)
}

type investasinontbkdetail struct {
	abstraction.Repository
}

func NewInvestasiNonTbkDetail(db *gorm.DB) *investasinontbkdetail {
	return &investasinontbkdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *investasinontbkdetail) Find(ctx *abstraction.Context, m *model.InvestasiNonTbkDetailFilterModel) (*[]model.InvestasiNonTbkDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.InvestasiNonTbkDetailEntityModel

	query := conn.Model(&model.InvestasiNonTbkDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *investasinontbkdetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.InvestasiNonTbkDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.InvestasiNonTbkDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code ILIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
