package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type MutasiPersediaanDetail interface {
	Find(ctx *abstraction.Context, m *model.MutasiPersediaanDetailFilterModel, p *abstraction.Pagination) (*[]model.MutasiPersediaanDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.MutasiPersediaanDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.MutasiPersediaanDetailEntityModel) (*model.MutasiPersediaanDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.MutasiPersediaanDetailEntityModel) (*model.MutasiPersediaanDetailEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.MutasiPersediaanDetailEntityModel) (*model.MutasiPersediaanDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiPersediaanDetailEntityModel, error)
	FindWithFormatter(ctx *abstraction.Context, m *model.MutasiPersediaanDetailFilterModel) (*[]model.MutasiPersediaanDetailFmtEntityModel, error)
	FindSummary(ctx *abstraction.Context, formatterID, formatterBridgesID *int, sort_id *float64) (*model.MutasiPersediaanDetailEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.MutasiPersediaanDetailFilterModel) (data *model.MutasiPersediaanDetailEntityModel, err error)
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
func (r *mutasipersediaandetail) FindByCriteria(ctx *abstraction.Context, filter *model.MutasiPersediaanDetailFilterModel) (data *model.MutasiPersediaanDetailEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.MutasiPersediaanDetailEntityModel{})
	query = r.Filter(ctx, query, *filter)
	if err = query.First(&data).Error; err != nil {
		return
	}
	// if data.ID == 0 {
	// 	err = errors.New("Data Not Found")
	// }
	return
}
func (r *mutasipersediaandetail) Find(ctx *abstraction.Context, m *model.MutasiPersediaanDetailFilterModel, p *abstraction.Pagination) (*[]model.MutasiPersediaanDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiPersediaanDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.MutasiPersediaanDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	//sort
	if p.Sort == nil {
		sort := "asc"
		p.Sort = &sort
	}
	if p.SortBy == nil {
		sortBy := "id"
		p.SortBy = &sortBy
	}

	sort := fmt.Sprintf("%s %s", *p.SortBy, *p.Sort)
	query = query.Order(sort)

	//pagination
	if p.Page == nil {
		page := 1
		p.Page = &page
	}
	if p.PageSize == nil {
		pageSize := 10
		p.PageSize = &pageSize
	}
	info = abstraction.PaginationInfo{
		Pagination: p,
	}
	limit := *p.PageSize
	offset := limit * (*p.Page - 1)
	query = query.Where("formatter_bridges_id IN (?)", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "MUTASI-PERSEDIAAN").Where("trx_ref_id = ?", m.MutasiPersediaanID))
	var totalData int64
	query = query.Count(&totalData).Limit(limit).Offset(offset)
	if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
	}

	info.Count = int(totalData)
	info.MoreRecords = false
	if len(datas) > *p.PageSize {
		info.MoreRecords = true
		// info.Count -= 1
		// datas = datas[:len(datas)-1]
	}

	return &datas, &info, nil
}

func (r *mutasipersediaandetail) FindByID(ctx *abstraction.Context, id *int) (*model.MutasiPersediaanDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.MutasiPersediaanDetailEntityModel
	if err := conn.Where("id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *mutasipersediaandetail) Create(ctx *abstraction.Context, e *model.MutasiPersediaanDetailEntityModel) (*model.MutasiPersediaanDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *mutasipersediaandetail) Update(ctx *abstraction.Context, id *int, e *model.MutasiPersediaanDetailEntityModel) (*model.MutasiPersediaanDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *mutasipersediaandetail) Delete(ctx *abstraction.Context, id *int, e *model.MutasiPersediaanDetailEntityModel) (*model.MutasiPersediaanDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *mutasipersediaandetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiPersediaanDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.MutasiPersediaanDetailEntityModel
	if err := conn.Where("code ILIKE ?", *code+"%").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *mutasipersediaandetail) FindWithFormatter(ctx *abstraction.Context, m *model.MutasiPersediaanDetailFilterModel) (*[]model.MutasiPersediaanDetailFmtEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data []model.MutasiPersediaanDetailFmtEntityModel
	// err := conn.Model(&model.FormatterBridgesEntityModel{}).
	// 	Joins("INNER JOIN m_formatter f ON formatter_bridges.formatter_id = f.id").
	// 	Joins("INNER JOIN m_formatter_detail fd ON f.id = fd.formatter_id").
	// 	Joins("INNER JOIN mutasi_persediaan_detail tb ON fd.code = tb.code AND tb.formatter_bridges_id = formatter_bridges.id").
	// 	Where("formatter_bridges_id = ?", m.FormatterBridgesID).Order("fd.sort_id ASC").
	// 	Select("tb.*, auto_summary, is_total, is_control, is_label, control_formula").Find(&data).Error
	// if err != nil {
	// 	return nil, err
	// }

	query := conn.Model(&model.FormatterBridgesEntityModel{}).
		Joins("INNER JOIN m_formatter f ON formatter_bridges.formatter_id = f.id").
		Joins("INNER JOIN m_formatter_detail fd ON f.id = fd.formatter_id").
		Joins("INNER JOIN mutasi_persediaan_detail tb ON fd.code = tb.code AND tb.formatter_bridges_id = formatter_bridges.id").
		Select("tb.*, auto_summary, is_total, is_control, is_label, control_formula").Order("fd.sort_id ASC")

	if m.FormatterBridgesID != nil && *m.FormatterBridgesID != 0 {
		query = query.Where("formatter_bridges_id = ?", m.FormatterBridgesID)
	}

	if m.MutasiPersediaanID != nil && *m.MutasiPersediaanID != 0 {
		query = query.Where("trx_ref_id = ? AND source = ?", m.MutasiPersediaanID, "MUTASI-PERSEDIAAN")
	}

	err := query.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *mutasipersediaandetail) FindSummary(ctx *abstraction.Context, formatterID, formatterBridgesID *int, sort_id *float64) (*model.MutasiPersediaanDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.MutasiPersediaanDetailEntityModel
	var tmp model.FormatterDetailEntityModel
	err := conn.Model(&model.FormatterDetailEntityModel{}).Where("formatter_id = ? AND auto_summary IS NOT NULL AND auto_summary = true AND sort_id < ?", formatterID, sort_id).Order("sort_id DESC").Limit(1).Find(&tmp).Error
	if err != nil {
		return nil, err
	}

	lastID := 0.0
	if tmp.SortID != 0 {
		lastID = tmp.SortID
	}

	err = conn.Raw("SELECT SUM( amount ) amount FROM (SELECT tb.* FROM formatter_bridges INNER JOIN m_formatter f ON formatter_bridges.formatter_id = f.ID INNER JOIN m_formatter_detail fd ON f.ID = fd.formatter_id INNER JOIN mutasi_persediaan_detail tb ON fd.code = tb.code AND tb.formatter_bridges_id = formatter_bridges.ID WHERE ( formatter_bridges_id = ? ) AND fd.sort_id < ? AND fd.sort_id > ? ORDER BY fd.sort_id ASC) tmp ", formatterBridgesID, sort_id, lastID).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}
