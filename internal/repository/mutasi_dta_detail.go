package repository

import (
	"fmt"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type MutasiDtaDetail interface {
	Find(ctx *abstraction.Context, m *model.MutasiDtaDetailFilterModel) (*[]model.MutasiDtaDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiDtaDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.MutasiDtaDetailEntityModel) (*model.MutasiDtaDetailEntityModel, error)
}

type mutasidtadetail struct {
	abstraction.Repository
}

func NewMutasiDtaDetail(db *gorm.DB) *mutasidtadetail {
	return &mutasidtadetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *mutasidtadetail) Create(ctx *abstraction.Context, e *model.MutasiDtaDetailEntityModel) (*model.MutasiDtaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *mutasidtadetail) Find(ctx *abstraction.Context, m *model.MutasiDtaDetailFilterModel) (*[]model.MutasiDtaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiDtaDetailEntityModel

	query := conn.Model(&model.MutasiDtaDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *mutasidtadetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiDtaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.MutasiDtaDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code ILIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
