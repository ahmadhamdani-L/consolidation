package repository

import (
	"fmt"
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"

	"gorm.io/gorm"
)

type JcteDetail interface {
	Find(ctx *abstraction.Context, m *model.JcteDetailFilterModel) (*[]model.JcteDetailEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.JcteDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.JcteDetailEntityModel, error)
}

type jctedetail struct {
	abstraction.Repository
}

func NewJcteDetail(db *gorm.DB) *jctedetail {
	return &jctedetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *jctedetail) Find(ctx *abstraction.Context, m *model.JcteDetailFilterModel) (*[]model.JcteDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.JcteDetailEntityModel

	query := conn.Model(&model.JcteDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Preload("Jcte").Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *jctedetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.JcteDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.JcteDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code LIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *jctedetail) FindByID(ctx *abstraction.Context, id *int) (*model.JcteDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JcteDetailEntityModel
	if err := conn.Where("id = ?", &id).Preload("Jcte").First(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
