package repository

import (
	"fmt"
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"

	"gorm.io/gorm"
)

type JpmDetail interface {
	Find(ctx *abstraction.Context, m *model.JpmDetailFilterModel) (*[]model.JpmDetailEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.JpmDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.JpmDetailEntityModel, error)
}

type jpmdetail struct {
	abstraction.Repository
}

func NewJpmDetail(db *gorm.DB) *jpmdetail {
	return &jpmdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *jpmdetail) Find(ctx *abstraction.Context, m *model.JpmDetailFilterModel) (*[]model.JpmDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.JpmDetailEntityModel

	query := conn.Model(&model.JpmDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *jpmdetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.JpmDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.JpmDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code LIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *jpmdetail) FindByID(ctx *abstraction.Context, id *int) (*model.JpmDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JpmDetailEntityModel
	if err := conn.Where("id = ?", &id).First(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
