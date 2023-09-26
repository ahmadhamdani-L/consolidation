package repository

import (
	"fmt"
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type MutasiPersediaanDetail interface {
	Find(ctx *abstraction.Context, m *model.MutasiPersediaanDetailFilterModel) (*[]model.MutasiPersediaanDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiPersediaanDetailEntityModel, error)
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

func (r *mutasipersediaandetail) Find(ctx *abstraction.Context, m *model.MutasiPersediaanDetailFilterModel) (*[]model.MutasiPersediaanDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiPersediaanDetailEntityModel

	query := conn.Model(&model.MutasiPersediaanDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	query = query.Where("formatter_bridges_id IN (?)", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "MUTASI-PERSEDIAAN").Where("trx_ref_id = ?", m.MutasiPersediaanID))

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
