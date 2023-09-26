package repository

import (
	"fmt"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type MutasiIaDetail interface {
	Find(ctx *abstraction.Context, m *model.MutasiIaDetailFilterModel) (*[]model.MutasiIaDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiIaDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.MutasiIaDetailEntityModel) (*model.MutasiIaDetailEntityModel, error)
}

type mutasiiadetail struct {
	abstraction.Repository
}

func NewMutasiIaDetail(db *gorm.DB) *mutasiiadetail {
	return &mutasiiadetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *mutasiiadetail) Create(ctx *abstraction.Context, e *model.MutasiIaDetailEntityModel) (*model.MutasiIaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *mutasiiadetail) Find(ctx *abstraction.Context, m *model.MutasiIaDetailFilterModel) (*[]model.MutasiIaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiIaDetailEntityModel

	query := conn.Model(&model.MutasiIaDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *mutasiiadetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiIaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.MutasiIaDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code ILIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
