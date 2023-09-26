package repository

import (
	"fmt"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type MutasiFaDetail interface {
	Find(ctx *abstraction.Context, m *model.MutasiFaDetailFilterModel) (*[]model.MutasiFaDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiFaDetailEntityModel, error)
}

type mutasifadetail struct {
	abstraction.Repository
}

func NewMutasiFaDetail(db *gorm.DB) *mutasifadetail {
	return &mutasifadetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *mutasifadetail) Find(ctx *abstraction.Context, m *model.MutasiFaDetailFilterModel) (*[]model.MutasiFaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiFaDetailEntityModel

	query := conn.Model(&model.MutasiFaDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	query = query.Where("formatter_bridges_id IN (?)", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "MUTASI-FA").Where("trx_ref_id = ?", m.MutasiFaID))
	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *mutasifadetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiFaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.MutasiFaDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code ILIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
