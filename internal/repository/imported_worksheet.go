package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type ImportedWorksheet interface {
	Create(ctx *abstraction.Context, e *model.ImportedWorksheetEntityModel) (*model.ImportedWorksheetEntityModel, error)
	FindCompany(ctx *abstraction.Context, m *model.TrialBalanceEntityModel) (*[]model.TrialBalanceEntityModel, error)
	GetVersionsTb(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.ImportedWorksheetEntityModel, error)
	FindByIDDetail(ctx *abstraction.Context, id *int) (*[]model.ImportedWorksheetDetailEntityModel, error)
	Find(ctx *abstraction.Context, m *model.ImportedWorksheetFilterModel, p *abstraction.Pagination) (*[]model.ImportedWorksheetEntityModel, *abstraction.PaginationInfo, error)
	FindByVCTrialBalance(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*model.TrialBalanceEntityModel, error)
	FindByVCAgingUtangPiutang(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel) (*model.AgingUtangPiutangEntityModel, error)
	FindByVCMutasiDta(ctx *abstraction.Context, m *model.MutasiDtaFilterModel) (*model.MutasiDtaEntityModel, error)
	FindByVCMutasiRua(ctx *abstraction.Context, m *model.MutasiRuaFilterModel) (*model.MutasiRuaEntityModel, error)
	FindByVCMutasiIa(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*model.MutasiIaEntityModel, error)
	FindByVCMutasiPersediaan(ctx *abstraction.Context, m *model.MutasiPersediaanFilterModel) (*model.MutasiPersediaanEntityModel, error)
	FindByVCInvestasiTbk(ctx *abstraction.Context, m *model.InvestasiTbkFilterModel) (*model.InvestasiTbkEntityModel, error)
	FindByVCInvestasiNonTbk(ctx *abstraction.Context, m *model.InvestasiNonTbkFilterModel) (*model.InvestasiNonTbkEntityModel, error)
	FindByVCMutasiFa(ctx *abstraction.Context, m *model.MutasiFaFilterModel) (*model.MutasiFaEntityModel, error)
	FindByVCPembelianPenjualanBerelasi(ctx *abstraction.Context, m *model.PembelianPenjualanBerelasiFilterModel) (*model.PembelianPenjualanBerelasiEntityModel, error)
	FindByVCEmployeeBenefit(ctx *abstraction.Context, m *model.EmployeeBenefitFilterModel) (*model.EmployeeBenefitEntityModel, error)
	FindByIDDetailS(ctx *abstraction.Context, id *int) (*[]model.ImportedWorksheetDetailEntityModel, error)
	DeleteAje(ctx *abstraction.Context, companyId *int, tbId *int) (*model.AdjustmentEntityModel, error)
	DeleteJcte(ctx *abstraction.Context, companyId *int, tbId *int) (*model.JcteEntityModel, error)
	DeleteJpm(ctx *abstraction.Context, companyId *int, tbId *int) (*model.JpmEntityModel, error)
	DeleteJelim(ctx *abstraction.Context, companyId *int, tbId *int) (*model.JelimEntityModel, error)
	DeleteTBDF(ctx *abstraction.Context, fbi *int) (*model.TrialBalanceDetailEntityModel, error)
	DeleteTBDFNew(ctx *abstraction.Context, id *int , e *model.TrialBalanceEntityModel) (*model.TrialBalanceEntityModel, error)
	DeleteMFDF(ctx *abstraction.Context, id *int , e *model.MutasiFaEntityModel) (*model.MutasiFaEntityModel, error)
	DeleteMDDF(ctx *abstraction.Context, id *int , e *model.MutasiDtaEntityModel) (*model.MutasiDtaEntityModel, error)
	DeleteMRDF(ctx *abstraction.Context, id *int , e *model.MutasiRuaEntityModel) (*model.MutasiRuaEntityModel, error)
	DeleteMPDF(ctx *abstraction.Context, id *int , e *model.MutasiPersediaanEntityModel) (*model.MutasiPersediaanEntityModel, error)
	DeleteMIDF(ctx *abstraction.Context, id *int , e *model.MutasiIaEntityModel) (*model.MutasiIaEntityModel, error)
	DeleteITDF(ctx *abstraction.Context, id *int , e *model.InvestasiTbkEntityModel) (*model.InvestasiTbkEntityModel, error)
	DeleteINTDF(ctx *abstraction.Context, id *int , e *model.InvestasiNonTbkEntityModel) (*model.InvestasiNonTbkEntityModel, error)
	DeleteAPDF(ctx *abstraction.Context, id *int , e *model.AgingUtangPiutangEntityModel) (*model.AgingUtangPiutangEntityModel, error)
	DeletePPBDF(ctx *abstraction.Context, id *int , e *model.PembelianPenjualanBerelasiEntityModel) (*model.PembelianPenjualanBerelasiEntityModel, error)
	DeleteEBDF(ctx *abstraction.Context, id *int , e *model.EmployeeBenefitEntityModel) (*model.EmployeeBenefitEntityModel, error)
	FindAjeWithTb(ctx *abstraction.Context, id *int) (*[]model.AdjustmentEntityModel, error)
	FindByFormatterBridges(ctx *abstraction.Context, id *int, c *string) (*model.FormatterBridgesEntityModel, error)
	Download(ctx *abstraction.Context, id *int) (*model.ImportedWorksheetDetailEntityModel, error)
	DownloadAll(ctx *abstraction.Context, id *int) (*[]model.ImportedWorksheetDetailEntityModel, error)
	FindByIDWorksheetDetail(ctx *abstraction.Context, id *int) (*model.ImportedWorksheetDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.ImportedWorksheetDetailEntityModel) (*model.ImportedWorksheetDetailEntityModel, error)
	DeleteImportedWorksheet(ctx *abstraction.Context, id *int, e *model.ImportedWorksheetEntityModel) (*model.ImportedWorksheetEntityModel, error)
	DeleteAjeWithTb(ctx *abstraction.Context, id *int, e *model.AdjustmentEntityModel) (*model.AdjustmentEntityModel, error)
}
type importedworksheet struct {
	abstraction.Repository
}

func NewImportedWorksheet(db *gorm.DB) *importedworksheet {
	return &importedworksheet{
		abstraction.Repository{
			Db: db,
		},
	}
}
func (r *importedworksheet) Update(ctx *abstraction.Context, id *int, e *model.ImportedWorksheetDetailEntityModel) (*model.ImportedWorksheetDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *importedworksheet) DownloadAll(ctx *abstraction.Context, id *int) (*[]model.ImportedWorksheetDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.ImportedWorksheetDetailEntityModel
	if err := conn.Where("imported_worksheet_id = ?", &id).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *importedworksheet) Download(ctx *abstraction.Context, id *int) (*model.ImportedWorksheetDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ImportedWorksheetDetailEntityModel
	if err := conn.Where("id = ?", &id).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *importedworksheet) FindByFormatterBridges(ctx *abstraction.Context, id *int, c *string) (*model.FormatterBridgesEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.FormatterBridgesEntityModel
	if err := conn.Where("trx_ref_id = ? AND source = ? ", &id, &c).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *importedworksheet) FindAjeWithTb(ctx *abstraction.Context, fbi *int) (*[]model.AdjustmentEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data *[]model.AdjustmentEntityModel
	if err := conn.Where("tb_id = ?", fbi).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *importedworksheet) DeleteAje(ctx *abstraction.Context, companyId *int, tbId *int) (*model.AdjustmentEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.AdjustmentEntityModel
	if err := conn.Model(data).Where("company_id = ? AND tb_id = ?", companyId, tbId).Delete(data).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *importedworksheet) DeleteJpm(ctx *abstraction.Context, companyId *int, tbId *int) (*model.JpmEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JpmEntityModel
	if err := conn.Model(data).Where("company_id = ? AND tb_id = ?", companyId, tbId).Delete(data).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *importedworksheet) DeleteJcte(ctx *abstraction.Context, companyId *int, tbId *int) (*model.JcteEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JcteEntityModel
	if err := conn.Model(data).Where("company_id = ? AND tb_id = ?", companyId, tbId).Delete(data).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *importedworksheet) DeleteJelim(ctx *abstraction.Context, companyId *int, tbId *int) (*model.JelimEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JelimEntityModel
	if err := conn.Model(data).Where("company_id = ? AND tb_id = ?", companyId, tbId).Delete(data).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return &data, nil
}
func (r *importedworksheet) DeleteImportedWorksheet(ctx *abstraction.Context, id *int, e *model.ImportedWorksheetEntityModel) (*model.ImportedWorksheetEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Update("status", 0).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}
func (r *importedworksheet) DeleteTBDF(ctx *abstraction.Context, fbi *int) (*model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.TrialBalanceDetailEntityModel
	if err := conn.Where("formatter_bridges_id =?", fbi).Delete(data).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *importedworksheet) DeleteAjeWithTb(ctx *abstraction.Context, id *int, e *model.AdjustmentEntityModel) (*model.AdjustmentEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	// e.UserCreatedString = e.UserCreated.Name
	// e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *importedworksheet) DeleteTBDFNew(ctx *abstraction.Context, id *int , e *model.TrialBalanceEntityModel ) (*model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)
	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	// if err := conn.Model(e).Where("id =?", id).Update("status", 4).WithContext(ctx.Request().Context()).Error; err != nil {
	// 	return nil, err
	// }
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}
func (r *importedworksheet) DeleteMFDF(ctx *abstraction.Context, id *int, e *model.MutasiFaEntityModel) (*model.MutasiFaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *importedworksheet) DeleteMDDF(ctx *abstraction.Context, id *int, e *model.MutasiDtaEntityModel) (*model.MutasiDtaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *importedworksheet) DeleteMRDF(ctx *abstraction.Context, id *int, e *model.MutasiRuaEntityModel) (*model.MutasiRuaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *importedworksheet) DeleteMPDF(ctx *abstraction.Context, id *int, e *model.MutasiPersediaanEntityModel ) (*model.MutasiPersediaanEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *importedworksheet) DeleteMIDF(ctx *abstraction.Context, id *int, e *model.MutasiIaEntityModel) (*model.MutasiIaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *importedworksheet) DeleteITDF(ctx *abstraction.Context, id *int, e *model.InvestasiTbkEntityModel) (*model.InvestasiTbkEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *importedworksheet) DeleteINTDF(ctx *abstraction.Context, id *int, e *model.InvestasiNonTbkEntityModel) (*model.InvestasiNonTbkEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *importedworksheet) DeleteAPDF(ctx *abstraction.Context, id *int, e *model.AgingUtangPiutangEntityModel) (*model.AgingUtangPiutangEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *importedworksheet) DeletePPBDF(ctx *abstraction.Context, id *int, e *model.PembelianPenjualanBerelasiEntityModel) (*model.PembelianPenjualanBerelasiEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *importedworksheet) DeleteEBDF(ctx *abstraction.Context, id *int, e *model.EmployeeBenefitEntityModel) (*model.EmployeeBenefitEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *importedworksheet) FindByVCTrialBalance(ctx *abstraction.Context, m *model.TrialBalanceFilterModel ) (*model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.TrialBalanceEntityModel
	if err := conn.Where("DATE(period) = ? AND company_id = ? AND versions = ?",m.Period, m.CompanyID, m.Versions).Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *importedworksheet) FindByVCAgingUtangPiutang(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel) (*model.AgingUtangPiutangEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.AgingUtangPiutangEntityModel
	if err := conn.Where("DATE(period) = ? AND company_id = ? AND versions = ?",m.Period, m.CompanyID, m.Versions).Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}



func (r *importedworksheet) FindByVCMutasiDta(ctx *abstraction.Context, m *model.MutasiDtaFilterModel) (*model.MutasiDtaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.MutasiDtaEntityModel
	if err := conn.Where("DATE(period) = ? AND company_id = ? AND versions = ?",m.Period, m.CompanyID, m.Versions).Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *importedworksheet) FindByVCMutasiRua(ctx *abstraction.Context, m *model.MutasiRuaFilterModel) (*model.MutasiRuaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.MutasiRuaEntityModel
	if err := conn.Where("DATE(period) = ? AND company_id = ? AND versions = ?",m.Period, m.CompanyID, m.Versions).Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *importedworksheet) FindByVCMutasiIa(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*model.MutasiIaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.MutasiIaEntityModel
	if err := conn.Where("DATE(period) = ? AND company_id = ? AND versions = ?",m.Period, m.CompanyID, m.Versions).Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *importedworksheet) FindByVCMutasiPersediaan(ctx *abstraction.Context, m *model.MutasiPersediaanFilterModel) (*model.MutasiPersediaanEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.MutasiPersediaanEntityModel
	if err := conn.Where("DATE(period) = ? AND company_id = ? AND versions = ?",m.Period, m.CompanyID, m.Versions).Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *importedworksheet) FindByVCInvestasiTbk(ctx *abstraction.Context, m *model.InvestasiTbkFilterModel) (*model.InvestasiTbkEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.InvestasiTbkEntityModel
	if err := conn.Where("DATE(period) = ? AND company_id = ? AND versions = ?",m.Period, m.CompanyID, m.Versions).Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *importedworksheet) FindByVCInvestasiNonTbk(ctx *abstraction.Context, m *model.InvestasiNonTbkFilterModel) (*model.InvestasiNonTbkEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.InvestasiNonTbkEntityModel
	if err := conn.Where("DATE(period) = ? AND company_id = ? AND versions = ?",m.Period, m.CompanyID, m.Versions).Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *importedworksheet) FindByVCMutasiFa(ctx *abstraction.Context, m *model.MutasiFaFilterModel) (*model.MutasiFaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.MutasiFaEntityModel
	if err := conn.Where("DATE(period) = ? AND company_id = ? AND versions = ?",m.Period, m.CompanyID, m.Versions).Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *importedworksheet) FindByVCPembelianPenjualanBerelasi(ctx *abstraction.Context, m *model.PembelianPenjualanBerelasiFilterModel) (*model.PembelianPenjualanBerelasiEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.PembelianPenjualanBerelasiEntityModel
	if err := conn.Where("DATE(period) = ? AND company_id = ? AND versions = ?",m.Period, m.CompanyID, m.Versions).Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *importedworksheet) FindByVCEmployeeBenefit(ctx *abstraction.Context, m *model.EmployeeBenefitFilterModel) (*model.EmployeeBenefitEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.EmployeeBenefitEntityModel
	if err := conn.Where("DATE(period) = ? AND company_id = ? AND versions = ?",m.Period, m.CompanyID, m.Versions).Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *importedworksheet) Find(ctx *abstraction.Context, m *model.ImportedWorksheetFilterModel, p *abstraction.Pagination) (*[]model.ImportedWorksheetEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ImportedWorksheetEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.ImportedWorksheetEntityModel{})
	//filter
	tableName := model.ImportedWorksheetEntityModel{}.TableName()
	query = r.FilterTable(ctx, query, *m, tableName)
	query = r.AllowedCompany(ctx, query, tableName)

	//filter custom
	tmp1 := m.CompanyCustomFilter
	queryCompany := conn.Model(&model.CompanyEntityModel{}).Select("id")
	query = r.FilterMultiCompany(ctx, query, queryCompany, tmp1, tableName)
	query = r.FilterMultiVersion(ctx, query, m.ImportedWorksheetFilter)
	queryUser := conn.Model(&model.UserEntityModel{}).Select("id")
	query = r.FilterUser(ctx, query, queryUser, m.Filter, tableName)
	// query = query.Where("status != 0")

	if m.Search != nil {
		query = query.Where("imported_worksheet.created_by IN (?)", conn.Model(&model.UserEntityModel{}).Select("id").Where("name ILIKE ?", "%"+*m.Search+"%"))
	}

	if p.Sort == nil {
		sort := "desc"
		p.Sort = &sort
	}
	if p.SortBy == nil {
		sortBy := "created_at"
		p.SortBy = &sortBy
	}
	tmpSortBy := p.SortBy
	if p.SortBy != nil && *p.SortBy == "company" {
		sortBy := "\"Company\".name"
		p.SortBy = &sortBy
	} else if p.SortBy != nil && *p.SortBy == "user" {
		sortBy := "\"UserCreated\".name"
		p.SortBy = &sortBy
	}
	if p.SortBy != nil && (tmpSortBy != nil && *tmpSortBy != "company" && *tmpSortBy != "user") {
		sortBy := fmt.Sprintf("imported_worksheet.%s", *p.SortBy)
		p.SortBy = &sortBy
	}
	sort := fmt.Sprintf("%s %s", *p.SortBy, *p.Sort)
	query = query.Order(sort)
	p.SortBy = tmpSortBy
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
	var totalData int64
	query = query.Count(&totalData).Limit(limit).Offset(offset)

	if err := query.Joins("Company").Joins("UserCreated", func(db *gorm.DB) *gorm.DB {
		db = db.Select("id, name")
		return db
	}).Preload("UserModified").Preload("ImportedWorksheetDetail").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
	}

	for i, v := range datas {
		datas[i].UserCreatedString = v.UserCreated.Name
		datas[i].UserModifiedString = &v.UserModified.Name
	}

	info.Count = int(totalData)
	info.MoreRecords = false
	if int(totalData) > *p.PageSize {
		info.MoreRecords = true
		// info.Count -= 1
		// datas = datas[:len(datas)-1]
	}

	return &datas, &info, nil
}

func (r *importedworksheet) GetVersionsTb(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.TrialBalanceEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.TrialBalanceEntityModel{})

	// filter
	query = r.Filter(ctx, query, m)

	// sort
	if p.Sort == nil {
		sort := "desc"
		p.Sort = &sort
	}
	if p.SortBy == nil {
		sortBy := "created_at"
		p.SortBy = &sortBy
	}
	sort := fmt.Sprintf("%s %s", *p.SortBy, *p.Sort)
	query = query.Order(sort)

	// pagination
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
	var totalData int64
	query = query.Count(&totalData).Limit(limit).Offset(offset)

	err := query.Where("company_id = ?", m.CompanyID).Find(&datas).
		WithContext(ctx.Request().Context()).Error
	conn.Find(&datas)
	if err != nil {
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

func (r *importedworksheet) Create(ctx *abstraction.Context, e *model.ImportedWorksheetEntityModel) (*model.ImportedWorksheetEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Preload("UserCreated").Preload("UserModified").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *importedworksheet) FindCompany(ctx *abstraction.Context, m *model.TrialBalanceEntityModel) (*[]model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.TrialBalanceEntityModel
	if err := conn.Where("company_id = ? AND period = ?", m.CompanyID, m.Period).Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, err
	}
	return &datas, nil
}
func (r *importedworksheet) FindByID(ctx *abstraction.Context, id *int) (*model.ImportedWorksheetEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ImportedWorksheetEntityModel
	if err := conn.Where("id = ?", &id).Preload("Company").Preload("ImportedWorksheetDetail").Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	data.UserCreatedString = data.UserCreated.Name
	data.UserModifiedString = &data.UserModified.Name
	return &data, nil
}

func (r *importedworksheet) FindByIDWorksheetDetail(ctx *abstraction.Context, id *int) (*model.ImportedWorksheetDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ImportedWorksheetDetailEntityModel
	if err := conn.Where("id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *importedworksheet) FindByIDDetail(ctx *abstraction.Context, id *int) (*[]model.ImportedWorksheetDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.ImportedWorksheetDetailEntityModel
	if err := conn.Where("imported_worksheet_id = ? AND status = ?", &id, 2).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *importedworksheet) FindByIDDetailS(ctx *abstraction.Context, id *int) (*[]model.ImportedWorksheetDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.ImportedWorksheetDetailEntityModel
	if err := conn.Where("imported_worksheet_id = ? AND status = ?", &id, 1).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
