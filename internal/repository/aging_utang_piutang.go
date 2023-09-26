package repository

import (
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type AgingUtangPiutang interface {
	Find(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel) (*[]model.AgingUtangPiutangEntityModel, error)
}

type agingutangpiutang struct {
	abstraction.Repository
}

func NewAgingUtangPiutang(db *gorm.DB) *agingutangpiutang {
	return &agingutangpiutang{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *agingutangpiutang) Find(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel) (*[]model.AgingUtangPiutangEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.AgingUtangPiutangEntityModel

	query := conn.Model(&model.AgingUtangPiutangEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}
