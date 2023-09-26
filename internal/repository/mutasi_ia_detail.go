package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type MutasiIaDetail interface {
	Find(ctx *abstraction.Context, m *model.MutasiIaDetailFilterModel, p *abstraction.Pagination) (*[]model.MutasiIaDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.MutasiIaDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.MutasiIaDetailEntityModel) (*model.MutasiIaDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.MutasiIaDetailEntityModel) (*model.MutasiIaDetailEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.MutasiIaDetailEntityModel) (*model.MutasiIaDetailEntityModel, error)
	FindByCode(ctx *abstraction.Context, fmtBridgesID *int, code *string) (*model.MutasiIaDetailEntityModel, error)
	FindWithFormatter(ctx *abstraction.Context, m *model.MutasiIaDetailFilterModel) (*[]model.MutasiIaDetailFmtEntityModel, error)
	FindSummary(ctx *abstraction.Context, formatterID, formatterBridgesID *int, sort_id *float64) (*model.MutasiIaDetailEntityModel, error)
	FindTotal(ctx *abstraction.Context, e *model.MutasiIaDetailFilterModel) (*model.MutasiIaDetailEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.MutasiIaDetailFilterModel) (data *model.MutasiIaDetailEntityModel, err error)
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

func (r *mutasiiadetail) FindByCriteria(ctx *abstraction.Context, filter *model.MutasiIaDetailFilterModel) (data *model.MutasiIaDetailEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.MutasiIaDetailEntityModel{})
	query = r.Filter(ctx, query, *filter)
	if err = query.First(&data).Error; err != nil {
		return
	}
	// if data.ID == 0 {
	// 	err = errors.New("Data Not Found")
	// }
	return
}
func (r *mutasiiadetail) Find(ctx *abstraction.Context, m *model.MutasiIaDetailFilterModel, p *abstraction.Pagination) (*[]model.MutasiIaDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiIaDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.MutasiIaDetailEntityModel{})
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
	query = query.Where("formatter_bridges_id IN (?)", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "MUTASI-IA").Where("trx_ref_id = ?", m.MutasiIaID))
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

func (r *mutasiiadetail) FindByID(ctx *abstraction.Context, id *int) (*model.MutasiIaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.MutasiIaDetailEntityModel
	if err := conn.Where("id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *mutasiiadetail) Create(ctx *abstraction.Context, e *model.MutasiIaDetailEntityModel) (*model.MutasiIaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *mutasiiadetail) Update(ctx *abstraction.Context, id *int, e *model.MutasiIaDetailEntityModel) (*model.MutasiIaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *mutasiiadetail) Delete(ctx *abstraction.Context, id *int, e *model.MutasiIaDetailEntityModel) (*model.MutasiIaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *mutasiiadetail) FindByCode(ctx *abstraction.Context, fmtBridgesID *int, code *string) (*model.MutasiIaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.MutasiIaDetailEntityModel
	if err := conn.Where("formatter_bridges_id = ?", *fmtBridgesID).Where("code ILIKE ?", *code).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *mutasiiadetail) FindWithFormatter(ctx *abstraction.Context, m *model.MutasiIaDetailFilterModel) (*[]model.MutasiIaDetailFmtEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data []model.MutasiIaDetailFmtEntityModel
	query := conn.Model(&model.FormatterBridgesEntityModel{}).
		Joins("INNER JOIN m_formatter f ON formatter_bridges.formatter_id = f.id").
		Joins("INNER JOIN m_formatter_detail fd ON f.id = fd.formatter_id").
		Joins("INNER JOIN mutasi_ia_detail tb ON fd.code = tb.code AND tb.formatter_bridges_id = formatter_bridges.id").
		Select("tb.*, auto_summary, is_total, is_control, is_label, control_formula").Order("fd.sort_id ASC")

	if m.FormatterBridgesID != nil && *m.FormatterBridgesID != 0 {
		query = query.Where("formatter_bridges_id = ?", m.FormatterBridgesID)
	}

	if m.MutasiIaID != nil && *m.MutasiIaID != 0 {
		query = query.Where("trx_ref_id = ? AND source = ?", m.MutasiIaID, "MUTASI-IA")
	}

	err := query.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *mutasiiadetail) FindSummary(ctx *abstraction.Context, formatterID, formatterBridgesID *int, sort_id *float64) (*model.MutasiIaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.MutasiIaDetailEntityModel
	var tmp model.FormatterDetailEntityModel
	err := conn.Model(&model.FormatterDetailEntityModel{}).Where("formatter_id = ? AND auto_summary IS NOT NULL AND auto_summary = true AND sort_id < ?", formatterID, sort_id).Order("sort_id DESC").Limit(1).Find(&tmp).Error
	if err != nil {
		return nil, err
	}

	lastID := 0.0
	if tmp.SortID != 0 {
		lastID = tmp.SortID
	}

	err = conn.Raw("SELECT SUM( beginning_balance ) beginning_balance, SUM( acquisition_of_subsidiary ) acquisition_of_subsidiary, SUM( additions ) additions, SUM( deductions ) deductions, SUM( reclassification ) reclassification, SUM( revaluation ) revaluation, SUM( ending_balance ) ending_balance FROM (SELECT tb.* FROM formatter_bridges INNER JOIN m_formatter f ON formatter_bridges.formatter_id = f.ID INNER JOIN m_formatter_detail fd ON f.ID = fd.formatter_id INNER JOIN mutasi_ia_detail tb ON fd.code = tb.code AND tb.formatter_bridges_id = formatter_bridges.ID WHERE ( formatter_bridges_id = ? ) AND fd.sort_id < ? AND fd.sort_id > ? ORDER BY fd.sort_id ASC) tmp ", formatterBridgesID, sort_id, lastID).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *mutasiiadetail) FindTotal(ctx *abstraction.Context, e *model.MutasiIaDetailFilterModel) (*model.MutasiIaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.MutasiIaDetailEntityModel
	query := conn.Model(&model.MutasiIaDetailEntityModel{})
	query = query.Where("source = ? AND trx_ref_id = ?", "MUTASI-IA", e.MutasiIaID)
	query = query.Joins("INNER JOIN formatter_bridges ON formatter_bridges.id = formatter_bridges_id AND code = ?", e.Code)

	if err := query.First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
