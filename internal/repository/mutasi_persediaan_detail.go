package repository

import (
	"fmt"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type MutasiPersediaanDetail interface {
	Find(ctx *abstraction.Context, m *model.MutasiPersediaanDetailFilterModel) (*[]model.MutasiPersediaanDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiPersediaanDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.MutasiPersediaanDetailEntityModel) (*model.MutasiPersediaanDetailEntityModel, error)
}

type mutasipersediaandetail struct {
	abstraction.Repository
}

func NewMutasiPersediaanDetail(db *gorm.DB) *mutasipersediaandetail {
	return &mutasipersediaandetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *mutasipersediaandetail) Create(ctx *abstraction.Context, e *model.MutasiPersediaanDetailEntityModel) (*model.MutasiPersediaanDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *mutasipersediaandetail) Find(ctx *abstraction.Context, m *model.MutasiPersediaanDetailFilterModel) (*[]model.MutasiPersediaanDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiPersediaanDetailEntityModel

	query := conn.Model(&model.MutasiPersediaanDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *mutasipersediaandetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiPersediaanDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.MutasiPersediaanDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code ILIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
