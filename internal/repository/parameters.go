package repository

import (
	"fmt"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type Parameter interface {
	Find(ctx *abstraction.Context, m *model.ParameterFilterModel) (*[]model.ParameterEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.ParameterEntityModel, error)
}

type parameter struct {
	abstraction.Repository
}

func NewParameter(db *gorm.DB) *parameter {
	return &parameter{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *parameter) Find(ctx *abstraction.Context, m *model.ParameterFilterModel) (*[]model.ParameterEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ParameterEntityModel

	query := conn.Model(&model.ParameterEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}


func (r *parameter) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.ParameterEntityModel, error) {
	conn := r.CheckTrx(ctx)
  
	var data []model.ParameterEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code LIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}