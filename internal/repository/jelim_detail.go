package repository

import (
	"fmt"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type JelimDetail interface {
	Find(ctx *abstraction.Context, m *model.JelimDetailFilterModel) (*[]model.JelimDetailEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.JelimDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.JelimDetailEntityModel, error)
}

type jelimdetail struct {
	abstraction.Repository
}

func NewJelimDetail(db *gorm.DB) *jelimdetail {
	return &jelimdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *jelimdetail) Find(ctx *abstraction.Context, m *model.JelimDetailFilterModel) (*[]model.JelimDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.JelimDetailEntityModel

	query := conn.Model(&model.JelimDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	if err := query.Preload("Jelim").Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *jelimdetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.JelimDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.JelimDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code LIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *jelimdetail) FindByID(ctx *abstraction.Context, id *int) (*model.JelimDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JelimDetailEntityModel
	if err := conn.Where("id = ?", &id).Preload("Jelim").First(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
