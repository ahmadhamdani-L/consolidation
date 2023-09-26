package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type InvestasiTbkDetail interface {
	Find(ctx *abstraction.Context, m *model.InvestasiTbkDetailFilterModel, p *abstraction.Pagination) (*[]model.InvestasiTbkDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.InvestasiTbkDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.InvestasiTbkDetailEntityModel) (*model.InvestasiTbkDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.InvestasiTbkDetailEntityModel) (*model.InvestasiTbkDetailEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.InvestasiTbkDetailEntityModel) (*model.InvestasiTbkDetailEntityModel, error)
	GetTotal(ctx *abstraction.Context, m *model.InvestasiTbkDetailFilterModel) (float64, float64, float64, float64, float64, error)
	FindWithFormatter(ctx *abstraction.Context, m *model.InvestasiTbkDetailFilterModel, p *abstraction.Pagination) (*[]model.InvestasiTbkDetailFmtEntityModel, *abstraction.PaginationInfo, error)
	FindByStock(ctx *abstraction.Context, fmtBridgesID *int, stock *string) (*model.InvestasiTbkDetailEntityModel, error)
}

type investasitbkdetail struct {
	abstraction.Repository
}

func NewInvestasiTbkDetail(db *gorm.DB) *investasitbkdetail {
	return &investasitbkdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *investasitbkdetail) Find(ctx *abstraction.Context, m *model.InvestasiTbkDetailFilterModel, p *abstraction.Pagination) (*[]model.InvestasiTbkDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.InvestasiTbkDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.InvestasiTbkDetailEntityModel{})
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
	query = query.Where("formatter_bridges_id IN (?)", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "INVESTASI-TBK").Where("trx_ref_id = ?", m.InvestasiTbkID))
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

func (r *investasitbkdetail) FindByID(ctx *abstraction.Context, id *int) (*model.InvestasiTbkDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.InvestasiTbkDetailEntityModel
	if err := conn.Where("id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *investasitbkdetail) Create(ctx *abstraction.Context, e *model.InvestasiTbkDetailEntityModel) (*model.InvestasiTbkDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *investasitbkdetail) Update(ctx *abstraction.Context, id *int, e *model.InvestasiTbkDetailEntityModel) (*model.InvestasiTbkDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *investasitbkdetail) Delete(ctx *abstraction.Context, id *int, e *model.InvestasiTbkDetailEntityModel) (*model.InvestasiTbkDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *investasitbkdetail) FindWithFormatter(ctx *abstraction.Context, m *model.InvestasiTbkDetailFilterModel, p *abstraction.Pagination) (*[]model.InvestasiTbkDetailFmtEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)
	var data []model.InvestasiTbkDetailFmtEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.FormatterBridgesEntityModel{})
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
	query = query.Joins("INNER JOIN m_formatter f ON formatter_bridges.formatter_id = f.id").Joins("INNER JOIN m_formatter_detail fd ON f.id = fd.formatter_id").
		Joins("INNER JOIN investasi_tbk_detail tb ON fd.code = tb.stock AND tb.formatter_bridges_id = formatter_bridges.id").Where("trx_ref_id = ? AND source = ?", m.InvestasiTbkID, "INVESTASI-TBK").
		Order("fd.sort_id ASC").Select("tb.*, auto_summary, is_total, is_control, is_label, control_formula")

	querycount := query
	var totalData int64
	if err := querycount.Count(&totalData).Error; err != nil {
		return &data, &info, err
	}
	query = query.Count(&totalData).Limit(limit).Offset(offset)
	if err := query.Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, &info, err
	}

	info.Count = int(totalData)
	info.MoreRecords = false
	if len(data) > *p.PageSize {
		info.MoreRecords = true
		info.Count -= 1
		data = data[:len(data)-1]
	}

	return &data, &info, nil
}

func (r *investasitbkdetail) GetTotal(ctx *abstraction.Context, m *model.InvestasiTbkDetailFilterModel) (float64, float64, float64, float64, float64, error) {
	conn := r.CheckTrx(ctx)

	type Total struct {
		TotalAmountCost     float64
		TotalAmountFv       float64
		TotalUnrealizedGain float64
		TotalRealizedGain   float64
		TotalFee            float64
	}
	var data Total

	query := conn.Model(&model.InvestasiTbkDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	query = query.Where("formatter_bridges_id IN (?)", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "INVESTASI-TBK").Where("trx_ref_id = ?", m.InvestasiTbkID))
	if err := query.Select("SUM(amount_cost) total_amount_cost, SUM(amount_fv) total_amount_fv, SUM(unrealized_gain) total_unrealized_gain, SUM(realized_gain) total_relized_gain, SUM(fee) total_fee").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return 0, 0, 0, 0, 0, err
	}

	return data.TotalAmountCost, data.TotalAmountFv, data.TotalUnrealizedGain, data.TotalRealizedGain, data.TotalFee, nil
}

func (r *investasitbkdetail) FindByStock(ctx *abstraction.Context, fmtBridgesID *int, stock *string) (*model.InvestasiTbkDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.InvestasiTbkDetailEntityModel

	query := conn.Model(&model.InvestasiTbkDetailEntityModel{})
	query = query.Where("stock = ?", stock).Where("formatter_bridges_id = ?", fmtBridgesID)
	if err := query.First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
