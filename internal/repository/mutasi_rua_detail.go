package repository

import (
	"fmt"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type MutasiRuaDetail interface {
	Find(ctx *abstraction.Context, m *model.MutasiRuaDetailFilterModel) (*[]model.MutasiRuaDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiRuaDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.MutasiRuaDetailEntityModel) (*model.MutasiRuaDetailEntityModel, error)
}

type mutasiruadetail struct {
	abstraction.Repository
}

func NewMutasiRuaDetail(db *gorm.DB) *mutasiruadetail {
	return &mutasiruadetail{
		abstraction.Repository{
			Db: db,
		},
	}
}


func (r *mutasiruadetail) Create(ctx *abstraction.Context, e *model.MutasiRuaDetailEntityModel) (*model.MutasiRuaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *mutasiruadetail) Find(ctx *abstraction.Context, m *model.MutasiRuaDetailFilterModel) (*[]model.MutasiRuaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiRuaDetailEntityModel

	query := conn.Model(&model.MutasiRuaDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *mutasiruadetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiRuaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.MutasiRuaDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code ILIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
