package repository

import (
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type AgingUtangPiutang interface {
	Find(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel) (*[]model.AgingUtangPiutangEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.AgingUtangPiutangEntityModel) (*model.AgingUtangPiutangEntityModel, error)
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

func (r *agingutangpiutang) Update(ctx *abstraction.Context, id *int, e *model.AgingUtangPiutangEntityModel) (*model.AgingUtangPiutangEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).Error; err != nil {
		return nil, err
	}
	return e, nil
}
