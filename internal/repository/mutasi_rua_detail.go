package repository

import (
	"fmt"
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type MutasiRuaDetail interface {
	Find(ctx *abstraction.Context, m *model.MutasiRuaDetailFilterModel) (*[]model.MutasiRuaDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiRuaDetailEntityModel, error)
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

func (r *mutasiruadetail) Find(ctx *abstraction.Context, m *model.MutasiRuaDetailFilterModel) (*[]model.MutasiRuaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiRuaDetailEntityModel

	query := conn.Model(&model.MutasiRuaDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	query = query.Where("formatter_bridges_id IN (?)", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "MUTASI-RUA").Where("trx_ref_id = ?", m.MutasiRuaID))
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
