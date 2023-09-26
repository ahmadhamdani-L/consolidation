package repository

import (
	"fmt"
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type MutasiIaDetail interface {
	Find(ctx *abstraction.Context, m *model.MutasiIaDetailFilterModel) (*[]model.MutasiIaDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiIaDetailEntityModel, error)
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

func (r *mutasiiadetail) Find(ctx *abstraction.Context, m *model.MutasiIaDetailFilterModel) (*[]model.MutasiIaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiIaDetailEntityModel

	query := conn.Model(&model.MutasiIaDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	query = query.Where("formatter_bridges_id IN (?)", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "MUTASI-IA").Where("trx_ref_id = ?", m.MutasiIaID))
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
