package repository

import (
	"fmt"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type Coa interface {
	Find(ctx *abstraction.Context, m *model.CoaFilterModel) (*[]model.CoaEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.CoaEntityModel, error)
}

type coa struct {
	abstraction.Repository
}

func NewCoa(db *gorm.DB) *coa {
	return &coa{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *coa) Find(ctx *abstraction.Context, m *model.CoaFilterModel) (*[]model.CoaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.CoaEntityModel

	query := conn.Model(&model.CoaEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *coa) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.CoaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.CoaEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code LIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
