package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type MutasiDtaDetail interface {
	Find(ctx *abstraction.Context, m *model.MutasiDtaDetailFilterModel, p *abstraction.Pagination) (*[]model.MutasiDtaDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.MutasiDtaDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.MutasiDtaDetailEntityModel) (*model.MutasiDtaDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.MutasiDtaDetailEntityModel) (*model.MutasiDtaDetailEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.MutasiDtaDetailEntityModel) (*model.MutasiDtaDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiDtaDetailEntityModel, error)
	FindWithFormatter(ctx *abstraction.Context, m *model.MutasiDtaDetailFilterModel) (*[]model.MutasiDtaDetailFmtEntityModel, error)
	FindSummary(ctx *abstraction.Context, formatterID, formatterBridgesID *int, sort_id *float64) (*model.MutasiDtaDetailEntityModel, error)
	FindTotal(ctx *abstraction.Context, e *model.MutasiDtaDetailFilterModel) (*model.MutasiDtaDetailEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.MutasiDtaDetailFilterModel) (data *model.MutasiDtaDetailEntityModel, err error)
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

func (r *mutasidtadetail) FindByCriteria(ctx *abstraction.Context, filter *model.MutasiDtaDetailFilterModel) (data *model.MutasiDtaDetailEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.MutasiDtaDetailEntityModel{})
	query = r.Filter(ctx, query, *filter)
	if err = query.First(&data).Error; err != nil {
		return
	}
	// if data.ID == 0 {
	// 	err = errors.New("Data Not Found")
	// }
	return
}
func (r *mutasidtadetail) Find(ctx *abstraction.Context, m *model.MutasiDtaDetailFilterModel, p *abstraction.Pagination) (*[]model.MutasiDtaDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiDtaDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.MutasiDtaDetailEntityModel{})
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
	query = query.Where("formatter_bridges_id IN (?)", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "MUTASI-DTA").Where("trx_ref_id = ?", m.MutasiDtaID))
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

func (r *mutasidtadetail) FindByID(ctx *abstraction.Context, id *int) (*model.MutasiDtaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.MutasiDtaDetailEntityModel
	if err := conn.Where("id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *mutasidtadetail) Create(ctx *abstraction.Context, e *model.MutasiDtaDetailEntityModel) (*model.MutasiDtaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *mutasidtadetail) Update(ctx *abstraction.Context, id *int, e *model.MutasiDtaDetailEntityModel) (*model.MutasiDtaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *mutasidtadetail) Delete(ctx *abstraction.Context, id *int, e *model.MutasiDtaDetailEntityModel) (*model.MutasiDtaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *mutasidtadetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiDtaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.MutasiDtaDetailEntityModel
	if err := conn.Where("code ILIKE ?", *code+"%").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *mutasidtadetail) FindWithFormatter(ctx *abstraction.Context, m *model.MutasiDtaDetailFilterModel) (*[]model.MutasiDtaDetailFmtEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data []model.MutasiDtaDetailFmtEntityModel
	// err := conn.Model(&model.FormatterBridgesEntityModel{}).
	// 	Joins("INNER JOIN m_formatter f ON formatter_bridges.formatter_id = f.id").
	// 	Joins("INNER JOIN m_formatter_detail fd ON f.id = fd.formatter_id").
	// 	Joins("INNER JOIN mutasi_dta_detail tb ON fd.code = tb.code AND tb.formatter_bridges_id = formatter_bridges.id").
	// 	Where("trx_ref_id = ? AND source = ?", m.MutasiDtaID, "MUTASI-DTA").
	// 	Order("fd.sort_id ASC").
	// 	Select("tb.*, auto_summary, is_total, is_control, is_label, control_formula").Find(&data).Error
	// if err != nil {
	// 	return nil, err
	// }

	query := conn.Model(&model.FormatterBridgesEntityModel{}).
		Joins("INNER JOIN m_formatter f ON formatter_bridges.formatter_id = f.id").
		Joins("INNER JOIN m_formatter_detail fd ON f.id = fd.formatter_id").
		Joins("INNER JOIN mutasi_dta_detail tb ON fd.code = tb.code AND tb.formatter_bridges_id = formatter_bridges.id").
		Select("tb.*, auto_summary, is_total, is_control, is_label, control_formula").Order("fd.sort_id ASC")

	if m.FormatterBridgesID != nil && *m.FormatterBridgesID != 0 {
		query = query.Where("formatter_bridges_id = ?", m.FormatterBridgesID)
	}

	if m.MutasiDtaID != nil && *m.MutasiDtaID != 0 {
		query = query.Where("trx_ref_id = ? AND source = ?", m.MutasiDtaID, "MUTASI-DTA")
	}

	err := query.Find(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *mutasidtadetail) FindSummary(ctx *abstraction.Context, formatterID, formatterBridgesID *int, sort_id *float64) (*model.MutasiDtaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.MutasiDtaDetailEntityModel
	var tmp model.FormatterDetailEntityModel
	err := conn.Model(&model.FormatterDetailEntityModel{}).Where("formatter_id = ? AND auto_summary IS NOT NULL AND auto_summary = true AND sort_id < ?", formatterID, sort_id).Order("sort_id DESC").Limit(1).Find(&tmp).Error
	if err != nil {
		return nil, err
	}

	lastID := 0.0
	if tmp.SortID != 0 {
		lastID = tmp.SortID
	}

	err = conn.Raw("SELECT SUM( saldo_awal ) saldo_awal, SUM( manfaat_beban_pajak ) manfaat_beban_pajak, SUM( oci ) oci, SUM( akuisisi_entitas_anak ) akuisisi_entitas_anak, SUM( dibebankan_ke_lr ) dibebankan_ke_lr, SUM( dibebankan_ke_oci ) dibebankan_ke_oci, SUM( saldo_akhir ) saldo_akhir FROM formatter_bridges INNER JOIN m_formatter f ON formatter_bridges.formatter_id = f.ID INNER JOIN m_formatter_detail fd ON f.ID = fd.formatter_id INNER JOIN mutasi_dta_detail tb ON fd.code = tb.code AND tb.formatter_bridges_id = formatter_bridges.ID WHERE ( formatter_bridges_id = ? ) AND fd.sort_id < ? AND fd.sort_id > ?", formatterBridgesID, sort_id, lastID).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *mutasidtadetail) FindTotal(ctx *abstraction.Context, e *model.MutasiDtaDetailFilterModel) (*model.MutasiDtaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.MutasiDtaDetailEntityModel
	query := conn.Model(&model.FormatterBridgesEntityModel{})
	query = query.Where("formatter_bridges_id = ?", e.FormatterBridgesID)
	query = query.Joins("INNER JOIN mutasi_dta_detail mdta ON formatter_bridges.id = mdta.formatter_bridges_id AND code = ?", e.Code)

	if err := query.First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
