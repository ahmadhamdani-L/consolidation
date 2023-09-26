package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type ConsolidationDetail interface {
	Find(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel, p *abstraction.Pagination) ([]model.ConsolidationDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.ConsolidationDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.ConsolidationDetailEntityModel) (*model.ConsolidationDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.ConsolidationDetailEntityModel) (*model.ConsolidationDetailEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.ConsolidationDetailEntityModel) (*model.ConsolidationDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.ConsolidationDetailEntityModel, error)
	Import(ctx *abstraction.Context, e *[]model.ConsolidationDetailEntityModel) (*[]model.ConsolidationDetailEntityModel, error)
	FindByFormatter(ctx *abstraction.Context, formatter *string) (*model.FormatterEntityModel, error)
	FindByFormatterDetail(ctx *abstraction.Context, id *int, parent *string) (*model.FormatterDetailEntityModel, error)
	FindAnakUsaha(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel) ([]model.ConsolidationBridgeEntityModel, error)
	FindCode(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel, p *abstraction.Pagination) ([]model.ConsolidationDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByAmount(ctx *abstraction.Context, id *int, code *string) (*[]model.ConsolidationBridgeDetailEntityModel, error)
	
	FindAnakUsahaOnly(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel) ([]model.ConsolidationBridgeEntityModel, error)
	FindAnakUsahaOnlys(ctx *abstraction.Context, consolidationID int , code string) ([]model.ConsolidationBridgeEntityModel, error)
	FindisNull(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel, p *abstraction.Pagination) ([]model.ConsolidationDetailEntityModel,  *abstraction.PaginationInfo, error)
	FindAnakUsahaOnlyNull(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel) ([]model.ConsolidationBridgeEntityModel, error)
	FindDetail(ctx *abstraction.Context, consolidationID *int, parentID *int) (*[]model.ConsolidationDetailFmtEntityModel, error)
	FindAllDetail(ctx *abstraction.Context, consolidationID *int) (*[]model.ConsolidationDetailFmtEntityModel, error)
	FindAllDetailAmounts(ctx *abstraction.Context, consolidationID *int, code *string) (*[]model.ConsolidationBridgeDetailEntityModel, error)
	FindAnakUsahaOnlyss(ctx *abstraction.Context, consolidationID int) ([]model.ConsolidationBridgeEntityModel, error) 
}

type consolidationdetail struct {
	abstraction.Repository
}

func NewConsolidationDetail(db *gorm.DB) *consolidationdetail {
	return &consolidationdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *consolidationdetail) FindAnakUsahaOnlyss(ctx *abstraction.Context, consolidationID int) ([]model.ConsolidationBridgeEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationBridgeEntityModel

	query := conn.Model(&model.ConsolidationBridgeEntityModel{})
	if err := query.Where("consolidation_id = ? ", consolidationID).Preload("Company").
	Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return datas, nil
	}
	
	// querycount := query
	// var totalData int64
	// if err := querycount.Count(&totalData).Error; err != nil {
	// 	return &datas, &info, err
	// }
	// if *p.PageSize != -1 {
	// 	query = query.Limit(limit).Offset(offset)
	// }
	// if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
	// 	return &datas, &info, err
	// }

	// info.Count = int(totalData)
	// info.MoreRecords = false
	// if len(datas) > *p.PageSize {
	// 	info.MoreRecords = true
	// 	info.Count -= 1
	// 	datas = datas[:len(datas)-1]
	// }

	return datas, nil
}

func (r *consolidationdetail) FindAllDetailAmounts(ctx *abstraction.Context, consolidationID *int, code *string) (*[]model.ConsolidationBridgeDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	// if err := conn.Model(&model.FormatterBridgesEntityModel{}).Where("trx_ref_id = ? AND source = 'TRIAL-BALANCE'", tb.ID).First(&fmtBridgeData).Error; err != nil {
	// 	return nil, err
	// }

	query := conn.Model(&model.ConsolidationBridgeDetailEntityModel{})
	query = query.Where("consolidation_bridge_id = ? AND code = ? ", consolidationID , code)
	var data []model.ConsolidationBridgeDetailEntityModel
	if err := query.First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}
func (r *consolidationdetail) FindAllDetail(ctx *abstraction.Context, consolidationID *int) (*[]model.ConsolidationDetailFmtEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var tb model.ConsolidationEntityModel
	// var fmtBridgeData model.FormatterBridgesEntityModel

	if err := conn.Model(&model.ConsolidationEntityModel{}).Where("id = ?", consolidationID).First(&tb).Error; err != nil {
		return nil, err
	}

	// if err := conn.Model(&model.FormatterBridgesEntityModel{}).Where("trx_ref_id = ? AND source = 'TRIAL-BALANCE'", tb.ID).First(&fmtBridgeData).Error; err != nil {
	// 	return nil, err
	// }
	tmpCon1 := "SUBSTR( consolidation_detail.code, 0, 7 ) = fd.code"
	tmpCon2 := "SUBSTR( consolidation_detail.code, 0, 4 ) = fd.code"

	querySelect := "SELECT consolidation_detail.id, fd.id as formatter_detail_id, consolidation_detail.consolidation_id,consolidation_detail.code,consolidation_detail.wp_reff,consolidation_detail.description,consolidation_detail.sort_id,consolidation_detail.amount_before_jpm,consolidation_detail.amount_jpm_dr,consolidation_detail.amount_jpm_cr,consolidation_detail.amount_after_jpm,consolidation_detail.amount_jcte_dr,consolidation_detail.amount_jcte_cr,consolidation_detail.amount_after_jcte,consolidation_detail.amount_combine_subsidiary,consolidation_detail.amount_jelim_dr,consolidation_detail.amount_jelim_cr,consolidation_detail.amount_console, fd.is_parent, parent_id, fd.show_group_coa"

	tmpTB := " FROM \"consolidation_detail\" INNER JOIN m_formatter_detail fd ON fd.formatter_id = 3 AND"

	queryStr := fmt.Sprintf("WITH tbxfmt AS (%s %s (consolidation_detail.code = fd.code OR "+tmpCon1+") WHERE consolidation_id = ?), tbxfmtc AS (%s %s (consolidation_detail.code = fd.code OR "+tmpCon2+") WHERE consolidation_id = ? AND fd.code IN ('950', '960'))", querySelect, tmpTB, querySelect, tmpTB)

	query := conn.Raw(fmt.Sprintf("%s SELECT * FROM tbxfmt UNION ALL SELECT * FROM tbxfmtc ORDER BY sort_id, id", queryStr), consolidationID, consolidationID)

	// query := conn.Model(&model.ConsolidationDetailEntityModel{}).Joins("INNER JOIN m_formatter_detail fd ON (consolidation_detail.code = fd.code OR SUBSTR(consolidation_detail.code, 0, 7) = fd.code) AND fd.formatter_id = ?", 3)
	// query = query.Select("consolidation_detail.id, fd.id as formatter_detail_id, consolidation_detail.consolidation_id,consolidation_detail.code,consolidation_detail.wp_reff,consolidation_detail.description,consolidation_detail.sort_id,consolidation_detail.amount_before_jpm,consolidation_detail.amount_jpm_dr,consolidation_detail.amount_jpm_cr,consolidation_detail.amount_after_jpm,consolidation_detail.amount_jcte_dr,consolidation_detail.amount_jcte_cr,consolidation_detail.amount_after_jcte,consolidation_detail.amount_combine_subsidiary,consolidation_detail.amount_jelim_dr,consolidation_detail.amount_jelim_cr,consolidation_detail.amount_console, fd.is_parent, parent_id")
	// query = query.Where("consolidation_id = ?", consolidationID).Order("fd.sort_id ASC")
	var data []model.ConsolidationDetailFmtEntityModel
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}
func (r *consolidationdetail) FindDetail(ctx *abstraction.Context, consolidationID *int, parentID *int) (*[]model.ConsolidationDetailFmtEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var tb model.ConsolidationEntityModel
	// var fmtBridgeData model.FormatterBridgesEntityModel

	if err := conn.Model(&model.ConsolidationEntityModel{}).Where("id = ?", consolidationID).Find(&tb).Error; err != nil {
		return nil, err
	}
	// if err := conn.Model(&model.FormatterBridgesEntityModel{}).Where("trx_ref_id = ? AND source = 'TRIAL-BALANCE'", tb.ID).Find(&fmtBridgeData).Error; err != nil {
	// 	return nil, err
	// }
	query := conn.Model(&model.ConsolidationDetailEntityModel{}).Joins("INNER JOIN m_formatter_detail fd ON (consolidation_detail.code = fd.code OR SUBSTR(consolidation_detail.code, 0, 7) = fd.code) AND fd.formatter_id = 3")
	query = query.Select("consolidation_detail.id, fd.id as formatter_detail_id, consolidation_detail.consolidation_id,consolidation_detail.code,consolidation_detail.wp_reff,consolidation_detail.description,consolidation_detail.sort_id,consolidation_detail.amount_before_jpm,consolidation_detail.amount_jpm_dr,consolidation_detail.amount_jpm_cr,consolidation_detail.amount_after_jpm,consolidation_detail.amount_jcte_dr,consolidation_detail.amount_jcte_cr,consolidation_detail.amount_after_jcte,consolidation_detail.amount_combine_subsidiary,consolidation_detail.amount_jelim_dr,consolidation_detail.amount_jelim_cr,consolidation_detail.amount_console,fd.is_parent, parent_id")
	if parentID != nil && *parentID != 0 {
		query = query.Where("parent_id = ?", parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	query = query.Where("consolidation_id = ?", tb.ID).Order("fd.sort_id ASC")
	var data []model.ConsolidationDetailFmtEntityModel
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}
func (r *consolidationdetail) FindAnakUsahaOnlyNull(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel) ([]model.ConsolidationBridgeEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationBridgeEntityModel

	query := conn.Model(&model.ConsolidationBridgeEntityModel{})
	if err := query.Where("consolidation_detail.consolidation_id = ? AND consolidation.consolidation_versions = ? AND m_formatter_detail.parent_id is null", m.ConsolidationID, m.VersionConsolidation).
	Joins("INNER JOIN consolidation  ON consolidation_bridge.consolidation_versions = consolidation.consolidation_versions").
	Joins("INNER JOIN consolidation_detail consolidation_detail ON consolidation_bridge.consolidation_id = consolidation_detail.consolidation_id").
	Joins("INNER JOIN m_formatter_detail ON m_formatter_detail.sort_id = consolidation_detail.sort_id").
	Joins("LEFT JOIN consolidation_bridge_detail ON consolidation_bridge.ID = consolidation_bridge_detail.consolidation_bridge_id AND consolidation_bridge_detail.code = consolidation_detail.code").
	Group("consolidation_bridge.id").Preload("ConsolidationBridgeDetail", "code = ?", m.Code).
	Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return datas, nil
	}
	
	// querycount := query
	// var totalData int64
	// if err := querycount.Count(&totalData).Error; err != nil {
	// 	return &datas, &info, err
	// }
	// if *p.PageSize != -1 {
	// 	query = query.Limit(limit).Offset(offset)
	// }
	// if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
	// 	return &datas, &info, err
	// }

	// info.Count = int(totalData)
	// info.MoreRecords = false
	// if len(datas) > *p.PageSize {
	// 	info.MoreRecords = true
	// 	info.Count -= 1
	// 	datas = datas[:len(datas)-1]
	// }

	return datas, nil
}
func (r *consolidationdetail) FindByAmount(ctx *abstraction.Context, id *int, code *string) (*[]model.ConsolidationBridgeDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.ConsolidationBridgeDetailEntityModel
	if err := conn.Where("consolidation_bridge_id = ? AND code = ? ", &id, &code).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *consolidationdetail) FindCode(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel, p *abstraction.Pagination) ([]model.ConsolidationDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.ConsolidationDetailEntityModel{})

	//filter
	// query = r.Filter(ctx, query, *m)

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
	// limit := *p.PageSize + 1
	// offset := limit * (*p.Page - 1)

	if err := query.Where("consolidation_id = ? AND m_formatter_detail.parent_id = ? ", m.ConsolidationID, m.ParentID).
		Joins(fmt.Sprintf("INNER JOIN m_formatter_detail ON m_formatter_detail.sort_id = consolidation_detail.sort_id")).Group("consolidation_detail.id").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return datas, &info, err
	}

	// querycount := query
	// var totalData int64
	// if err := querycount.Count(&totalData).Error; err != nil {
	// 	return &datas, &info, err
	// }
	// if *p.PageSize != -1 {
	// 	query = query.Limit(limit).Offset(offset)
	// }
	// if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
	// 	return &datas, &info, err
	// }

	// info.Count = int(totalData)
	// info.MoreRecords = false
	// if len(datas) > *p.PageSize {
	// 	info.MoreRecords = true
	// 	info.Count -= 1
	// 	datas = datas[:len(datas)-1]
	// }

	return datas, &info, nil
}

func (r *consolidationdetail) Find(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel, p *abstraction.Pagination) ([]model.ConsolidationDetailEntityModel,  *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.ConsolidationDetailEntityModel{})

	//filter
	// query = r.Filter(ctx, query, *m)

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
	query2 := fmt.Sprintf("SELECT distinct cd.id, cd.consolidation_id, cd.code, cd.wp_reff, cd.description, cd.sort_id,cd.amount_before_jpm, cd.amount_jpm_cr, cd.amount_jpm_dr, cd.amount_after_jpm, cd.amount_jcte_cr, cd.amount_jcte_dr, cd.amount_after_jcte,cd,amount_combine_subsidiary, cd.amount_jelim_cr, cd.amount_jelim_dr, cd.amount_console , mfd.is_parent FROM consolidation_detail cd INNER JOIN m_formatter_detail mfd ON (cd.code = mfd.code OR SUBSTR(cd.code, 0, 7) = mfd.code) WHERE cd.consolidation_id = ? AND mfd.parent_id = ? ")
	if err := conn.Raw(query2, m.ConsolidationID, m.ParentID, m.ConsolidationID).Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return datas, &info, err
	}

	// querycount := query
	// var totalData int64
	// if err := querycount.Count(&totalData).Error; err != nil {
	// 	return &datas, &info, err
	// }
	// if *p.PageSize != -1 {
	// 	query = query.Limit(limit).Offset(offset)
	// }
	// if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
	// 	return &datas, &info, err
	// }

	// info.Count = int(totalData)
	// info.MoreRecords = false
	// if len(datas) > *p.PageSize {
	// 	info.MoreRecords = true
	// 	info.Count -= 1
	// 	datas = datas[:len(datas)-1]
	// }

	return datas,&info, nil
}

func (r *consolidationdetail) FindisNull(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel, p *abstraction.Pagination) ([]model.ConsolidationDetailEntityModel,  *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.ConsolidationDetailEntityModel{})

	//filter
	// query = r.Filter(ctx, query, *m)

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
	query2 := fmt.Sprintf("SELECT distinct cd.id, cd.consolidation_id, cd.code, cd.wp_reff, cd.description, cd.sort_id,cd.amount_before_jpm, cd.amount_jpm_cr, cd.amount_jpm_dr, cd.amount_after_jpm, cd.amount_jcte_cr, cd.amount_jcte_dr, cd.amount_after_jcte,cd,amount_combine_subsidiary, cd.amount_jelim_cr, cd.amount_jelim_dr, cd.amount_console , cd.is_parent FROM consolidation_detail cd INNER JOIN m_formatter_detail mfd ON (cd.code = mfd.code OR SUBSTR(cd.code, 0, 7) = mfd.code) WHERE cd.consolidation_id = ? AND mfd.parent_id is null ")
	if err := conn.Raw(query2, m.ConsolidationID, m.ConsolidationID).Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return datas, &info, err
	}

	// querycount := query
	// var totalData int64
	// if err := querycount.Count(&totalData).Error; err != nil {
	// 	return &datas, &info, err
	// }
	// if *p.PageSize != -1 {
	// 	query = query.Limit(limit).Offset(offset)
	// }
	// if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
	// 	return &datas, &info, err
	// }

	// info.Count = int(totalData)
	// info.MoreRecords = false
	// if len(datas) > *p.PageSize {
	// 	info.MoreRecords = true
	// 	info.Count -= 1
	// 	datas = datas[:len(datas)-1]
	// }

	return datas,&info, nil
}

func (r *consolidationdetail) FindAnakUsaha(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel) ([]model.ConsolidationBridgeEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationBridgeEntityModel

	query := conn.Model(&model.ConsolidationBridgeEntityModel{}).Preload("ConsolidationBridgeDetail", "code = ?", m.Code).Preload("Company")
	if err := query.Where("consolidation_detail.consolidation_id = ? AND consolidation.consolidation_versions = ? AND m_formatter_detail.parent_id = ? AND consolidation_detail.code = ? ", m.ConsolidationID, m.VersionConsolidation, m.ParentID, m.Code).
	Joins(fmt.Sprintf("INNER JOIN consolidation  ON consolidation_bridge.consolidation_versions = consolidation.consolidation_versions")).
	Joins(fmt.Sprintf("INNER JOIN consolidation_detail consolidation_detail ON consolidation_bridge.consolidation_id = consolidation_detail.consolidation_id")).
	Joins(fmt.Sprintf("INNER JOIN m_formatter_detail ON m_formatter_detail.sort_id = consolidation_detail.sort_id")).
	Joins(fmt.Sprintf("LEFT JOIN consolidation_bridge_detail ON consolidation_bridge.ID = consolidation_bridge_detail.consolidation_bridge_id AND consolidation_bridge_detail.code = consolidation_detail.code")).
	Group("consolidation_bridge.id").
	Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return datas, nil
	}
	
	// querycount := query
	// var totalData int64
	// if err := querycount.Count(&totalData).Error; err != nil {
	// 	return &datas, &info, err
	// }
	// if *p.PageSize != -1 {
	// 	query = query.Limit(limit).Offset(offset)
	// }
	// if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
	// 	return &datas, &info, err
	// }

	// info.Count = int(totalData)
	// info.MoreRecords = false
	// if len(datas) > *p.PageSize {
	// 	info.MoreRecords = true
	// 	info.Count -= 1
	// 	datas = datas[:len(datas)-1]
	// }

	return datas, nil
}
func (r *consolidationdetail) FindAnakUsahaOnly(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel) ([]model.ConsolidationBridgeEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationBridgeEntityModel

	query := conn.Model(&model.ConsolidationBridgeEntityModel{})
	if err := query.Where("consolidation_detail.consolidation_id = ? ", m.ConsolidationID).
	Joins(fmt.Sprintf("INNER JOIN consolidation  ON consolidation_bridge.consolidation_versions = consolidation.consolidation_versions")).
	Joins(fmt.Sprintf("INNER JOIN consolidation_detail consolidation_detail ON consolidation_bridge.consolidation_id = consolidation_detail.consolidation_id")).
	Joins(fmt.Sprintf("INNER JOIN m_formatter_detail ON m_formatter_detail.sort_id = consolidation_detail.sort_id")).
	Joins(fmt.Sprintf("LEFT JOIN consolidation_bridge_detail ON consolidation_bridge.ID = consolidation_bridge_detail.consolidation_bridge_id AND consolidation_bridge_detail.code = consolidation_detail.code")).
	Group("consolidation_bridge.id").Preload("ConsolidationBridgeDetail", "code = ?", m.Code).Preload("Company").
	Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return datas, nil
	}
	
	// querycount := query
	// var totalData int64
	// if err := querycount.Count(&totalData).Error; err != nil {
	// 	return &datas, &info, err
	// }
	// if *p.PageSize != -1 {
	// 	query = query.Limit(limit).Offset(offset)
	// }
	// if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
	// 	return &datas, &info, err
	// }

	// info.Count = int(totalData)
	// info.MoreRecords = false
	// if len(datas) > *p.PageSize {
	// 	info.MoreRecords = true
	// 	info.Count -= 1
	// 	datas = datas[:len(datas)-1]
	// }

	return datas, nil
}
func (r *consolidationdetail) FindAnakUsahaOnlys(ctx *abstraction.Context, consolidationID int , code string) ([]model.ConsolidationBridgeEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationBridgeEntityModel

	query := conn.Model(&model.ConsolidationBridgeEntityModel{})
	if err := query.Where("consolidation_id = ? ", consolidationID).
	// Joins(fmt.Sprintf("INNER JOIN consolidation  ON consolidation_bridge.consolidation_versions = consolidation.consolidation_versions")).
	// Joins("INNER JOIN consolidation_detail consolidation_detail ON consolidation_bridge.consolidation_id = consolidation_detail.consolidation_id").
	// Joins("INNER JOIN m_formatter_detail ON m_formatter_detail.sort_id = consolidation_detail.sort_id").
	// Joins(fmt.Sprintf("LEFT JOIN consolidation_bridge_detail ON consolidation_bridge.ID = consolidation_bridge_detail.consolidation_bridge_id AND consolidation_bridge_detail.code = consolidation_detail.code")).
	Preload("ConsolidationBridgeDetail", "code = ?", code).Preload("Company").
	Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return datas, nil
	}
	
	// querycount := query
	// var totalData int64
	// if err := querycount.Count(&totalData).Error; err != nil {
	// 	return &datas, &info, err
	// }
	// if *p.PageSize != -1 {
	// 	query = query.Limit(limit).Offset(offset)
	// }
	// if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
	// 	return &datas, &info, err
	// }

	// info.Count = int(totalData)
	// info.MoreRecords = false
	// if len(datas) > *p.PageSize {
	// 	info.MoreRecords = true
	// 	info.Count -= 1
	// 	datas = datas[:len(datas)-1]
	// }

	return datas, nil
}
func (r *consolidationdetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.ConsolidationDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code LIKE ?", tmp+"%").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *consolidationdetail) FindByID(ctx *abstraction.Context, id *int) (*model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ConsolidationDetailEntityModel
	if err := conn.Where("id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *consolidationdetail) FindByFormatter(ctx *abstraction.Context, formatter *string) (*model.FormatterEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.FormatterEntityModel
	if err := conn.Where("formatter_for = ?", &formatter).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *consolidationdetail) FindByFormatterDetail(ctx *abstraction.Context, id *int, parent *string) (*model.FormatterDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.FormatterDetailEntityModel
	if err := conn.Where("formatter_id = ? AND code = ?", &id, &parent).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *consolidationdetail) Create(ctx *abstraction.Context, e *model.ConsolidationDetailEntityModel) (*model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *consolidationdetail) Update(ctx *abstraction.Context, id *int, e *model.ConsolidationDetailEntityModel) (*model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *consolidationdetail) Delete(ctx *abstraction.Context, id *int, e *model.ConsolidationDetailEntityModel) (*model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *consolidationdetail) Import(ctx *abstraction.Context, e *[]model.ConsolidationDetailEntityModel) (*[]model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}
