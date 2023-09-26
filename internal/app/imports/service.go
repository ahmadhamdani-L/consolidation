package imports

import (
	"errors"
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"mcash-finance-console-core/pkg/util/trxmanager"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx *abstraction.Context, payload *dto.ImportedWorksheetCreateRequest) (*dto.ImportedWorksheetCreateResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.ImportReUploadRequest) (*dto.ImportReUploadResponse, error)
	FindByIDDetail(ctx *abstraction.Context, payload *dto.ImportReUploadDetailRequest) (*dto.ImportReUploadDetailResponse, error)
	FindByIDDetailS(ctx *abstraction.Context, id *int) (*dto.ImportReUploadDetailResponse, error)
	DeleteTBDF(ctx *abstraction.Context, fbi *int) (*model.TrialBalanceDetailEntityModel, error)
	FindByFormatterBridges(ctx *abstraction.Context, fbi *int, source *string) (*model.FormatterBridgesEntityModel, error)
	DeletedFormatterBridges(ctx *abstraction.Context, trx *int, source *string) (*model.FormatterBridgesEntityModel, error)
	DeletedAJE(ctx *abstraction.Context, companyId *int, tbId *int) (*model.AdjustmentEntityModel, error)
	DeletedJELIM(ctx *abstraction.Context, companyId *int, tbId *int) (*model.JelimEntityModel, error)
	DeletedJPM(ctx *abstraction.Context, companyId *int, tbId *int) (*model.JpmEntityModel, error)
	DeletedJCTE(ctx *abstraction.Context, companyId *int, tbId *int) (*model.JcteEntityModel, error)
	Download(ctx *abstraction.Context, payload *dto.ImportedWorksheetDetailGetByIDRequest) (*dto.ImportedWorksheetDetailGetByIDResponse, error)
}

type service struct {
	Repository                                 repository.ImportedWorksheet
	Db                                         *gorm.DB
	CompanyRepository                          repository.Company
	NotificationRepository                     repository.Notification
	ImportedWorksheetRepository                repository.ImportedWorksheet
	ConsolidationRepository                    repository.Consolidation
	AgingUtangPiutangRepository                repository.AgingUtangPiutang
	AgingUtangPiutangDetailRepository          repository.AgingUtangPiutangDetail
	FormatterRepository                        repository.Formatter
	FormatterBridgesRepository                 repository.FormatterBridges
	EmployeeBenefitRepository                  repository.EmployeeBenefit
	EmployeeBenefitDetailRepository            repository.EmployeeBenefitDetail
	AdjustmentRepository                       repository.Adjustment
	AdjustmentDetailRepository                 repository.AdjustmentDetail
	InvestasiNonTbkRepository                  repository.InvestasiNonTbk
	InvestasiNonTbkDetailRepository            repository.InvestasiNonTbkDetail
	InvestasiTbkRepository                     repository.InvestasiTbk
	InvestasiTbkDetailRepository               repository.InvestasiTbkDetail
	MutasiDtaRepository                        repository.MutasiDta
	MutasiDtaDetailRepository                  repository.MutasiDtaDetail
	MutasiFaRepository                         repository.MutasiFa
	MutasiFaDetailRepository                   repository.MutasiFaDetail
	MutasiIaRepository                         repository.MutasiIa
	MutasiIaDetailRepository                   repository.MutasiIaDetail
	MutasiPersediaanRepository                 repository.MutasiPersediaan
	MutasiPersediaanDetailRepository           repository.MutasiPersediaanDetail
	MutasiRuaRepository                        repository.MutasiRua
	MutasiRuaDetailRepository                  repository.MutasiRuaDetail
	PembelianPenjualanBerelasiRepository       repository.PembelianPenjualanBerelasi
	PembelianPenjualanBerelasiDetailRepository repository.PembelianPenjualanBerelasiDetail
	TrialBalanceRepository                     repository.TrialBalance
	TrialBalanceDetailRepository               repository.TrialBalanceDetail
	FormatterDetailRepository                  repository.FormatterDetail
	CoaRepository                              repository.Coa
}

func NewService(f *factory.Factory) *service {
	repository := f.ImportedWorksheetRepository
	db := f.Db
	notifRepo := f.NotificationRepository
	importedWorksheetRepo := f.ImportedWorksheetRepository
	consolidationRepo := f.ConsolidationRepository
	agingUtangPiutangRepo := f.AgingUtangPiutangRepository
	agingUtangPiutangDetailRepo := f.AgingUtangPiutangDetailRepository
	formatterRepo := f.FormatterRepository
	formatterBridgesRepo := f.FormatterBridgesRepository
	employeeBenefitRepo := f.EmployeeBenefitRepository
	employeeBenefitDetailRepo := f.EmployeeBenefitDetailRepository
	adjustmentRepo := f.AdjustmentRepository
	adjustmentDetailRepo := f.AdjustmentDetailRepository
	investasiNonTbkRepo := f.InvestasiNonTbkRepository
	investasiNonTbkDetailRepo := f.InvestasiNonTbkDetailRepository
	investasiTbkRepo := f.InvestasiTbkRepository
	investasiTbkDetailRepo := f.InvestasiTbkDetailRepository
	mutasiDtaRepo := f.MutasiDtaRepository
	mutasiDtaDetailRepo := f.MutasiDtaDetailRepository
	mutasiFaRepo := f.MutasiFaRepository
	mutasiFaDetailRepo := f.MutasiFaDetailRepository
	mutasiIaRepo := f.MutasiIaRepository
	mutasiIaDetailRepo := f.MutasiIaDetailRepository
	mutasiPersediaanRepo := f.MutasiPersediaanRepository
	mutasiPersediaanDetailRepo := f.MutasiPersediaanDetailRepository
	mutasiRuaRepo := f.MutasiRuaRepository
	mutasiRuaDetailRepo := f.MutasiRuaDetailRepository
	pembelianPenjualanBerelasiRepo := f.PembelianPenjualanBerelasiRepository
	pembelianPenjualanBerelasiDetailRepo := f.PembelianPenjualanBerelasiDetailRepository
	trialBalanceRepo := f.TrialBalanceRepository
	trialBalanceDetailRepo := f.TrialBalanceDetailRepository
	formatterDetailRepo := f.FormatterDetailRepository
	companyRepo := f.CompanyRepository
	coaRepo := f.CoaRepository

	return &service{
		Repository:                                 repository,
		Db:                                         db,
		NotificationRepository:                     notifRepo,
		ImportedWorksheetRepository:                importedWorksheetRepo,
		ConsolidationRepository:                    consolidationRepo,
		AgingUtangPiutangRepository:                agingUtangPiutangRepo,
		AgingUtangPiutangDetailRepository:          agingUtangPiutangDetailRepo,
		FormatterRepository:                        formatterRepo,
		FormatterBridgesRepository:                 formatterBridgesRepo,
		EmployeeBenefitRepository:                  employeeBenefitRepo,
		EmployeeBenefitDetailRepository:            employeeBenefitDetailRepo,
		AdjustmentRepository:                       adjustmentRepo,
		AdjustmentDetailRepository:                 adjustmentDetailRepo,
		InvestasiNonTbkRepository:                  investasiNonTbkRepo,
		InvestasiNonTbkDetailRepository:            investasiNonTbkDetailRepo,
		InvestasiTbkRepository:                     investasiTbkRepo,
		InvestasiTbkDetailRepository:               investasiTbkDetailRepo,
		MutasiDtaRepository:                        mutasiDtaRepo,
		MutasiDtaDetailRepository:                  mutasiDtaDetailRepo,
		MutasiFaRepository:                         mutasiFaRepo,
		MutasiFaDetailRepository:                   mutasiFaDetailRepo,
		MutasiIaRepository:                         mutasiIaRepo,
		MutasiIaDetailRepository:                   mutasiIaDetailRepo,
		MutasiPersediaanRepository:                 mutasiPersediaanRepo,
		MutasiPersediaanDetailRepository:           mutasiPersediaanDetailRepo,
		MutasiRuaRepository:                        mutasiRuaRepo,
		MutasiRuaDetailRepository:                  mutasiRuaDetailRepo,
		PembelianPenjualanBerelasiRepository:       pembelianPenjualanBerelasiRepo,
		PembelianPenjualanBerelasiDetailRepository: pembelianPenjualanBerelasiDetailRepo,
		TrialBalanceRepository:                     trialBalanceRepo,
		TrialBalanceDetailRepository:               trialBalanceDetailRepo,
		FormatterDetailRepository:                  formatterDetailRepo,
		CompanyRepository:                          companyRepo,
		CoaRepository:                              coaRepo}
}
func (s *service) Update(ctx *abstraction.Context, id int) (*model.ImportedWorksheetDetailEntityModel, error) {
	var data model.ImportedWorksheetDetailEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		byIDDetail, err := s.Repository.FindByIDWorksheetDetail(ctx, &id)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err)
		}
		data.Context = ctx
		data.ImportedWorksheetDetailEntity = byIDDetail.ImportedWorksheetDetailEntity
		data.ImportedWorksheetDetailEntity.Status = 0
		result, err := s.Repository.Update(ctx, &id, &data)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
		}
		data = *result
		return nil
	}); err != nil {
		return &model.ImportedWorksheetDetailEntityModel{}, err
	}

	return &data, nil
}
func (s *service) FindByIDWorksheetDetail(ctx *abstraction.Context, id int) (*model.ImportedWorksheetDetailEntityModel, error) {
	data, err := s.Repository.FindByIDWorksheetDetail(ctx, &id)
	if err != nil {
		return &model.ImportedWorksheetDetailEntityModel{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := data
	return result, nil
}
func (s *service) Create(ctx *abstraction.Context, payload *dto.ImportedWorksheetCreateRequest) (*dto.ImportedWorksheetCreateResponse, error) {

	var data model.ImportedWorksheetEntityModel

	if err = trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {

		data.Context = ctx
		var tb model.TrialBalanceEntityModel
		tb.CompanyID = payload.CompanyID
		tb.Period = payload.Period

		getVersions, err := s.Repository.FindCompany(ctx, &tb)
		if err != nil {
			return err
		}
		format := "2006-01-02"

		hminsatubulan, err := time.Parse(format, payload.Period)
		if err != nil {
			return nil
		}

		hminbulansaatini := int(hminsatubulan.Month())

		now := time.Now()
		bulanSaatIni := int(now.Month())

		if hminbulansaatini == bulanSaatIni {
			return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Memasukan Period Saat Ini", "Tidak Dapat Memasukan Period Saat Ini")
		}

		data.ImportedWorksheetEntity = model.ImportedWorksheetEntity{
			Versions:  len(*getVersions) + 1,
			CompanyID: payload.CompanyID,
			Period:    payload.Period,
			Status:    1,
		}
		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
		}
		findByID, err := s.Repository.FindByID(ctx, &result.ID)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
		}
		allowed := helper.CompanyValidation(ctx.Auth.ID, findByID.CompanyID)
		if !allowed {

			return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Akses Company", "Tidak Dapat Akses Company")
		}

		data = *result
		return nil
	}); err != nil {
		return &dto.ImportedWorksheetCreateResponse{}, err
	}
	result := &dto.ImportedWorksheetCreateResponse{
		ImportedWorksheetEntityModel: data,
	}
	return result, nil
}
func (s *service) FindByID(ctx *abstraction.Context, payload *dto.ImportReUploadRequest) (*dto.ImportReUploadResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.ImportReUploadResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := &dto.ImportReUploadResponse{
		ImportedWorksheetEntityModel: *data,
	}
	return result, nil
}

func (s *service) Download(ctx *abstraction.Context, payload *dto.ImportedWorksheetDetailGetByIDRequest) (*dto.ImportedWorksheetDetailGetByIDResponse, error) {
	data, err := s.Repository.Download(ctx, &payload.ID)
	if err != nil {
		return &dto.ImportedWorksheetDetailGetByIDResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := &dto.ImportedWorksheetDetailGetByIDResponse{
		ImportedWorksheetDetailEntityModel: *data,
	}
	return result, nil
}
func (s *service) DownloadAll(ctx *abstraction.Context, payload *dto.ImportedWorksheetGetByIDRequest) (*dto.ImportedWorksheetGetByIDDownloadAllResponse, error) {

	data, err := s.Repository.DownloadAll(ctx, &payload.ID)
	if err != nil {
		return &dto.ImportedWorksheetGetByIDDownloadAllResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	var Filename []string

	for _, v := range *data {
		b := v.Note
		Filename = append(Filename, b)
	}

	result := &dto.ImportedWorksheetGetByIDDownloadAllResponse{
		Datas:    *data,
		FileName: Filename,
	}
	return result, nil
}
func (s *service) FindByFormatterBridges(ctx *abstraction.Context, fbi *int, source *string) (*model.FormatterBridgesEntityModel, error) {
	data, err := s.Repository.FindByFormatterBridges(ctx, fbi, source)
	if err != nil {
		return &model.FormatterBridgesEntityModel{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	return data, nil
}

func (s *service) DeletedAJE(ctx *abstraction.Context, companyId *int, tbId *int) (*model.AdjustmentEntityModel, error) {
	data, err := s.Repository.DeleteAje(ctx, companyId, tbId)
	if err != nil {
		return &model.AdjustmentEntityModel{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	return data, nil
}
func (s *service) DeletedJelim(ctx *abstraction.Context, companyId *int, tbId *int) (*model.JelimEntityModel, error) {
	data, err := s.Repository.DeleteJelim(ctx, companyId, tbId)
	if err != nil {
		return &model.JelimEntityModel{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	return data, nil
}
func (s *service) DeletedJcte(ctx *abstraction.Context, companyId *int, tbId *int) (*model.JcteEntityModel, error) {
	data, err := s.Repository.DeleteJcte(ctx, companyId, tbId)
	if err != nil {
		return &model.JcteEntityModel{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	return data, nil
}
func (s *service) DeletedJpm(ctx *abstraction.Context, companyId *int, tbId *int) (*model.JpmEntityModel, error) {
	data, err := s.Repository.DeleteJpm(ctx, companyId, tbId)
	if err != nil {
		return &model.JpmEntityModel{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	return data, nil
}

func (s *service) DeleteTBDF(ctx *abstraction.Context, fbi *int) (*model.TrialBalanceDetailEntityModel, error) {
	data, err := s.Repository.DeleteTBDF(ctx, fbi)
	if err != nil {
		return &model.TrialBalanceDetailEntityModel{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	return data, nil
}

func (s *service) FindByVCTrialBalance(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*model.TrialBalanceEntityModel, error) {
	data, err := s.Repository.FindByVCTrialBalance(ctx, m)
	if err != nil {
		return &model.TrialBalanceEntityModel{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	return data, nil
}

func (s *service) FindByVCAgingUtangPiutang(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel) (*model.AgingUtangPiutangEntityModel, error) {
	data, err := s.Repository.FindByVCAgingUtangPiutang(ctx, m)
	if err != nil {
		return &model.AgingUtangPiutangEntityModel{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	return data, nil
}

func (s *service) FindByVCMutasiDta(ctx *abstraction.Context, m *model.MutasiDtaFilterModel) (*model.MutasiDtaEntityModel, error) {
	data, err := s.Repository.FindByVCMutasiDta(ctx, m)
	if err != nil {
		return &model.MutasiDtaEntityModel{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	return data, nil
}

func (s *service) FindByVCMutasiRua(ctx *abstraction.Context, m *model.MutasiRuaFilterModel) (*model.MutasiRuaEntityModel, error) {
	data, err := s.Repository.FindByVCMutasiRua(ctx, m)
	if err != nil {
		return &model.MutasiRuaEntityModel{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	return data, nil
}

func (s *service) FindByVCMutasiIa(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*model.MutasiIaEntityModel, error) {
	data, err := s.Repository.FindByVCMutasiIa(ctx, m)
	if err != nil {
		return &model.MutasiIaEntityModel{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	return data, nil
}

func (s *service) FindByVCMutasiPersediaan(ctx *abstraction.Context, m *model.MutasiPersediaanFilterModel) (*model.MutasiPersediaanEntityModel, error) {
	data, err := s.Repository.FindByVCMutasiPersediaan(ctx, m)
	if err != nil {
		return &model.MutasiPersediaanEntityModel{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	return data, nil
}

func (s *service) FindByVCInvestasiTbk(ctx *abstraction.Context, m *model.InvestasiTbkFilterModel) (*model.InvestasiTbkEntityModel, error) {
	data, err := s.Repository.FindByVCInvestasiTbk(ctx, m)
	if err != nil {
		return &model.InvestasiTbkEntityModel{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	return data, nil
}

func (s *service) FindByVCInvestasiNonTbk(ctx *abstraction.Context, m *model.InvestasiNonTbkFilterModel) (*model.InvestasiNonTbkEntityModel, error) {
	data, err := s.Repository.FindByVCInvestasiNonTbk(ctx, m)
	if err != nil {
		return &model.InvestasiNonTbkEntityModel{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	return data, nil
}

func (s *service) FindByVCMutasiFa(ctx *abstraction.Context, m *model.MutasiFaFilterModel) (*model.MutasiFaEntityModel, error) {
	data, err := s.Repository.FindByVCMutasiFa(ctx, m)
	if err != nil {
		return &model.MutasiFaEntityModel{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	return data, nil
}

func (s *service) FindByVCPembelianPenjualanBerelasi(ctx *abstraction.Context, m *model.PembelianPenjualanBerelasiFilterModel) (*model.PembelianPenjualanBerelasiEntityModel, error) {
	data, err := s.Repository.FindByVCPembelianPenjualanBerelasi(ctx, m)
	if err != nil {
		return &model.PembelianPenjualanBerelasiEntityModel{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	return data, nil
}

func (s *service) FindByVCEmployeeBenefit(ctx *abstraction.Context, m *model.EmployeeBenefitFilterModel) (*model.EmployeeBenefitEntityModel, error) {
	data, err := s.Repository.FindByVCEmployeeBenefit(ctx, m)
	if err != nil {
		return &model.EmployeeBenefitEntityModel{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	return data, nil
}

func (s *service) FindByIDDetail(ctx *abstraction.Context, id *int) (*dto.ImportReUploadDetailResponse, error) {
	data, err := s.Repository.FindByIDDetail(ctx, id)
	if err != nil {
		return &dto.ImportReUploadDetailResponse{}, err
	}
	result := &dto.ImportReUploadDetailResponse{
		Data: *data,
	}
	return result, nil
}

func (s *service) FindByIDDetailS(ctx *abstraction.Context, id *int) (*dto.ImportReUploadDetailResponse, error) {
	data, err := s.Repository.FindByIDDetailS(ctx, id)
	if err != nil {
		return &dto.ImportReUploadDetailResponse{}, err
	}
	result := &dto.ImportReUploadDetailResponse{
		Data: *data,
	}
	return result, nil
}

func (s *service) ExportAllTemplate(ctx *abstraction.Context) (*string, error) {

	f := excelize.NewFile()
	currentSheet := f.GetSheetName(f.GetActiveSheetIndex())
	f.DeleteSheet(currentSheet)

	errorMessages := []string{}

	a, err := s.ExportTrialBalance(ctx, f)
	if err != nil {
		return nil, err
	}
	f = a

	b, err := s.ExportAgingUtangPiutang(ctx, f)
	if err != nil {
		return nil, err
	}
	f = b

	h, err := s.ExportMutasiFa(ctx, f)
	if err != nil {
		return nil, err
	}
	f = h

	i, err := s.ExportMutasiIa(ctx, f)
	if err != nil {
		return nil, err
	}
	f = i

	k, err := s.ExportMutasiRua(ctx, f)
	if err != nil {
		return nil, err
	}
	f = k

	g, err := s.ExportMutasiDta(ctx, f)
	if err != nil {
		return nil, err
	}
	f = g

	j, err := s.ExportMutasiPersediaan(ctx, f)
	if err != nil {
		return nil, err
	}
	f = j

	d, err := s.ExportInvestasiNonTbk(ctx, f)
	if err != nil {
		return nil, err
	}
	f = d

	e, err := s.ExportInvestasiTbk(ctx, f)
	if err != nil {
		return nil, err
	}
	f = e

	c, err := s.ExportEmployeeBenefit(ctx, f)
	if err != nil {
		return nil, err
	}
	f = c

	l, err := s.ExportPembelianPenjualanBerelasi(ctx, f)
	if err != nil {
		return nil, err
	}
	f = l
	if len(errorMessages) > 0 {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New(strings.Join(errorMessages, ", ")))
	}

	tmpFolder := fmt.Sprintf("assets/%s", "template")
	_, err = os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	fileName := "all-worksheet.xlsx"
	fileLoc := fmt.Sprintf("assets/%s/%s", "template", fileName)
	err = f.SaveAs(fileLoc)
	if err != nil {
		return nil, err
	}
	f, err = excelize.OpenFile(fileLoc)
	if err != nil {
		return nil, err
	}

	rows, err := f.GetRows("TRIAL_BALANCE")
	if err != nil {
		fmt.Println("Error reading rows:", err)
		return nil, err
	}
	// startRow := 9
	// endRow := 15
	
	t := true
	formatterID := 3
	var criteriaFormatterGrouping model.FormatterDetailFilterModel
	criteriaFormatterGrouping.FormatterID = &formatterID
	criteriaFormatterGrouping.IsLabel = &t
	// criteriaFormatterGrouping.IsTotal = &fa
	pagesize := 100000
	tmpStr := "sort_id"
	tmpStr1 := "ASC"
	paginationTB := abstraction.Pagination{
		PageSize: &pagesize,
		SortBy:   &tmpStr,
		Sort:     &tmpStr1,
	}
	dataGrouping, _, err := s.FormatterDetailRepository.FindGroup(ctx, &criteriaFormatterGrouping, &paginationTB)
	if err != nil {
		return nil, err
	}
	

	lineCode := make(map[string]int)
	for i, row := range rows {
		
		if len(row) < 2 {
			cellValue := "test"
			lineCode[cellValue] = i
		}else {
			cellValue := row[2]
			lineCode[cellValue] = i
		}
		i++
		// Simpan nomor baris dalam map dengan kunci (key) nilai sel
	}
	for i, d := range *dataGrouping {

		if *d.IsTotal == t {
			continue
		}
		if _,ok := lineCode[d.Description]; ok{
		
			startRow := lineCode[d.Description]
			startRow = startRow + 2

			if i+1 < len(*dataGrouping) {
				secondElement := (*dataGrouping)[i+1]
				endRow := lineCode[secondElement.Description]
				endRow = endRow - 1

				for row := startRow; row <= endRow; row++ {
					if err := f.SetRowOutlineLevel("TRIAL_BALANCE", row, 1); err != nil {
						return nil, err
					}
				}
			}
		}
	}
	tmpFolder = fmt.Sprintf("assets/%s", "template")
	_, err = os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	fileName = "all-worksheet.xlsx"
	fileLoc = fmt.Sprintf("assets/%s/%s", "template", fileName)
	err = f.SaveAs(fileLoc)
	if err != nil {
		return nil, err
	}
	return &fileLoc, nil
}

func (s *service) ExportAdjustment(ctx *abstraction.Context, f *excelize.File) (*excelize.File, error) {

	sheet := "ADJUSTMENT"
	f.NewSheet(sheet)
	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 6.50},
		{"COL": "B", "WIDTH": 15.38},
		{"COL": "C", "WIDTH": 7.25},
		{"COL": "D", "WIDTH": 2.14},
		{"COL": "E", "WIDTH": 57.45},
		{"COL": "F", "WIDTH": 7.25},
		{"COL": "G", "WIDTH": 15.38},
		{"COL": "H", "WIDTH": 15.38},
		{"COL": "I", "WIDTH": 15.38},
		{"COL": "J", "WIDTH": 15.38},
	}

	for _, v := range arrStyleColWidth {
		tmpColWidth := fmt.Sprintf("%f", v["WIDTH"])
		colWidth, _ := strconv.ParseFloat(tmpColWidth, 64)
		f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
	}

	styleDefault, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		return nil, err
	}

	f.SetColStyle(sheet, "A:Z", styleDefault)
	stylingBorderLROnly, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
		},
	})
	if err != nil {
		return nil, err
	}

	stylingBorderAll, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
		},
	})
	if err != nil {
		return nil, err
	}

	stylingHeader, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#fac090"},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}

	f.MergeCell(sheet, "A6", "A7")
	f.MergeCell(sheet, "B6", "B7")
	f.MergeCell(sheet, "C6", "C7")
	f.MergeCell(sheet, "D6", "E7")
	f.MergeCell(sheet, "F6", "F7")
	f.MergeCell(sheet, "G6", "H6")
	f.MergeCell(sheet, "I6", "j6")

	stylingCurrency, err := f.NewStyle(&excelize.Style{
		NumFmt: 7,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	stylingSubTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
		},

		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#00ff00"},
		},

		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return nil, err
	}

	stylingSubTotal2, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
		},

		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#99ff33"},
		},

		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return nil, err
	}

	f.SetCellValue(sheet, "A2", "Company")
	f.SetCellValue(sheet, "B2", ": "+"tb.Company.CompanyEntity.Name")
	f.SetCellValue(sheet, "A3", "Date")
	f.SetCellValue(sheet, "B3", ": "+"")
	f.SetCellValue(sheet, "A4", "Subject")
	f.SetCellValue(sheet, "B4", ": AJE")

	f.SetCellStyle(sheet, "A6", "J7", stylingHeader)
	f.SetCellValue(sheet, "A6", "NO")
	f.SetCellValue(sheet, "B6", "COA")
	f.SetCellValue(sheet, "C6", "AJE")
	f.SetCellValue(sheet, "D6", "DESCRIPTION")
	f.SetCellValue(sheet, "F6", "WP")
	f.SetCellValue(sheet, "F7", "REFF")
	f.SetCellValue(sheet, "G6", "Balance Sheet")
	f.SetCellValue(sheet, "G7", "DR")
	f.SetCellValue(sheet, "H7", "CR")
	f.SetCellValue(sheet, "I6", "Income Stat")
	f.SetCellValue(sheet, "I7", "DR")
	f.SetCellValue(sheet, "J7", "CR")

	row := 8
	// counterRef := 1
	rowBefore := row

	f.SetCellStyle(sheet, "A8", fmt.Sprintf("B%d", row), stylingBorderLROnly)
	f.SetCellStyle(sheet, "C8", fmt.Sprintf("C%d", row), stylingBorderLROnly)
	f.SetCellStyle(sheet, "F8", fmt.Sprintf("F%d", row), stylingBorderLROnly)
	f.SetCellStyle(sheet, "G8", fmt.Sprintf("J%d", row), stylingCurrency)
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("F%d", row+1), stylingBorderAll)
	f.MergeCell(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("E%d", row))
	f.MergeCell(sheet, fmt.Sprintf("D%d", row+1), fmt.Sprintf("E%d", row+1))

	f.SetSheetFormatPr(sheet, excelize.DefaultRowHeight(12.85))

	f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "Total PM")
	f.SetCellFormula(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("=SUM(G%d:G%d)", rowBefore, row-1))
	f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("=SUM(H%d:H%d)", rowBefore, row-1))
	f.SetCellFormula(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("=SUM(I%d:I%d)", rowBefore, row-1))
	f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=SUM(J%d:J%d)", rowBefore, row-1))
	f.SetCellFormula(sheet, fmt.Sprintf("I%d", row+1), fmt.Sprintf("=G%d+I%d", row, row))
	f.SetCellFormula(sheet, fmt.Sprintf("J%d", row+1), fmt.Sprintf("=H%d+J%d", row, row))
	f.SetCellFormula(sheet, fmt.Sprintf("J%d", row+2), fmt.Sprintf("=I%d-J%d", row+1, row+1))

	f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), stylingSubTotal)
	f.SetCellStyle(sheet, fmt.Sprintf("G%d", row+1), fmt.Sprintf("J%d", row+1), stylingSubTotal2)

	f.SetDefaultFont("Arial")

	return f, nil
}
func (s *service) ExportMutasiFa(ctx *abstraction.Context, f *excelize.File) (*excelize.File, error) {
	sheet := "MUTASI_FA"
	f.NewSheet(sheet)

	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 8.21},
		{"COL": "B", "WIDTH": 3.74},
		{"COL": "C", "WIDTH": 33.57},
		{"COL": "D", "WIDTH": 18.21},
		{"COL": "E", "WIDTH": 16.60},
		{"COL": "F", "WIDTH": 16.60},
		{"COL": "G", "WIDTH": 19.64},
		{"COL": "H", "WIDTH": 19.64},
		{"COL": "I", "WIDTH": 18.21},
		{"COL": "J", "WIDTH": 17.68},
		{"COL": "K", "WIDTH": 16.78},
	}
	for _, v := range arrStyleColWidth {
		tmpColWidth := fmt.Sprintf("%f", v["WIDTH"])
		colWidth, _ := strconv.ParseFloat(tmpColWidth, 64)
		err := f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			return nil, err
		}
	}
	stylingControl, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#FFFF00"},
		},
	})
	if err != nil {
		return nil, err
	}
	styleHeader, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#66ff33"},
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}
	numberFormat := "#,##"
	styleCurrency, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}

	styleCurrencyTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#66ff33"},
		},
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}

	styleLabelTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#66ff33"},
		},
	})
	if err != nil {
		return nil, err
	}

	styleLabel, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	defaultStyle, err := f.NewStyle(&excelize.Style{
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}

	styleCurrencyWoBorder, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}

	formatterCode := []string{"MUTASI-FA-COST", "MUTASI-FA-ACCUMULATED-DEPRECATION"}
	formatterTitle := []string{"Mutasi Fixed Assets (FA)", ""}
	row, rowStart := 7, 7

	f.SetCellValue(sheet, "B2", formatterTitle[0])

	f.MergeCell(sheet, "B4", "C6")
	f.MergeCell(sheet, "D4", "K4")
	f.MergeCell(sheet, "D5", "D6")
	f.MergeCell(sheet, "E5", "E6")
	f.MergeCell(sheet, "F5", "F6")
	f.MergeCell(sheet, "G5", "G6")
	f.MergeCell(sheet, "H5", "H6")
	f.MergeCell(sheet, "I5", "I6")
	f.MergeCell(sheet, "J5", "J6")
	f.MergeCell(sheet, "K5", "K6")

	f.SetCellStyle(sheet, "B4", "K6", styleHeader)
	f.SetCellFormula(sheet, "B4" , "=TRIAL_BALANCE!D2")
	f.SetCellValue(sheet, "D4", "")
	f.SetCellValue(sheet, "D5", "Beginning Balance")
	f.SetCellValue(sheet, "E5", "Acquisition of Subsidiary")
	f.SetCellValue(sheet, "F5", "Additions (+)")
	f.SetCellValue(sheet, "G5", "Deductions (-)")
	f.SetCellValue(sheet, "H5", "Reclassification")
	f.SetCellValue(sheet, "I5", "Revaluation")
	f.SetCellValue(sheet, "J5", "Ending balance")
	f.SetCellValue(sheet, "K5", "Control")
	rowCode := make(map[string]int)
	for _, formatter := range formatterCode {

		var criteria dto.FormatterGetRequest
		criteria.FormatterFilterModel.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
		if err != nil {
			return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
		}

		partRowStart := row
		for _, v := range data.FormatterDetail {
			valuekosong := 0.0
			rowCode[v.Code] = row
			f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabel)
			// f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("K%d", row), styleCurrency)

			if strings.Contains(strings.ToLower(v.Code), "blank") {
				row++
				continue
			}

			if strings.ToUpper(v.Code) == "ACCUMULATED_DEPRECIATION" {
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
				row++
				continue
			}
			if strings.ToUpper(v.Code) == "COST:" {
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
				row++
				continue
			}
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), valuekosong)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), valuekosong)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), valuekosong)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), valuekosong)
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), valuekosong)
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), valuekosong)
			f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("I%d", row), styleCurrency)

			if v.ControlFormula != "" {

				if v.AutoSummary != nil && *v.AutoSummary {
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
					f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), styleCurrencyTotal)
					f.MergeCell(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row))
					if strings.ToUpper(v.Code) == "NET_BOOK_VALUE" {
						f.SetCellFormula(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("=SUM(D%d:D%d)", partRowStart, row-1))
						f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=SUM(J%d:J%d)", partRowStart, row-1))
					} else {
						for chr := 'D'; chr <= 'J'; chr++ {
							f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
						}
					}

				}
				f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)

				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)

				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					//cari jml berdasarkan code
					if _, ok := rowCode[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("J%d", rowCode[vMatch]))
					}
					if _, ok := tbRowCode[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCode[vMatch]))
						f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCode[vMatch]), "control")
					}

				}
				f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("=%s", formula))
				f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)

			}
			if v.IsTotal != nil && *v.IsTotal {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("K%d", row), styleCurrencyTotal)
				if v.FxSummary == "" {
					row++
					continue
				}
				arrChr := []string{"D", "E", "F", "G", "H", "I", "J", "K"}
				if strings.ToUpper(v.Code) == "NET_BOOK_VALUE" {
					arrChr = []string{"D", "J"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_1" {
					arrChr = []string{"F"}
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("K%d", row), styleCurrency)
					f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingControl)
				}
				if strings.ToUpper(v.Code) == "CONTROL_2" {
					arrChr = []string{"J"}
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("K%d", row), defaultStyle)
				}

				for _, chr := range arrChr {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						//cari jml berdasarkan code
						if rowCode[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCode[vMatch]))
						}
						if _, ok := tbRowCode[vMatch]; ok {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCode[vMatch]))
							f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCode[vMatch]), "control")
						}

					}
					if strings.ToUpper(v.Code) == "CONTROL_2" {
						row = row - 1
						f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)
						f.SetCellFormula(sheet, fmt.Sprintf("%s%d", "K", row), fmt.Sprintf("=%s", formula))
						continue

					}
					f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row), fmt.Sprintf("=%s", formula))
				}
				row++
				continue
			}

			if v.AutoSummary != nil && *v.AutoSummary {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), styleCurrencyTotal)
				f.MergeCell(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row))
				if strings.ToUpper(v.Code) == "NET_BOOK_VALUE" {
					f.SetCellFormula(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("=SUM(D%d:D%d)", partRowStart, row-1))
					f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=SUM(J%d:J%d)", partRowStart, row-1))
				} else {
					for chr := 'D'; chr <= 'J'; chr++ {
						f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
					}
				}
				row++
				partRowStart = row
				continue
			}
			f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=D%d+E%d+F%d-G%d+H%d+I%d", row, row, row, row, row, row))
			f.SetCellStyle(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("J%d", row), styleCurrencyWoBorder)
			f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)

			row++
		}
		rowStart = row
		row = rowStart
	}

	//Penambahan detail pengurangan
	row += 2
	var criteria dto.FormatterGetRequest
	tmpStr := "MUTASI-DETAIL-PENGURANGAN"
	criteria.FormatterFilterModel.FormatterFor = &tmpStr

	data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "Detail Pengurangan:")
	f.SetCellValue(sheet, fmt.Sprintf("E%d", row), "Penjualan")
	f.SetCellValue(sheet, fmt.Sprintf("F%d", row), "Penghapusan")
	row += 1
	partRowStart := row
	
	for _, v := range data.FormatterDetail {
		rowCode[v.Code] = row
		f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("F%d", row), styleCurrencyWoBorder)

		if strings.Contains(strings.ToLower(v.Code), "blank") {
			row++
			continue
		}

		rowKosong := 0.0
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), rowKosong)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), rowKosong)

		if v.Code == "CONTROL_1" || v.Code == "CONTROL_2" {
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), "")
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), "")
		}
		if v.IsTotal != nil && *v.IsTotal {
			if v.FxSummary == "" {
				row++
				continue
			}
			arrChr := []string{"E", "F"}
			

			// if strings.ToUpper(v.Code) == "CONTROL_2" {
			// 	arrChr = []string{"F"}
			// 	f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), styleCurrency)
			// 	f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), stylingControl)
			// }
			for _, chr := range arrChr {
				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					//cari jml berdasarkan code
					if rowCode[vMatch] != 0 {
						if v.IsCoa != nil && *v.IsCoa {
							if chr == "E" {
								formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("G%d", rowCode[vMatch]))
							} else {
								formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("I%d", rowCode[vMatch]))
							}
						} else {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCode[vMatch]))
						}
					}
					if _, ok := tbRowCode[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCode[vMatch]))
						f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCode[vMatch]), "control")
					}
				}
				if strings.ToUpper(v.Code) == "CONTROL_1" {
					// arrChr = []string{"E"}
					f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("F%d", row), styleCurrency)
					f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("F%d", row), stylingControl)
				}
				// if strings.ToUpper(v.Code) == "CONTROL_2" {
				// 	row = row - 1
				// 	f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), styleCurrency)
				// 	f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingControl)
				// }
				err = f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row), fmt.Sprintf("=%s", formula))
				if err != nil {
					fmt.Println(err)
				}
			}
			row++
			continue
		}

		if v.AutoSummary != nil && *v.AutoSummary {
			for chr := 'E'; chr <= 'F'; chr++ {
				f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
			}
			row++
			partRowStart = row
			continue
		}

		row++
	}

	return f, nil
}
func (s *service) ExportMutasiIa(ctx *abstraction.Context, f *excelize.File) (*excelize.File, error) {
	sheet := "MUTASI_IA"
	f.NewSheet(sheet)

	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 8.21},
		{"COL": "B", "WIDTH": 3.74},
		{"COL": "C", "WIDTH": 33.57},
		{"COL": "D", "WIDTH": 18.21},
		{"COL": "E", "WIDTH": 16.60},
		{"COL": "F", "WIDTH": 16.60},
		{"COL": "G", "WIDTH": 19.64},
		{"COL": "H", "WIDTH": 19.64},
		{"COL": "I", "WIDTH": 18.21},
		{"COL": "J", "WIDTH": 17.68},
		{"COL": "K", "WIDTH": 16.78},
	}
	for _, v := range arrStyleColWidth {
		tmpColWidth := fmt.Sprintf("%f", v["WIDTH"])
		colWidth, _ := strconv.ParseFloat(tmpColWidth, 64)
		err := f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			return nil, err
		}
	}
	stylingControl, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#FFFF00"},
		},
	})
	if err != nil {
		return nil, err
	}
	styleHeader, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#66ff33"},
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}
	numberFormat := "#,##"
	styleCurrency, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}

	styleCurrencyTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#66ff33"},
		},
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}

	styleLabelTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#66ff33"},
		},
	})
	if err != nil {
		return nil, err
	}

	styleLabel, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	styleCurrencyWoBorder, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}
	defaultStyle, err := f.NewStyle(&excelize.Style{
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}
	formatterCode := []string{"MUTASI-IA-COST", "MUTASI-IA-ACCUMULATED-DEPRECATION"}
	formatterTitle := []string{"Mutasi Intangible Assets (IA)", ""}
	row, rowStart := 7, 7

	f.SetCellValue(sheet, "B2", formatterTitle[0])

	f.MergeCell(sheet, "B4", "C6")
	f.MergeCell(sheet, "D4", "K4")
	f.MergeCell(sheet, "D5", "D6")
	f.MergeCell(sheet, "E5", "E6")
	f.MergeCell(sheet, "F5", "F6")
	f.MergeCell(sheet, "G5", "G6")
	f.MergeCell(sheet, "H5", "H6")
	f.MergeCell(sheet, "I5", "I6")
	f.MergeCell(sheet, "J5", "J6")
	f.MergeCell(sheet, "K5", "K6")

	f.SetCellStyle(sheet, "B4", "K6", styleHeader)
	f.SetCellFormula(sheet, "B4" , "=TRIAL_BALANCE!D2")
	f.SetCellValue(sheet, "D4", "")
	f.SetCellValue(sheet, "D5", "Beginning Balance")
	f.SetCellValue(sheet, "E5", "Acquisition of Subsidiary")
	f.SetCellValue(sheet, "F5", "Additions (+)")
	f.SetCellValue(sheet, "G5", "Deductions (-)")
	f.SetCellValue(sheet, "H5", "Reclassification")
	f.SetCellValue(sheet, "I5", "Revaluation")
	f.SetCellValue(sheet, "J5", "Ending balance")
	f.SetCellValue(sheet, "K5", "Control")
	rowCode := make(map[string]int)
	for _, formatter := range formatterCode {

		var criteria dto.FormatterGetRequest
		criteria.FormatterFilterModel.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
		if err != nil {
			return nil, helper.ErrorHandler(err)
		}

		partRowStart := row
		for _, v := range data.FormatterDetail {
			valuekosong := 0.0
			rowCode[v.Code] = row
			f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabel)
			// f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("K%d", row), styleCurrency)

			if strings.Contains(strings.ToLower(v.Code), "blank") {
				row++
				continue
			}
			if strings.ToUpper(v.Code) == "ACCUMULATED_DEPRECIATION" {
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
				row++
				continue
			}
			if strings.ToUpper(v.Code) == "COST:" {
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
				row++
				continue
			}

			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), valuekosong)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), valuekosong)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), valuekosong)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), valuekosong)
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), valuekosong)
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), valuekosong)
			f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("I%d", row), styleCurrency)

			if v.ControlFormula != "" {

				if v.AutoSummary != nil && *v.AutoSummary {
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
					f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), styleCurrencyTotal)
					f.MergeCell(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row))
					if strings.ToUpper(v.Code) == "NET_BOOK_VALUE" {
						f.SetCellFormula(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("=SUM(D%d:D%d)", partRowStart, row-1))
						f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=SUM(J%d:J%d)", partRowStart, row-1))
					} else {
						for chr := 'D'; chr <= 'J'; chr++ {
							f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
						}
					}

				}
				f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)

				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)

				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					//cari jml berdasarkan code
					if _, ok := rowCode[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("J%d", rowCode[vMatch]))
					}
					if _, ok := tbRowCode[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCode[vMatch]))
						f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCode[vMatch]), "control")
					}

				}
				f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("=%s", formula))
				f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)

			}
			if v.IsTotal != nil && *v.IsTotal {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("K%d", row), styleCurrencyTotal)
				if v.FxSummary == "" {
					row++
					continue
				}
				arrChr := []string{"D", "E", "F", "G", "H", "I", "J", "K"}
				if strings.ToUpper(v.Code) == "NET_BOOK_VALUE" {
					arrChr = []string{"D", "J"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_1" {
					arrChr = []string{"F"}
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("K%d", row), styleCurrency)
					f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingControl)
				}
				if strings.ToUpper(v.Code) == "CONTROL_2" {
					arrChr = []string{"J"}
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("K%d", row), defaultStyle)

				}

				for _, chr := range arrChr {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						//cari jml berdasarkan code
						if rowCode[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCode[vMatch]))
						}
						if _, ok := tbRowCode[vMatch]; ok {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCode[vMatch]))
							f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCode[vMatch]), "control")
						}

					}
					if strings.ToUpper(v.Code) == "CONTROL_2" {
						row = row - 1
						f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)
						f.SetCellFormula(sheet, fmt.Sprintf("%s%d", "K", row), fmt.Sprintf("=%s", formula))
						continue

					}
					f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row), fmt.Sprintf("=%s", formula))
				}
				row++
				continue
			}

			if v.AutoSummary != nil && *v.AutoSummary {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), styleCurrencyTotal)
				f.MergeCell(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row))
				if strings.ToUpper(v.Code) == "NET_BOOK_VALUE" {
					f.SetCellFormula(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("=SUM(D%d:D%d)", partRowStart, row-1))
					f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=SUM(J%d:J%d)", partRowStart, row-1))
				} else {
					for chr := 'D'; chr <= 'J'; chr++ {
						f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
					}
				}
				row++
				partRowStart = row
				continue
			}
			f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=D%d+E%d+F%d-G%d+H%d+I%d", row, row, row, row, row, row))
			f.SetCellStyle(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("J%d", row), styleCurrencyWoBorder)
			f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)
			// }
			row++
		}
		rowStart = row
		row = rowStart
	}

	//Penambahan detail pengurangan
	row += 2
	var criteria dto.FormatterGetRequest
	tmpStr := "MUTASI-DETAIL-PENGURANGAN"
	criteria.FormatterFilterModel.FormatterFor = &tmpStr

	data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
	if err != nil {
		return nil, err
	}
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "Detail Pengurangan:")
	f.SetCellValue(sheet, fmt.Sprintf("E%d", row), "Penjualan")
	f.SetCellValue(sheet, fmt.Sprintf("F%d", row), "Penghapusan")
	row += 1
	partRowStart := row
	
	for _, v := range data.FormatterDetail {
		rowCode[v.Code] = row
		f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("F%d", row), styleCurrencyWoBorder)

		if strings.Contains(strings.ToLower(v.Code), "blank") {
			row++
			continue
		}

		rowKosong := 0.0
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), rowKosong)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), rowKosong)
		if v.Code == "CONTROL_1" || v.Code == "CONTROL_2" {
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), "")
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), "")
		}
		if v.IsTotal != nil && *v.IsTotal {
			if v.FxSummary == "" {
				row++
				continue
			}
			if v.Code == "CONTROL_1" || v.Code == "CONTROL_2" {
				continue
			}
			arrChr := []string{"E", "F"}
			for _, chr := range arrChr {
				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z_~]+`)
				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					//cari jml berdasarkan code
					if rowCode[vMatch] != 0 {
						if v.IsCoa != nil && *v.IsCoa {
							if chr == "E" {
								formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("G%d", rowCode[vMatch]))
							} else {
								formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("I%d", rowCode[vMatch]))
							}
						} else {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCode[vMatch]))
						}
					}

				}
				err = f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row), fmt.Sprintf("=%s", formula))
				if err != nil {
					fmt.Println(err)
				}
			}
			row++
			continue
		}

		if v.AutoSummary != nil && *v.AutoSummary {
			for chr := 'E'; chr <= 'F'; chr++ {
				f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
			}
			row++
			partRowStart = row
			continue
		}

		if v.ControlFormula != "" {
			// f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
			// f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styleCurrencyTotal)
			if v.FxSummary == "" {
				row++
				continue
			}
			f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)
			formula := v.FxSummary
			reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)

			match := reg.FindAllString(formula, -1)
			for _, vMatch := range match {
				//cari jml berdasarkan code
				if _, ok := rowCode[vMatch]; ok {
					formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("J%d", rowCode[vMatch]))
				}
				if _, ok := tbRowCode[vMatch]; ok {
					formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCode[vMatch]))
					f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCode[vMatch]), "control")
				}

			}
			f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("=%s", formula))
			row++
			continue
		}
		row++
	}

	return f, nil
}
func (s *service) ExportMutasiRua(ctx *abstraction.Context, f *excelize.File) (*excelize.File, error) {
	sheet := "MUTASI_RUA"
	f.NewSheet(sheet)

	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 8.21},
		{"COL": "B", "WIDTH": 3.74},
		{"COL": "C", "WIDTH": 33.57},
		{"COL": "D", "WIDTH": 18.21},
		{"COL": "E", "WIDTH": 16.60},
		{"COL": "F", "WIDTH": 16.60},
		{"COL": "G", "WIDTH": 19.64},
		{"COL": "H", "WIDTH": 19.64},
		{"COL": "I", "WIDTH": 18.21},
		{"COL": "J", "WIDTH": 17.68},
		{"COL": "K", "WIDTH": 16.78},
	}
	for _, v := range arrStyleColWidth {
		tmpColWidth := fmt.Sprintf("%f", v["WIDTH"])
		colWidth, _ := strconv.ParseFloat(tmpColWidth, 64)
		err := f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			return nil, err
		}
	}
	stylingControl, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#FFFF00"},
		},
	})
	if err != nil {
		return nil, err
	}
	styleHeader, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#66ff33"},
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}
	numberFormat := "#,##"
	styleCurrency, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}

	styleCurrencyTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#66ff33"},
		},
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}
	defaultStyle, err := f.NewStyle(&excelize.Style{
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}
	styleLabelTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#66ff33"},
		},
	})
	if err != nil {
		return nil, err
	}

	styleLabel, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	styleCurrencyWoBorder, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}

	formatterCode := []string{"MUTASI-RUA-COST", "MUTASI-RUA-ACCUMULATED-DEPRECATION"}
	formatterTitle := []string{"Mutasi Right of Used Assets (RUA)"}
	row, rowStart := 7, 7

	f.SetCellValue(sheet, "B2", formatterTitle[0])

	f.MergeCell(sheet, "B4", "C6")
	f.MergeCell(sheet, "D4", "K4")
	f.MergeCell(sheet, "D5", "D6")
	f.MergeCell(sheet, "E5", "E6")
	f.MergeCell(sheet, "F5", "F6")
	f.MergeCell(sheet, "G5", "G6")
	f.MergeCell(sheet, "H5", "H6")
	f.MergeCell(sheet, "I5", "I6")
	f.MergeCell(sheet, "J5", "J6")
	f.MergeCell(sheet, "K5", "K6")

	f.SetCellStyle(sheet, "B4", "K6", styleHeader)
	f.SetCellFormula(sheet, "B4" , "=TRIAL_BALANCE!D2")
	f.SetCellValue(sheet, "D4", "")
	f.SetCellValue(sheet, "D5", "Beginning Balance")
	f.SetCellValue(sheet, "E5", "Acquisition of Subsidiary")
	f.SetCellValue(sheet, "F5", "Additions (+)")
	f.SetCellValue(sheet, "G5", "Deductions (-)")
	f.SetCellValue(sheet, "H5", "Reclassification")
	f.SetCellValue(sheet, "I5", "Revaluation")
	f.SetCellValue(sheet, "J5", "Ending balance")
	f.SetCellValue(sheet, "K5", "Control")
	rowCode := make(map[string]int)
	for _, formatter := range formatterCode {

		var criteria dto.FormatterGetRequest
		criteria.FormatterFilterModel.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
		if err != nil {
			return nil, helper.ErrorHandler(err)
		}

		partRowStart := row
		for _, v := range data.FormatterDetail {
			rowCode[v.Code] = row
			valuekosong := 0.0
			f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabel)
			// f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("K%d", row), styleCurrency)

			if strings.Contains(strings.ToLower(v.Code), "blank") {
				row++
				continue
			}
			if strings.ToUpper(v.Code) == "ACCUMULATED_DEPRECIATION" {
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
				row++
				continue
			}
			if strings.ToUpper(v.Code) == "COST:" {
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
				row++
				continue
			}

			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), valuekosong)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), valuekosong)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), valuekosong)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), valuekosong)
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), valuekosong)
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), valuekosong)
			f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("I%d", row), styleCurrency)

			if v.ControlFormula != "" {

				if v.AutoSummary != nil && *v.AutoSummary {
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
					f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), styleCurrencyTotal)
					f.MergeCell(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row))
					if strings.ToUpper(v.Code) == "NET_BOOK_VALUE" {
						f.SetCellFormula(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("=SUM(D%d:D%d)", partRowStart, row-1))
						f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=SUM(J%d:J%d)", partRowStart, row-1))
					} else {
						for chr := 'D'; chr <= 'J'; chr++ {
							f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
						}
					}

				}
				f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)

				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)

				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					//cari jml berdasarkan code
					if _, ok := rowCode[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("J%d", rowCode[vMatch]))
					}
					if _, ok := tbRowCode[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCode[vMatch]))
						f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCode[vMatch]), "control")
					}

				}
				f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("=%s", formula))
				f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)

			}
			if v.IsTotal != nil && *v.IsTotal {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("K%d", row), styleCurrencyTotal)
				if v.FxSummary == "" {
					row++
					continue
				}
				arrChr := []string{"D", "E", "F", "G", "H", "I", "J", "K"}
				if strings.ToUpper(v.Code) == "NET_BOOK_VALUE" {
					arrChr = []string{"D", "J"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_1" {
					arrChr = []string{"F"}
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("K%d", row), styleCurrency)
					f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingControl)
				}
				if strings.ToUpper(v.Code) == "CONTROL_2" {
					arrChr = []string{"J"}
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("K%d", row), defaultStyle)
				}

				for _, chr := range arrChr {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						//cari jml berdasarkan code
						if rowCode[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCode[vMatch]))
						}
						if _, ok := tbRowCode[vMatch]; ok {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCode[vMatch]))
							f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCode[vMatch]), "control")
						}

					}
					if strings.ToUpper(v.Code) == "CONTROL_2" {
						row = row - 1
						f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)
						f.SetCellFormula(sheet, fmt.Sprintf("%s%d", "K", row), fmt.Sprintf("=%s", formula))
						continue

					}
					f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row), fmt.Sprintf("=%s", formula))
				}
				row++
				continue
			}

			if v.AutoSummary != nil && *v.AutoSummary {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), styleCurrencyTotal)
				f.MergeCell(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row))
				if strings.ToUpper(v.Code) == "NET_BOOK_VALUE" {
					f.SetCellFormula(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("=SUM(D%d:D%d)", partRowStart, row-1))
					f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=SUM(J%d:J%d)", partRowStart, row-1))
				} else {
					for chr := 'D'; chr <= 'J'; chr++ {
						f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
					}
				}
				row++
				partRowStart = row
				continue
			}
			f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=D%d+E%d+F%d-G%d+H%d+I%d", row, row, row, row, row, row))
			f.SetCellStyle(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("J%d", row), styleCurrencyWoBorder)
			f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)

			row++
		}
		rowStart = row
		row = rowStart
	}

	//Penambahan detail pengurangan
	row += 2
	var criteria dto.FormatterGetRequest
	tmpStr := "MUTASI-DETAIL-PENGURANGAN"
	criteria.FormatterFilterModel.FormatterFor = &tmpStr

	data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
	if err != nil {
		return nil, err
	}

	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "Detail Pengurangan:")
	f.SetCellValue(sheet, fmt.Sprintf("E%d", row), "Penjualan")
	f.SetCellValue(sheet, fmt.Sprintf("F%d", row), "Penghapusan")
	row += 1
	partRowStart := row
	
	for _, v := range data.FormatterDetail {
		rowCode[v.Code] = row
		f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("F%d", row), styleCurrencyWoBorder)

		if strings.Contains(strings.ToLower(v.Code), "blank") {
			row++
			continue
		}

		rowKosong := 0.0
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), rowKosong)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), rowKosong)
		if v.Code == "CONTROL_1" || v.Code == "CONTROL_2" {
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), "")
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), "")
		}

		if v.IsTotal != nil && *v.IsTotal {
			if v.FxSummary == "" {
				row++
				continue
			}
			if v.Code == "CONTROL_1" || v.Code == "CONTROL_2" {
				continue
			}
			arrChr := []string{"E", "F"}
			for _, chr := range arrChr {
				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z_~]+`)
				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					//cari jml berdasarkan code
					if rowCode[vMatch] != 0 {
						if v.IsCoa != nil && *v.IsCoa {
							if chr == "E" {
								formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("G%d", rowCode[vMatch]))
							} else {
								formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("I%d", rowCode[vMatch]))
							}
						} else {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCode[vMatch]))
						}
					}

				}
				err = f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row), fmt.Sprintf("=%s", formula))
				if err != nil {
					fmt.Println(err)
				}
			}
			row++
			continue
		}

		if v.AutoSummary != nil && *v.AutoSummary {
			for chr := 'E'; chr <= 'F'; chr++ {
				f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
			}
			row++
			partRowStart = row
			continue
		}

		row++
	}

	return f, nil
}

var tbRowCode = make(map[string]int)

func (s *service) ExportTrialBalance(ctx *abstraction.Context, f *excelize.File) (*excelize.File, error) {
	sheet := "TRIAL_BALANCE"
	indexSheet := f.NewSheet(sheet)
	f.SetActiveSheet(indexSheet)
	var (
		criteriaFormatter model.FormatterDetailFilterModel
	)
	formatterID := 3
	t := true
	criteriaFormatter.IsShowExport = &t
	criteriaFormatter.FormatterID = &formatterID
	pagesize := 100000
	tmpStr := "sort_id"
	tmpStr1 := "ASC"
	paginationTB := abstraction.Pagination{
		PageSize: &pagesize,
		SortBy:   &tmpStr,
		Sort:     &tmpStr1,
	}

	data, _, err := s.FormatterDetailRepository.Find(ctx, &criteriaFormatter, &paginationTB)
	if err != nil {
		return nil, err
	}
	f.SetCellValue(sheet, "B2", "Company")
	f.SetCellValue(sheet, "C2", ": ")
	f.SetCellValue(sheet, "B3", "Date")
	f.SetCellValue(sheet, "C3", ": ")
	f.SetCellValue(sheet, "B4", "Subject")
	f.SetCellValue(sheet, "C4", ": DETAIL ASET, LIABILITAS & EKUITAS")

	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 0.83},
		{"COL": "B", "WIDTH": 15.38},
		{"COL": "C", "WIDTH": 2.14},
		{"COL": "D", "WIDTH": 2.14},
		{"COL": "E", "WIDTH": 57.45},
		{"COL": "F", "WIDTH": 6.43},
		{"COL": "G", "WIDTH": 17.65},
		{"COL": "H", "WIDTH": 10.71},
		{"COL": "I", "WIDTH": 16.83},
		{"COL": "J", "WIDTH": 10.10},
		{"COL": "K", "WIDTH": 17.65},
		{"COL": "L", "WIDTH": 22.14},
	}
	for _, v := range arrStyleColWidth {
		tmpColWidth := fmt.Sprintf("%f", v["WIDTH"])
		colWidth, _ := strconv.ParseFloat(tmpColWidth, 64)
		err = f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			return nil, err
		}
	}

	styleDefault, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		return nil, err
	}
	err = f.SetColStyle(sheet, "A:Z", styleDefault)
	if err != nil {
		return nil, err
	}
	numberFormat := "#,##"
	stylingBorderRightOnly, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	stylingBorderLeftOnly, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
		},
	})

	stylingBorderLROnly, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
		},
	})

	stylingBorderTopOnly, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	stylingHeader, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#fac090"},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}

	stylingSubTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return nil, err
	}

	stylingSubTotalCurrency, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}

	stylingTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#ccff33"},
		},
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}
	stylingTotalControl, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#3ada24"},
		},
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}

	err = f.MergeCell(sheet, "B6", "B8")
	if err != nil {
		return nil, err
	}
	err = f.MergeCell(sheet, "C6", "E8")
	if err != nil {
		return nil, err
	}
	err = f.MergeCell(sheet, "F6", "F8")
	if err != nil {
		return nil, err
	}
	err = f.MergeCell(sheet, "H6", "K7")
	if err != nil {
		return nil, err
	}

	stylingCurrency, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	stylingCurrency2, err := f.NewStyle(&excelize.Style{
		NumFmt: 7,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}
	err = f.SetCellStyle(sheet, "B6", "L8", stylingHeader)
	if err != nil {
		return nil, err
	}

	err = f.SetCellValue(sheet, "B6", "No Akun")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "C6", "Keterangan")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "F6", "WP Reff")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "G6", "")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "G7", "Unaudited")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "G8", "")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "H6", "Adjustment Journal Entry")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "I8", "Debet")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "K8", "Kredit")
	if err != nil {
		return nil, err
	}
	err = f.SetCellFormula(sheet, "L6", "=G6")
	if err != nil {
		return nil, err
	}
	err = f.SetCellFormula(sheet, "L7", "=G7")
	if err != nil {
		return nil, err
	}
	err = f.SetCellFormula(sheet, "L8", "=G8")
	if err != nil {
		return nil, err
	}

	row := 9

	rowCode := make(map[string]int)
	isAutoSum := make(map[string]bool)
	customRow := make(map[string]string)
	for _, v := range *data {
		rowCode[v.Code] = row
		if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("L%d", row), stylingCurrency); err != nil {
			return nil, err
		}
		if err = f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingBorderLROnly); err != nil {
			return nil, err
		}
		if v.AutoSummary != nil && *v.AutoSummary {
			isAutoSum[v.Code] = true
		}
		if !(v.IsTotal != nil && *v.IsTotal) && v.IsLabel != nil && *v.IsLabel {
			if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("L%d", row), stylingCurrency); err != nil {
				return nil, err
			}
			if err = f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingBorderLROnly); err != nil {
				return nil, err
			}
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)
		}

		if v.IsCoa != nil && *v.IsCoa {
			rowBefore := row
			tbdetails, err := s.CoaRepository.FindWithCode(ctx, &v.Code)
			if err != nil {
				return nil, err
			}
			for _, vTbDetail := range *tbdetails {
				if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("L%d", row), stylingCurrency); err != nil {
					return nil, err
				}
				if err = f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingBorderLROnly); err != nil {
					return nil, err
				}
				if strings.Contains(strings.ToLower(vTbDetail.Code), "subtotal") {
					continue
				}

				tbRowCode[vTbDetail.Code] = row
				if vTbDetail.Code == "310401004" || vTbDetail.Code == "310402002" || vTbDetail.Code == "310501002" || vTbDetail.Code == "310502002" || vTbDetail.Code == "310503002" {
					f.SetCellValue(sheet, fmt.Sprintf("M%d", row), "control")
				}
				valueKosong := 0.0
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vTbDetail.Code)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), vTbDetail.Name)
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), valueKosong)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), valueKosong)
				f.SetCellValue(sheet, fmt.Sprintf("K%d", row), valueKosong)
				if err = f.SetCellStyle(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("L%d", row), stylingCurrency); err != nil {
					return nil, err
				}
				if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), stylingCurrency); err != nil {
					return nil, err
				}

				if err = f.SetCellStyle(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), stylingCurrency2); err != nil {
					return nil, err
				}
				if err = f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingCurrency2); err != nil {
					return nil, err
				}
				tmpHeadCoa := fmt.Sprintf("%c", vTbDetail.Code[0])
				if tmpHeadCoa == "9" {
					tmpHeadCoa = vTbDetail.Code[:1]
				}
				switch tmpHeadCoa {
				case "1", "5", "6", "7", "91", "92":
					f.SetCellFormula(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("=G%d+I%d-K%d", row, row, row))
				default:
					f.SetCellFormula(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("=G%d-I%d+K%d", row, row, row))
				}
				row++
			}
			rowAfter := row - 1
			rowTB := len(*tbdetails)
			if v.AutoSummary != nil && *v.AutoSummary && rowTB > 1 {
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "Subtotal")
				f.SetCellFormula(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("=SUM(G%d:G%d)", rowBefore, rowAfter))
				f.SetCellFormula(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("=SUM(I%d:I%d)", rowBefore, rowAfter))
				f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("=SUM(K%d:K%d)", rowBefore, rowAfter))
				f.SetCellFormula(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("=SUM(L%d:L%d)", rowBefore, rowAfter))
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("F%d", row), stylingSubTotal)
				if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("L%d", row), stylingSubTotalCurrency); err != nil {
					return nil, err
				}
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)] = row
				tbRowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)] = row
				row++
				if err = f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("L%d", row), stylingBorderLROnly); err != nil {
					return nil, err
				}
			}
		}

		if v.IsTotal != nil && *v.IsTotal {
			tbRowCode[v.Code] = row
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)
			if v.Code == "TOTAL_INVESTASI_JANGKA_PENDEK" || v.Code == "TOTAL_PIUTANG_LAIN~LAIN_~_PIHAK_KETIGA~PIUTANG_USAHA" || v.Code == "TOTAL_CIP~ASET_TETAP" || v.Code == "TOTAL_PIUTANG_LAIN~LAIN" || v.Code == "TOTAL_CIP~ASET_TAK_BERWUJUD" || v.Code == "NCA_~_RUA_~_AKUMULASI_PENYUSUTAN_~_LAIN~LAIN" || v.Code == "TOTAL_INVESTASI_JANGKA_PANJANG" || v.Code == "TOTAL_PIUTANG_LAIN~LAIN_~_JANGKA_PANJANG" || v.Code == "TOTAL_ASET" || v.Code == "TOTAL_UTANG_USAHA" || v.Code == "TOTAL_UTANG_LAIN~LAIN_JANGKA_PENDEK" || v.Code == "TOTAL_LIABILITAS_IMBALAN_KERJA_~_JANGKA_PENDEK" || v.Code == "TOTAL_UTANG_LAIN~LAIN_JANGKA_PANJANG" {
				f.SetCellValue(sheet, fmt.Sprintf("M%d", row), "control")
			}
			if v.Code == "TOTAL_LIABILITAS_DAN_EKUITAS" {
				rowAset := row
				if _, ok := rowCode["TOTAL_ASET"]; ok {
					rowAset = rowCode["TOTAL_ASET"]
				}
				f.SetCellFormula(sheet, "G5", fmt.Sprintf("=G%d-G%d", rowAset, row))
				f.SetCellFormula(sheet, "L5", fmt.Sprintf("=L%d-L%d", rowAset, row))
			}

			//show control aje
			if v.Code == "CONTROL" {
				f.SetCellFormula(sheet, "I5", fmt.Sprintf("=I%d", row))
				f.SetCellFormula(sheet, "K5", fmt.Sprintf("=K%d", row))
			}
			if v.Code == "CONTROL_TO_ADJUSTMENT_SHEET" {
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)
				dbt := 0.0

				cdt := 0.0

				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), dbt)
				f.SetCellValue(sheet, fmt.Sprintf("K%d", row), cdt)
				f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=I%d-K%d", row, row))
			}
			if v.Code == "TOTAL_JOURNAL_IN_WP" {
				f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=I%d-K%d", row, row))
			}

			if v.Code != "TOTAL_JOURNAL_IN_WP" && v.Code != "CONTROL_TO_ADJUSTMENT_SHEET" && v.Code != "CONTROL" && v.Code != "CONTROL_TO_WBS_1" {
				if v.IsControl != nil && *v.IsControl {
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("L%d", row), stylingTotalControl)
				} else {
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("L%d", row), stylingTotal)
				}
			}
			if v.FxSummary == "" {
				row++
				continue
			}
			for chr := 'G'; chr <= 'L'; chr++ {
				if fmt.Sprintf("%c", chr) == "H" || fmt.Sprintf("%c", chr) == "J" || ((fmt.Sprintf("%c", chr) == "G" || fmt.Sprintf("%c", chr) == "L") && (v.Code == "TOTAL_JOURNAL_IN_WP" || v.Code == "CONTROL_TO_ADJUSTMENT_SHEET" || v.Code == "CONTROL")) {
					continue
				}
				formula := v.FxSummary
				// reg := regexp.MustCompile(`[A-Za-z_]+|[0-9]+\d{3,}`)
				reg := regexp.MustCompile(`[A-Za-z0-9_~:()]+|[0-9]+\d{2,}`)
				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					if len(vMatch) < 3 {
						continue
					}
					if isAutoSum[vMatch] {
						if rowCode[fmt.Sprintf("%s_SUBTOTAL", vMatch)] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%c%d", chr, rowCode[fmt.Sprintf("%s_SUBTOTAL", vMatch)]))
						} else {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%c%d", chr, rowCode[vMatch]))
						}
					} else {
						if rowCode[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%c%d", chr, rowCode[vMatch]))
						}
					}
				}
				f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=%s", formula))
			}
			row++
			continue
		}
		row++
	}

	customRow["310401004"] = "=LABA_BERSIH"
	customRow["310402002"] = "=TOTAL_PENGHASILAN_KOMPREHENSIF_LAIN~BS-SUM(310501002,310502002,310503002)"
	customRow["310501002"] = "=950101001"
	customRow["310502002"] = "=950301001+950301002"
	customRow["310503002"] = "=950401001+950401002"
	for key, nRow := range tbRowCode {
		if strings.Contains(customRow["310401004"], key) {
			customRow["310401004"] = strings.ReplaceAll(customRow["310401004"], key, fmt.Sprintf("@%d", nRow))
		}
		if strings.Contains(customRow["310402002"], key) && key != "RE" {
			customRow["310402002"] = strings.ReplaceAll(customRow["310402002"], key, fmt.Sprintf("@%d", nRow))
		}
		if strings.Contains(customRow["310501002"], key) {
			customRow["310501002"] = strings.ReplaceAll(customRow["310501002"], key, fmt.Sprintf("@%d", nRow))
		}
		if strings.Contains(customRow["310502002"], key) {
			customRow["310502002"] = strings.ReplaceAll(customRow["310502002"], key, fmt.Sprintf("@%d", nRow))
		}
		if strings.Contains(customRow["310503002"], key) {
			customRow["310503002"] = strings.ReplaceAll(customRow["310503002"], key, fmt.Sprintf("@%d", nRow))
		}
	}

	for key, vCustomRow := range customRow {
		if val, ok := tbRowCode[key]; ok {
			f.SetCellFormula(sheet, fmt.Sprintf("G%d", val), strings.ReplaceAll(vCustomRow, "@", "G"))
			f.SetCellFormula(sheet, fmt.Sprintf("I%d", val), strings.ReplaceAll(vCustomRow, "@", "I"))
			f.SetCellFormula(sheet, fmt.Sprintf("K%d", val), strings.ReplaceAll(vCustomRow, "@", "K"))
			// f.SetCellFormula(sheet, fmt.Sprintf("L%d", val), strings.ReplaceAll(vCustomRow, "@", "L"))
		}
	}

	if err = f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("L%d", row), stylingBorderTopOnly); err != nil {
		return nil, err
	}

	if err = f.SetCellStyle(sheet, "A9", fmt.Sprintf("A%d", row-1), stylingBorderRightOnly); err != nil {
		return nil, err
	}
	if err = f.SetCellStyle(sheet, "M9", fmt.Sprintf("M%d", row-1), stylingBorderLeftOnly); err != nil {
		return nil, err
	}

	if err = f.SetSheetFormatPr(sheet, excelize.DefaultRowHeight(12.85)); err != nil {
		return nil, err
	}
	f.SetDefaultFont("Arial")

	return f, nil
}
func (s *service) ExportInvestasiTbk(ctx *abstraction.Context, f *excelize.File) (*excelize.File, error) {
	sheet := "INVESTASI_TBK"
	f.NewSheet(sheet)

	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 8.21},
		{"COL": "B", "WIDTH": 4.10},
		{"COL": "C", "WIDTH": 33.39},
		{"COL": "D", "WIDTH": 11.78},
		{"COL": "E", "WIDTH": 12.31},
		{"COL": "F", "WIDTH": 15.71},
		{"COL": "G", "WIDTH": 12.78},
		{"COL": "H", "WIDTH": 15.35},
		{"COL": "I", "WIDTH": 15.35},
		{"COL": "J", "WIDTH": 15.71},
		{"COL": "K", "WIDTH": 14.10},
	}
	for _, v := range arrStyleColWidth {
		tmpColWidth := fmt.Sprintf("%f", v["WIDTH"])
		colWidth, _ := strconv.ParseFloat(tmpColWidth, 64)
		err := f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			return nil, err
		}
	}

	styleHeader, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold:   true,
			Family: "Arial",
			Size:   10,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#fcd5b4"},
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}
	numberFormat := "#,##"
	styleCurrency, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &numberFormat,
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		return nil, err
	}
	styleCurrencySum, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		NumFmt: 41,
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		return nil, err
	}
	styleCurrencyTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold:   true,
			Family: "Arial",
			Size:   10,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#99ff33"},
		},
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}

	styleLabelTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold:   true,
			Family: "Arial",
			Size:   10,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#99ff33"},
		},
	})
	if err != nil {
		return nil, err
	}

	styleLabel, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		return nil, err
	}

	f.SetCellStyle(sheet, "B4", "K6", styleHeader)
	f.SetCellValue(sheet, "B4", "No")
	f.SetCellValue(sheet, "C4", "Stock")
	f.SetCellValue(sheet, "D4", "Ending Share")
	f.SetCellValue(sheet, "E4", "AVG Price")
	f.SetCellValue(sheet, "F4", "Amount (Cost)")
	f.SetCellValue(sheet, "G4", fmt.Sprintf("Closing Price"))
	f.SetCellValue(sheet, "H4", "Amount (FV)")
	f.SetCellValue(sheet, "I4", "Unrealized Gain(oss)")
	f.SetCellValue(sheet, "J4", "Realized Gain(loss)")
	f.SetCellValue(sheet, "K4", "Fee")

	formatterCode := []string{"INVESTASI-TBK"}
	formatterTitle := []string{"Summary Investasi Tbk"}
	row, rowStart := 5, 5

	f.SetCellValue(sheet, "B2", formatterTitle[0])

	for _, formatter := range formatterCode {

		var criteria dto.FormatterGetRequest
		criteria.FormatterFilterModel.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
		if err != nil {
			return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
		}

		rowCode := make(map[string]int)
		partRowStart := row
		for _, v := range data.FormatterDetail {
			rowCode[v.Code] = row
			f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabel)
			f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("K%d", row), styleCurrency)

			if strings.Contains(strings.ToLower(v.Code), "blank") {
				row++
				continue
			}
			valueKosong := 0.0
			if strings.ToUpper(v.Code) == "TOTAL" {
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)
				f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), styleCurrencySum)
				f.SetCellStyle(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), styleCurrencySum)
				f.SetCellStyle(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), styleCurrencySum)
			} else {
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), row-4)
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), valueKosong)
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), valueKosong)
				f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("=G%d*D%d", row, row))
				f.SetCellFormula(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("=D%d*E%d", row, row))
				f.SetCellFormula(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("=H%d-F%d", row, row))
				f.SetCellValue(sheet, fmt.Sprintf("J%d", row), valueKosong)
				f.SetCellValue(sheet, fmt.Sprintf("K%d", row), valueKosong)
				f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), styleCurrencySum)
				f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), styleCurrencySum)
				f.SetCellStyle(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), styleCurrencySum)
				f.SetCellStyle(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), styleCurrencySum)
			}

			if v.IsTotal != nil && *v.IsTotal {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("K%d", row), styleCurrencyTotal)
				if v.FxSummary == "" {
					row++
					continue
				}
				for chr := 'D'; chr <= 'K'; chr++ {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[0-9]+\d{3,}`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						//cari jml berdasarkan code
						if rowCode[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%c%d", chr, rowCode[vMatch]))
						}

					}
					f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=%s", formula))
				}
				row++
				continue
			}

			if v.AutoSummary != nil && *v.AutoSummary {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("K%d", row), styleCurrencyTotal)
				f.SetCellFormula(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("=SUM(F%d:F%d)", partRowStart, row-1))
				for chr := 'H'; chr <= 'K'; chr++ {
					f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
				}
				row++
				partRowStart = row
				continue
			}
			row++

		}
		rowStart = row
		row = rowStart
	}

	return f, nil
}
func (s *service) ExportMutasiDta(ctx *abstraction.Context, f *excelize.File) (*excelize.File, error) {
	sheet := "MUTASI_DTA"
	f.NewSheet(sheet)

	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 8.21},
		{"COL": "B", "WIDTH": 3.50},
		{"COL": "C", "WIDTH": 33.39},
		{"COL": "D", "WIDTH": 15.35},
		{"COL": "E", "WIDTH": 15.35},
		{"COL": "F", "WIDTH": 13.74},
		{"COL": "G", "WIDTH": 13.74},
		{"COL": "H", "WIDTH": 13.74},
		{"COL": "I", "WIDTH": 15.71},
	}
	for _, v := range arrStyleColWidth {
		tmpColWidth := fmt.Sprintf("%f", v["WIDTH"])
		colWidth, _ := strconv.ParseFloat(tmpColWidth, 64)
		err := f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			return nil, err
		}
	}
	stylingControl, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#FFFF00"},
		},
	})
	if err != nil {
		return nil, err
	}
	styleHeader, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold:   true,
			Family: "Arial",
			Size:   10,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#fac090"},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}
	numberFormat := "#,##"
	styleCurrency, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &numberFormat,
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		return nil, err
	}

	styleCurrencyTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold:   true,
			Family: "Arial",
			Size:   10,
		},
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}

	styleLabelTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold:   true,
			Family: "Arial",
			Size:   10,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}

	styleLabel, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		return nil, err
	}

	stylingDefault, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}

	formatterCode := []string{"MUTASI-DTA"}
	formatterTitle := []string{"Mutasi DTA"}
	row, rowStart := 7, 7

	f.SetCellValue(sheet, "B2", formatterTitle[0])

	f.MergeCell(sheet, "B4", "B6")
	f.MergeCell(sheet, "C4", "C6")
	// f.MergeCell(sheet, "D4", "D6")
	f.MergeCell(sheet, "E4", "G4")
	f.MergeCell(sheet, "E5", "E6")
	f.MergeCell(sheet, "F5", "F6")
	f.MergeCell(sheet, "G5", "G6")
	f.MergeCell(sheet, "H4", "I4")
	f.MergeCell(sheet, "H5", "H6")
	f.MergeCell(sheet, "I5", "I6")
	// f.MergeCell(sheet, "J4", "J6")

	f.SetCellStyle(sheet, "B4", "J6", styleHeader)
	f.SetCellValue(sheet, "B4", "NO")
	f.SetCellValue(sheet, "C4", "Description")
	f.SetCellValue(sheet, "D4", "Saldo Awal")
	f.SetCellFormula(sheet, "D5" , "=TRIAL_BALANCE!G6")
	f.SetCellValue(sheet, "D6", "31-Dec-21")
	f.SetCellValue(sheet, "E4", "Penambahan (Pengurangan)")
	f.SetCellValue(sheet, "E5", "Manfaat (beban) pajak tangguhan")
	f.SetCellValue(sheet, "F5", "OCI")
	f.SetCellValue(sheet, "G5", "Akuisisi Entitas anak")
	f.SetCellValue(sheet, "H4", "Dampak perubahan tariff pajak")
	f.SetCellValue(sheet, "H5", "Dibebankan ke laba rugi")
	f.SetCellValue(sheet, "I5", "Dibebankan ke OCI")
	f.SetCellValue(sheet, "J4", "Saldo Akhir")
	f.SetCellValue(sheet, "J5", "Company")
	f.SetCellValue(sheet, "J6", "Periode")

	for _, formatter := range formatterCode {

		var criteria dto.FormatterGetRequest
		criteria.FormatterFilterModel.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
		if err != nil {
			return nil, helper.ErrorHandler(err)
		}

		rowCode := make(map[string]int)
		partRowStart := row
		for _, v := range data.FormatterDetail {
			rowCode[v.Code] = row
			f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabel)
			f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), styleCurrency)

			if strings.Contains(strings.ToLower(v.Code), "blank") {
				row++
				continue
			}

			valueKosong := 0.0
			if strings.ToUpper(v.Code) == "TOTALS" || strings.ToUpper(v.Code) == "CONTROL_1" || strings.ToUpper(v.Code) == "CONTROL_2" || strings.ToUpper(v.Code) == "CONTROL_3" {
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)
				// f.SetCellValue(sheet, fmt.Sprintf("D%d", row), valueKosong)
				// f.SetCellValue(sheet, fmt.Sprintf("E%d", row), valueKosong)
				// f.SetCellValue(sheet, fmt.Sprintf("F%d", row), valueKosong)
				// f.SetCellValue(sheet, fmt.Sprintf("G%d", row), valueKosong)
				// f.SetCellValue(sheet, fmt.Sprintf("H%d", row), valueKosong)
				// f.SetCellValue(sheet, fmt.Sprintf("I%d", row), valueKosong)
			} else {
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), row-7)
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), valueKosong)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), valueKosong)
				f.SetCellValue(sheet, fmt.Sprintf("F%d", row), valueKosong)
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), valueKosong)
				f.SetCellValue(sheet, fmt.Sprintf("H%d", row), valueKosong)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), valueKosong)
				f.SetCellStyle(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("J%d", row), styleCurrencyTotal)
			}

			if v.IsTotal != nil && *v.IsTotal {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), stylingDefault)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), stylingDefault)
				if v.FxSummary == "" {
					row++
					continue
				}
				arrChr := []string{"D", "E", "F", "G", "H", "I", "J", "K"}

				if strings.ToUpper(v.Code) == "CONTROL_1" {
					arrChr = []string{"E"}
					// f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("D%d", row), stylingDefault)
					// f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("I%d", row), stylingDefault)
				}
				if strings.ToUpper(v.Code) == "CONTROL_2" {
					arrChr = []string{"F"}

				}
				if strings.ToUpper(v.Code) == "CONTROL_3" {
					arrChr = []string{"J"}
				}
				for _, chr := range arrChr {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						//cari jml berdasarkan code
						if rowCode[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCode[vMatch]))
						}
						if _, ok := tbRowCode[vMatch]; ok {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCode[vMatch]))
							f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCode[vMatch]), "control")
						}

					}
					if strings.ToUpper(v.Code) == "CONTROL_1" {
						row = row
						f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_2" {
						row = row - 1
						f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_3" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("J%d", row), stylingControl)
					}
					f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row), fmt.Sprintf("=%s", formula))
				}
				row++
				continue
			}

			if v.AutoSummary != nil && *v.AutoSummary {
				f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), styleCurrencyTotal)
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)

				for chr := 'D'; chr <= 'J'; chr++ {
					f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
				}
				row++
				partRowStart = row
				continue
			}
			f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=SUM(D%d:I%d)", row, row))
			row++
		}
		rowStart = row
		row = rowStart
	}

	return f, nil
}
func (s *service) ExportAgingUtangPiutang(ctx *abstraction.Context, f *excelize.File) (*excelize.File, error) {
	sheet := "AGING_UTANG_PIUTANG"
	f.NewSheet(sheet)
	f.SetColWidth(sheet, "A", "A", 8.29)
	f.SetColWidth(sheet, "B", "B", 26.43)
	f.SetColWidth(sheet, "C", "S", 17.86)

	styleHeader, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#f8cbad"},
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}
	numberFormat := "#,##"
	styleCurrency, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}

	styleCurrencyTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#99ff66"},
		},
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}

	styleLabelTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#99ff66"},
		},
	})
	if err != nil {
		return nil, err
	}

	styleLabel, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	stylingControl, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#FFFF00"},
		},
	})
	if err != nil {
		return nil, err
	}

	stylingDefault, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}

	f.SetColWidth(sheet, "K", "K", 2)

	formatterCode := []string{"AGING-UTANG-PIUTANG", "AGING-UTANG-PIUTANG-MUTASI-ECL"}
	formatterTitle := []string{"Detail Aging", "Mutasi ECL's"}
	row, rowStart := 5, 5
	for i, formatter := range formatterCode {
		f.SetCellValue(sheet, fmt.Sprintf("B%d", (rowStart-3)), formatterTitle[i])
		f.SetRowHeight(sheet, (rowStart - 2), 43.50)
		f.MergeCell(sheet, fmt.Sprintf("B%d", (rowStart-2)), fmt.Sprintf("B%d", (rowStart-1)))
		f.MergeCell(sheet, fmt.Sprintf("G%d", (rowStart-2)), fmt.Sprintf("G%d", (rowStart-1)))
		f.MergeCell(sheet, fmt.Sprintf("J%d", (rowStart-2)), fmt.Sprintf("J%d", (rowStart-1)))
		f.MergeCell(sheet, fmt.Sprintf("P%d", (rowStart-2)), fmt.Sprintf("P%d", (rowStart-1)))
		f.MergeCell(sheet, fmt.Sprintf("S%d", (rowStart-2)), fmt.Sprintf("S%d", (rowStart-1)))

		f.MergeCell(sheet, fmt.Sprintf("C%d", (rowStart-2)), fmt.Sprintf("D%d", (rowStart-2)))
		f.MergeCell(sheet, fmt.Sprintf("E%d", (rowStart-2)), fmt.Sprintf("F%d", (rowStart-2)))
		f.MergeCell(sheet, fmt.Sprintf("H%d", (rowStart-2)), fmt.Sprintf("I%d", (rowStart-2)))
		f.MergeCell(sheet, fmt.Sprintf("L%d", (rowStart-2)), fmt.Sprintf("M%d", (rowStart-2)))
		f.MergeCell(sheet, fmt.Sprintf("N%d", (rowStart-2)), fmt.Sprintf("O%d", (rowStart-2)))
		f.MergeCell(sheet, fmt.Sprintf("Q%d", (rowStart-2)), fmt.Sprintf("R%d", (rowStart-2)))
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", (rowStart-2)), fmt.Sprintf("J%d", (rowStart-1)), styleHeader)
		f.SetCellStyle(sheet, fmt.Sprintf("L%d", (rowStart-2)), fmt.Sprintf("S%d", (rowStart-1)), styleHeader)

		f.SetCellValue(sheet, fmt.Sprintf("B%d", (rowStart-2)), "Description")
		f.SetCellValue(sheet, fmt.Sprintf("C%d", (rowStart-2)), "Piutang Usaha")
		f.SetCellValue(sheet, fmt.Sprintf("E%d", (rowStart-2)), "Piutang lain-lain jangka pendek")
		f.SetCellValue(sheet, fmt.Sprintf("G%d", (rowStart-2)), "Piutang pihak berelasi jangka pendek")
		f.SetCellValue(sheet, fmt.Sprintf("H%d", (rowStart-2)), "Piutang lain-lain jangka panjang")
		f.SetCellValue(sheet, fmt.Sprintf("J%d", (rowStart-2)), "Piutang pihak berelasi jangka panjang")
		f.SetCellValue(sheet, fmt.Sprintf("L%d", (rowStart-2)), "Utang usaha")
		f.SetCellValue(sheet, fmt.Sprintf("N%d", (rowStart-2)), "Utang lain-lain jangka pendek")
		f.SetCellValue(sheet, fmt.Sprintf("P%d", (rowStart-2)), "Utang pihak berelasi jangka pendek")
		f.SetCellValue(sheet, fmt.Sprintf("Q%d", (rowStart-2)), "Utang lain-lain jangka panjang")
		f.SetCellValue(sheet, fmt.Sprintf("S%d", (rowStart-2)), "Utang pihak berelasi jangka panjang")

		header1 := []string{"C", "D", "E", "F", "H", "I", "L", "M", "N", "O", "Q", "R"}
		for i, v := range header1 {
			if (i+1)%2 == 0 {
				f.SetCellValue(sheet, fmt.Sprintf("%s%d", v, (rowStart-1)), "Pihak Berelasi")
			} else {
				f.SetCellValue(sheet, fmt.Sprintf("%s%d", v, (rowStart-1)), "Pihak Ketiga")
			}
		}

		var criteria dto.FormatterGetRequest
		criteria.FormatterFilterModel.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
		if err != nil {
			return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
		}

		rowCode := make(map[string]int)
		partRowStart := row
		for _, v := range data.FormatterDetail {
			rowCode[v.Code] = row
			f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabel)
			f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("J%d", row), styleCurrency)
			f.SetCellStyle(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("S%d", row), styleCurrency)

			if strings.Contains(strings.ToLower(v.Code), "blank") {
				row++
				continue
			}

			rowKosong := 0
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
			
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), rowKosong)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), rowKosong)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), rowKosong)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), rowKosong)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), rowKosong)
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), rowKosong)
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), rowKosong)
			f.SetCellValue(sheet, fmt.Sprintf("J%d", row), rowKosong)
			f.SetCellValue(sheet, fmt.Sprintf("L%d", row), rowKosong)
			f.SetCellValue(sheet, fmt.Sprintf("M%d", row), rowKosong)
			f.SetCellValue(sheet, fmt.Sprintf("N%d", row), rowKosong)
			f.SetCellValue(sheet, fmt.Sprintf("O%d", row), rowKosong)
			f.SetCellValue(sheet, fmt.Sprintf("P%d", row), rowKosong)
			f.SetCellValue(sheet, fmt.Sprintf("P%d", row), rowKosong)
			f.SetCellValue(sheet, fmt.Sprintf("Q%d", row), rowKosong)
			f.SetCellValue(sheet, fmt.Sprintf("S%d", row), rowKosong)

			if v.IsTotal != nil && *v.IsTotal {
				if v.FxSummary == "" {
					row++
					continue
				}
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("J%d", row), stylingDefault)
				f.SetCellStyle(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("S%d", row), stylingDefault)

				arrChr := []string{"E", "F"}
				if strings.ToUpper(v.Code) == "CONTROL_1" {
					arrChr = []string{"C"}
					f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), stylingControl)

				}
				if strings.ToUpper(v.Code) == "CONTROL_17" {
					arrChr = []string{"C"}

					f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), stylingControl)
				}
				if strings.ToUpper(v.Code) == "CONTROL_18" {
					arrChr = []string{"C"}

					f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), stylingControl)
				}
				if strings.ToUpper(v.Code) == "CONTROL_2" {
					arrChr = []string{"D"}

				}
				if strings.ToUpper(v.Code) == "CONTROL_3" {
					arrChr = []string{"E"}

				}
				if strings.ToUpper(v.Code) == "CONTROL_4" {
					arrChr = []string{"F"}

				}
				if strings.ToUpper(v.Code) == "CONTROL_5" {
					arrChr = []string{"G"}

				}
				if strings.ToUpper(v.Code) == "CONTROL_6" {
					arrChr = []string{"H"}

				}
				if strings.ToUpper(v.Code) == "CONTROL_7" {
					arrChr = []string{"I"}

				}
				if strings.ToUpper(v.Code) == "CONTROL_8" {
					arrChr = []string{"J"}

				}
				if strings.ToUpper(v.Code) == "CONTROL_9" {
					arrChr = []string{"L"}

				}
				if strings.ToUpper(v.Code) == "CONTROL_10" {
					arrChr = []string{"M"}

				}
				if strings.ToUpper(v.Code) == "CONTROL_11" {
					arrChr = []string{"N"}

				}
				if strings.ToUpper(v.Code) == "CONTROL_12" {
					arrChr = []string{"O"}

				}
				if strings.ToUpper(v.Code) == "CONTROL_13" {
					arrChr = []string{"P"}

				}
				if strings.ToUpper(v.Code) == "CONTROL_14" {
					arrChr = []string{"Q"}

				}
				if strings.ToUpper(v.Code) == "CONTROL_15" {
					arrChr = []string{"R"}

				}
				if strings.ToUpper(v.Code) == "CONTROL_16" {
					arrChr = []string{"S"}

				}
				if strings.ToUpper(v.Code) == "CONTROL_17" {
					arrChr = []string{"C"}

				}
				if strings.ToUpper(v.Code) == "CONTROL_18" {
					arrChr = []string{"C"}

				}
				for _, chr := range arrChr {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[A-Za-z0-9_~#:'()]+|[0-9]+\d{3,}`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						//cari jml berdasarkan code
						if rowCode[vMatch] != 0 {
							if rowCode[vMatch] != 0 {
								formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCode[vMatch]))
							}
						}
						if _, ok := tbRowCode[vMatch]; ok {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCode[vMatch]))
							f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCode[vMatch]), "control")
						}
					}

					if strings.ToUpper(v.Code) == "CONTROL_2" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_3" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_4" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_5" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_6" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_7" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_8" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("J%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_9" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("L%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_10" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("M%d", row), fmt.Sprintf("M%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_11" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("N%d", row), fmt.Sprintf("N%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_12" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("O%d", row), fmt.Sprintf("O%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_13" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("P%d", row), fmt.Sprintf("P%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_14" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("Q%d", row), fmt.Sprintf("Q%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_15" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("R%d", row), fmt.Sprintf("R%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_16" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("S%d", row), fmt.Sprintf("S%d", row), stylingControl)
					}
					err = f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row), fmt.Sprintf("=%s", formula))
					if err != nil {
						fmt.Println(err)
					}
				}
				row++
				continue
			}

			if v.AutoSummary != nil && *v.AutoSummary {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("J%d", row), styleCurrencyTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("S%d", row), styleCurrencyTotal)

				for chr := 'C'; chr <= 'S'; chr++ {
					if chr == 'K' {
						continue
					}
					f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
				}
				row++
				partRowStart = row
				continue
			}

			row++
		}
		rowStart = row + 4
		row = rowStart
	}

	return f, nil
}
func (s *service) ExportPembelianPenjualanBerelasi(ctx *abstraction.Context, f *excelize.File) (*excelize.File, error) {
	sheet := "PEMBELIAN_PENJUALAN_BERELASI"
	f.NewSheet(sheet)
	err := f.SetColWidth(sheet, "D", "D", 66)
	if err != nil {
		return nil, err
	}
	err = f.SetColWidth(sheet, "E", "E", 13)
	if err != nil {
		return nil, err
	}
	err = f.SetColWidth(sheet, "F", "F", 13)
	if err != nil {
		return nil, err
	}
	styleHeader, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#f8cbad"},
		},
		Alignment: &excelize.Alignment{
			WrapText: true,
		},
	})
	if err != nil {
		return nil, err
	}

	numberFormat := "#,##"
	styleCurrency, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}

	styleCurrencyTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#99ff66"},
		},
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}

	styleLabelTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#99ff66"},
		},
	})
	if err != nil {
		return nil, err
	}

	styleLabel, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	f.SetCellStyle(sheet, "B4", "F4", styleHeader)
	f.SetCellValue(sheet, "C2", "List pembelian dan penjualan berelasi")
	f.SetCellValue(sheet, "B4", "NO")
	f.SetCellValue(sheet, "C4", "code")
	f.SetCellValue(sheet, "D4", "PT")
	f.SetCellValue(sheet, "E4", "Pembelian")
	f.SetCellValue(sheet, "F4", "Penjualan")

	t := true
	data, err := s.CompanyRepository.FindIsActive(ctx, &t)
	if err != nil {
		return nil, err
	}

	row, rowStart := 5, 5
	for i, v := range *data {
		valueKosong := 0.0
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), (i + 1))
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Code)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), v.Name)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), valueKosong)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), valueKosong)

		row++
	}

	f.SetCellStyle(sheet, "B5", fmt.Sprintf("D%d", row), styleLabel)
	f.SetCellStyle(sheet, "E5", fmt.Sprintf("F%d", row), styleCurrency)
	f.SetCellStyle(sheet, fmt.Sprintf("B%d", row+1), fmt.Sprintf("D%d", row+1), styleLabelTotal)
	f.SetCellStyle(sheet, fmt.Sprintf("E%d", row+1), fmt.Sprintf("F%d", row+1), styleCurrencyTotal)

	f.SetCellValue(sheet, fmt.Sprintf("D%d", row+1), "Total")
	f.SetCellFormula(sheet, fmt.Sprintf("E%d", row+1), fmt.Sprintf("=SUM(E%d:E%d)", rowStart, row))
	f.SetCellFormula(sheet, fmt.Sprintf("F%d", row+1), fmt.Sprintf("=SUM(F%d:F%d)", rowStart, row))

	return f, nil
}

var rowCodePersediaan = make(map[string]int)

func (s *service) ExportMutasiPersediaan(ctx *abstraction.Context, f *excelize.File) (*excelize.File, error) {
	sheet := "MUTASI_PERSEDIAAN"
	f.NewSheet(sheet)
	f.SetColWidth(sheet, "B", "B", 31.01)
	f.SetColWidth(sheet, "C", "C", 13)
	f.SetColWidth(sheet, "D", "D", 13)

	styleHeader, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#f8cbad"},
		},
		Alignment: &excelize.Alignment{
			WrapText: true,
		},
	})
	if err != nil {
		return nil, err
	}

	numberFormat := "#,##"
	styleCurrency, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}
	stylingControl, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#FFFF00"},
		},
	})
	if err != nil {
		return nil, err
	}
	styleCurrencyTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#99ff66"},
		},
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}

	styleLabelTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#99ff66"},
		},
	})
	if err != nil {
		return nil, err
	}

	styleLabel, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	formatterCode := []string{"MUTASI-PERSEDIAAN", "MUTASI-CADANGAN-PENGHAPUSAN-PERSEDIAAN"}
	formatterTitle := []string{"Mutasi Persediaan", "Mutasi Cadangan penghapusan persediaan"}
	row, rowStart := 4, 4

	for i, formatter := range formatterCode {
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", (rowStart-1)), fmt.Sprintf("C%d", (rowStart-1)), styleHeader)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", (rowStart-2)), formatterTitle[i])
		f.SetCellValue(sheet, fmt.Sprintf("B%d", (rowStart-1)), "Description")
		f.SetCellValue(sheet, fmt.Sprintf("C%d", (rowStart-1)), "Amount")

		var criteria model.FormatterFilterModel
		criteria.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria)
		if err != nil {
			return nil, helper.ErrorHandler(err)
		}

		partRowStart := row
		for _, v := range data.FormatterDetail {
			rowCodePersediaan[v.Code] = row
			f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabel)
			f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleCurrency)

			if strings.Contains(strings.ToLower(v.Code), "blank") {
				row++
				continue
			}
			valueKosong := 0.0
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), valueKosong)
			f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleCurrency)

			if v.IsTotal != nil && *v.IsTotal {
				if v.FxSummary == "" {
					row++
					continue
				}
				arrChr := []string{"D", "E", "F", "G", "H", "I", "J", "K"}

				if strings.ToUpper(v.Code) == "CONTROL_1" {
					arrChr = []string{"C"}
					f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), stylingControl)
				}
				for _, chr := range arrChr {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						//cari jml berdasarkan code
						if rowCodePersediaan[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCodePersediaan[vMatch]))
						}
						if _, ok := tbRowCode[vMatch]; ok {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCode[vMatch]))
							f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCode[vMatch]), "control")
						}

					}
					f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row), fmt.Sprintf("=%s", formula))
				}
				row++
				continue
			}
			if v.AutoSummary != nil && *v.AutoSummary {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleCurrencyTotal)

				f.SetCellFormula(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("=SUM(C%d:C%d)", partRowStart, row-1))
				row++
				partRowStart = row
				continue
			}
			if v.ControlFormula != "" {

				// if v.FxSummary == "" {
				// 	row++
				// 	continue
				// }
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), stylingControl)
				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					//cari jml berdasarkan code
					if _, ok := rowCodePersediaan[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("C%d", rowCodePersediaan[vMatch]))
					}
					if _, ok := tbRowCode[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCode[vMatch]))
						f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCode[vMatch]), "control")
					}

				}
				f.SetCellFormula(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("=%s", formula))
			}
			row++
		}
		rowStart = row + 3
		row = rowStart
	}

	return f, nil
}

var rowCodeEm = make(map[string]int)

func (s *service) ExportEmployeeBenefit(ctx *abstraction.Context, f *excelize.File) (*excelize.File, error) {
	sheet := "EMPLOYEE_BENEFIT"
	f.NewSheet(sheet)

	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 4.30},
		{"COL": "B", "WIDTH": 3.71},
		{"COL": "C", "WIDTH": 8.43},
		{"COL": "D", "WIDTH": 8.43},
		{"COL": "E", "WIDTH": 8.43},
		{"COL": "F", "WIDTH": 21.71},
		{"COL": "G", "WIDTH": 19.14},
	}
	for _, v := range arrStyleColWidth {
		tmpColWidth := fmt.Sprintf("%f", v["WIDTH"])
		colWidth, _ := strconv.ParseFloat(tmpColWidth, 64)
		err := f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			return nil, err
		}
	}

	styleHeader, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#f8cbad"},
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}
	numberFormat := "#,##"
	styleCurrency, err := f.NewStyle(&excelize.Style{
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}
	stylingControl, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#FFFF00"},
		},
	})
	if err != nil {
		return nil, err
	}
	styleCurrencyTotal, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}

	styleLabelTotal, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return nil, err
	}

	formatterCode := []string{"EMPLOYEE-BENEFIT-ASUMSI", "EMPLOYEE-BENEFIT-REKONSILIASI", "EMPLOYEE-BENEFIT-RINCIAN-LAPORAN", "EMPLOYEE-BENEFIT-RINCIAN-EKUITAS", "EMPLOYEE-BENEFIT-MUTASI", "EMPLOYEE-BENEFIT-INFORMASI", "EMPLOYEE-BENEFIT-ANALISIS"}
	formatterTitle := []string{"Asumsi-asumsi yang digunakan:", "Rekonsiliasi jumlah liabilitas imbalan kerja karyawan pada laporan posisi keuangan adalah sebagai berikut:", "Rincian beban imbalan kerja karyawan yang diakui dalam laporan laba rugi dan penghasilan komprehensif lain adalah sebagai berikut:", "Rincian beban imbalan kerja karyawan yang diakui pada ekuitas dalam penghasilan komprehensif lain adalah sebagai berikut:", "Mutasi liabilitas imbalan kerja karyawan adalah sebagai berikut:", "Informasi historis dari nilai kini liabilitas imbalan pasti, nilai wajar aset program dan penyesuaian adalah sebagai berikut:", "Analisis sensitivitas dari perubahan asumsi-asumsi utama terhadap liabilitas imbalan kerja", ""}
	row, rowStart := 7, 7

	for i, formatter := range formatterCode {
		f.SetCellValue(sheet, fmt.Sprintf("B%d", (rowStart-3)), formatterTitle[i])
		f.MergeCell(sheet, fmt.Sprintf("B%d", (rowStart-2)), fmt.Sprintf("F%d", (rowStart-1)))
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", (rowStart-2)), fmt.Sprintf("G%d", (rowStart-1)), styleHeader)
		// f.SetCellStyle(sheet, fmt.Sprintf("A%d", (rowStart-3)), fmt.Sprintf("A%d", (rowStart-3)), styleLabel)

		f.SetCellValue(sheet, fmt.Sprintf("B%d", (rowStart-2)), "Description")
		// f.SetCellValue(sheet, fmt.Sprintf("G%d", (rowStart-2)), "period")
		// f.SetCellValue(sheet, fmt.Sprintf("G%d", (rowStart-1)), "employeeBenefit.Company.Name")
		f.SetCellFormula(sheet, fmt.Sprintf("G%d", (rowStart-1)), "=TRIAL_BALANCE!G6")
		// f.SetCellFormula(sheet, fmt.Sprintf("G%d", (rowStart-1)), "=TRIAL_BALANCE!G8")
		f.SetCellValue(sheet, fmt.Sprintf("G%d", (rowStart-2)), "31-Dec-21")

		var criteria dto.FormatterGetRequest
		criteria.FormatterFilterModel.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
		if err != nil {
			return nil, helper.ErrorHandler(err)
		}

		partRowStart := row
		for _, v := range data.FormatterDetail {
			if strings.Contains(strings.ToLower(v.Code), "blank") {
				row++
				continue
			}
			if rowCodeEm[v.Code] == 0 {
				rowCodeEm[v.Code] = row
			}
			rowKosong := 0.0
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
			// f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabel)
			f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styleCurrency)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), rowKosong)

			if v.ControlFormula != "" {
				if v.AutoSummary != nil && *v.AutoSummary {

					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
					f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styleCurrencyTotal)
					f.SetCellFormula(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("=SUM($G%d:$G%d)", partRowStart, row-1))
					partRowStart = row
				}
				if v.FxSummary == "" {
					row++
					continue
				}
				f.SetCellStyle(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), stylingControl)
				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					//cari jml berdasarkan code
					if _, ok := rowCodeEm[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("G%d", rowCodeEm[vMatch]))
					}
					if _, ok := tbRowCode[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCode[vMatch]))
						f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCode[vMatch]), "control")
					}

				}
				f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("=%s", formula))
				row++
				continue
			}
			if v.IsTotal != nil && *v.IsTotal {

				// if v.FxSummary == "" {
				// 	row++
				// 	continue
				// }

				arrChr := []string{"G"}

				if strings.ToUpper(v.Code) == "CONTROL_1" {
					arrChr = []string{"G"}
					f.SetCellStyle(sheet, fmt.Sprintf("G%d", row-16), fmt.Sprintf("G%d", row-16), styleCurrency)
					f.SetCellStyle(sheet, fmt.Sprintf("G%d", row-16), fmt.Sprintf("G%d", row-16), stylingControl)
				}
				if strings.ToUpper(v.Code) == "CONTROL_2" {
					arrChr = []string{"G"}
					f.SetCellStyle(sheet, fmt.Sprintf("G%d", row-28), fmt.Sprintf("G%d", row-28), styleCurrency)
					f.SetCellStyle(sheet, fmt.Sprintf("G%d", row-28), fmt.Sprintf("G%d", row-28), stylingControl)
				}
				for _, chr := range arrChr {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						//cari jml berdasarkan code
						if rowCodeEm[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCodeEm[vMatch]))
						}
						if _, ok := tbRowCode[vMatch]; ok {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCode[vMatch]))
							f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCode[vMatch]), "control")
						}

					}

					f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row), fmt.Sprintf("=%s", formula))
					if strings.ToUpper(v.Code) == "CONTROL_1" {
						f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row-16), fmt.Sprintf("=%s", formula))
					}
					if strings.ToUpper(v.Code) == "CONTROL_2" {
						f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row-28), fmt.Sprintf("=%s", formula))
					}
				}
				row++
				continue
			}
			if v.AutoSummary != nil && *v.AutoSummary {

				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styleCurrencyTotal)
				f.SetCellFormula(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("=SUM($G%d:$G%d)", partRowStart, row-1))
				row++
				partRowStart = row
				continue
			}

			row++
		}
		rowStart = row + 5
		row = rowStart
	}

	return f, nil
}
func (s *service) ExportInvestasiNonTbk(ctx *abstraction.Context, f *excelize.File) (*excelize.File, error) {
	sheet := "INVESTASI_NON_TBK"
	f.NewSheet(sheet)

	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 8.21},
		{"COL": "B", "WIDTH": 3.50},
		{"COL": "C", "WIDTH": 13.71},
		{"COL": "D", "WIDTH": 13.71},
		{"COL": "E", "WIDTH": 13.71},
		{"COL": "F", "WIDTH": 13.71},
		{"COL": "G", "WIDTH": 13.71},
		{"COL": "H", "WIDTH": 13.71},
		{"COL": "I", "WIDTH": 13.71},
	}
	for _, v := range arrStyleColWidth {
		tmpColWidth := fmt.Sprintf("%f", v["WIDTH"])
		colWidth, _ := strconv.ParseFloat(tmpColWidth, 64)
		err := f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("internal server error"))
		}
	}

	styleHeader, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold:   true,
			Family: "Arial",
			Size:   10,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#fac090"},
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("internal server error"))
	}
	numberFormat := "#,##"
	// styleCurrencyPercentage, err := f.NewStyle(&excelize.Style{
	// 	Border: []excelize.Border{
	// 		{Type: "top", Color: "000000", Style: 1},
	// 		{Type: "bottom", Color: "000000", Style: 1},
	// 		{Type: "left", Color: "000000", Style: 1},
	// 		{Type: "right", Color: "000000", Style: 1},
	// 	},
	// 	CustomNumFmt: &numberFormat,
	// 	Font: &excelize.Font{
	// 		Family: "Arial",
	// 		Size:   10,
	// 	},
	// })
	// if err != nil {
	// 	return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("internal server error"))
	// }
	styleCurrency, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &numberFormat,
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("internal server error"))
	}

	styleCurrencyAccounting, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		NumFmt: 41,
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("internal server error"))

	}
	styleLabel, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("internal server error"))
	}

	// allowed := helper.CompanyValidation(ctx.Auth.ID, investasiNonTbk.CompanyID)
	// if !allowed {
	// 	return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	// }

	// datePeriod, err := time.Parse(time.RFC3339, investasiNonTbk.Period)
	// if err != nil {
	// 	return nil, err
	// }

	row := 5

	f.SetCellValue(sheet, "B2", "Detail investasi anak usaha Non TBK")

	f.SetCellStyle(sheet, "B4", "J4", styleHeader)
	f.SetCellValue(sheet, "B4", "No")
	f.SetCellValue(sheet, "C4", "Code")
	f.SetCellValue(sheet, "D4", "Lembar saham dimiliki")
	f.SetCellValue(sheet, "E4", "Total lembar saham")
	f.SetCellValue(sheet, "F4", "% Ownership")
	f.SetCellValue(sheet, "G4", "Harga Par")
	f.SetCellValue(sheet, "H4", "Total harga Par")
	f.SetCellValue(sheet, "I4", "Harga beli")
	f.SetCellValue(sheet, "J4", "Total Harga beli")

	// criteriaDetail := model.InvestasiNonTbkDetailFilterModel{}

	// criteriaDetail.InvestasiNonTbkID = &investasiNonTbk.ID

	// paginationDetail := abstraction.Pagination{}
	// pagesize := 10000
	// paginationDetail.PageSize = &pagesize

	// detail, _, err := s.InvestasiNonTbkDetailRepository.Find(ctx, &criteriaDetail, &paginationDetail)
	// if err != nil {
	// 	return nil, helper.ErrorHandler(err)
	// }

	valueKosong := 0.0

	f.SetCellValue(sheet, fmt.Sprintf("D%d", row), valueKosong)
	f.SetCellValue(sheet, fmt.Sprintf("E%d", row), valueKosong)
	f.SetCellFormula(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("=IFERROR(D%d/E%d,%d)", row, row, 0))
	f.SetCellValue(sheet, fmt.Sprintf("G%d", row), valueKosong)
	f.SetCellValue(sheet, fmt.Sprintf("I%d", row), valueKosong)

	f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("=D%d*G%d", row, row))
	f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=D%d*I%d", row, row))

	f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabel)
	f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("E%d", row), styleCurrency)
	f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), styleCurrencyAccounting)
	f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styleCurrency)
	f.SetCellStyle(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), styleCurrency)
	f.SetCellStyle(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("J%d", row), styleCurrencyAccounting)
	f.SetCellStyle(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), styleCurrencyAccounting)

	row++

	return f, nil
}

func (s *service) ExportAje(ctx *abstraction.Context, filePath string, payload dto.AjeTemplateRequest) (*[]model.AdjustmentDetailEntity, error) {

	var AdjustmentDetail []model.AdjustmentDetailEntity
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	sheet := "AJE"
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	rows = rows[7:][:]
	for _, row := range rows {

		if len(row) == 0 {
			continue
		}
		if row[1] == "" {
			continue
		}

		if row[1] != "" && len(row) == 7 {
			if row[1] == "310401004" || row[1] == "310501002" || row[1] == "310502002" || row[1] == "310503002" || row[1] == "310402002" {
				return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+row[1])
			}
			coa := row[1]
			codeCoa := coa[:1]
			
			if codeCoa == "4"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "5"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "6"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "7"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "8"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "9"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			findByCodeCoa, err := s.CoaRepository.FindWithCode(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			findByCodeCoas, err := s.CoaRepository.FindWithCodes(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			if len(*findByCodeCoa) == 0 {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Tidak Terdaftar di Master Coa "+row[1], "Coa Tidak Terdaftar di Master Coa "+row[1])
			}
			if row[6] == "" {
				row[6] = strings.Replace(strings.ToUpper(row[6]), "", "0", -1)
			}

			BlnceDr, err := strconv.ParseFloat(row[6], 64)
			if err != nil {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Tolong Masukan Nominal yg Sesuai "+row[6], "Tolong Masukan Nominal yg Sesuai "+row[6])
			}
			BlnceCr := 0.0
			if err != nil {
				return nil, err
			}
			IncmDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmCr := 0.0
			if err != nil {
				return nil, err
			}
			data := model.AdjustmentDetailEntity{
			
				CoaCode:           row[1],
				ReffNumber:        &row[2],
				Description:       &findByCodeCoas.Name,
				BalanceSheetDr:    &BlnceDr,
				BalanceSheetCr:    &BlnceCr,
				IncomeStatementDr: &IncmDr,
				IncomeStatementCr: &IncmCr,
			}
			AdjustmentDetail = append(AdjustmentDetail, data)
		}
		if row[1] != "" && len(row) == 8 {
			if row[1] == "310401004" || row[1] == "310501002" || row[1] == "310502002" || row[1] == "310503002" || row[1] == "310402002" {
				return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+row[1])
			}
			coa := row[1]
			codeCoa := coa[:1]
			
			if codeCoa == "4"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "5"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "6"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "7"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "8"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "9"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			findByCodeCoa, err := s.CoaRepository.FindWithCode(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			findByCodeCoas, err := s.CoaRepository.FindWithCodes(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			if len(*findByCodeCoa) == 0 {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Tidak Terdaftar di Master Coa "+row[1], "Coa Tidak Terdaftar di Master Coa "+row[1])
			}
			if row[7] == "" {
				row[7] = strings.Replace(strings.ToUpper(row[7]), "", "0", -1)
			}

			BlnceCr, err := strconv.ParseFloat(row[7], 64)
			if err != nil {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Tolong Masukan Nominal yg Sesuai "+row[7], "Tolong Masukan Nominal yg Sesuai "+row[7])
			}
			BlnceDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmCr := 0.0
			if err != nil {
				return nil, err
			}
			data := model.AdjustmentDetailEntity{
				
				CoaCode:           row[1],
				ReffNumber:        &row[2],
				Description:       &findByCodeCoas.Name,
				BalanceSheetDr:    &BlnceDr,
				BalanceSheetCr:    &BlnceCr,
				IncomeStatementDr: &IncmDr,
				IncomeStatementCr: &IncmCr,
			}
			AdjustmentDetail = append(AdjustmentDetail, data)
		}
		if row[1] != "" && len(row) == 9 {
			if row[1] == "310401004" || row[1] == "310501002" || row[1] == "310502002" || row[1] == "310503002" || row[1] == "310402002" {
				return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+row[1])
			}
			coa := row[1]
			codeCoa := coa[:1]
			if codeCoa == "1"{
				if row[8] != ""  {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "2"{
				if row[8] != ""  {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "3"{
				if row[8] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "4"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "5"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "6"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "7"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "8"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "9"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			findByCodeCoa, err := s.CoaRepository.FindWithCode(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			findByCodeCoas, err := s.CoaRepository.FindWithCodes(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			if len(*findByCodeCoa) == 0 {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Tidak Terdaftar di Master Coa "+row[1], "Coa Tidak Terdaftar di Master Coa "+row[1])
			}
			if row[8] == "" {
				row[8] = strings.Replace(strings.ToUpper(row[8]), "", "0", -1)
			}

			IncmDr, err := strconv.ParseFloat(row[8], 64)
			if err != nil {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Tolong Masukan Nominal yg Sesuai "+row[8], "Tolong Masukan Nominal yg Sesuai "+row[8])
			}
			BlnceCr := 0.0
			if err != nil {
				return nil, err
			}
			BlnceDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmCr := 0.0
			if err != nil {
				return nil, err
			}
			data := model.AdjustmentDetailEntity{
				
				CoaCode:           row[1],
				ReffNumber:        &row[2],
				Description:       &findByCodeCoas.Name,
				BalanceSheetDr:    &BlnceDr,
				BalanceSheetCr:    &BlnceCr,
				IncomeStatementDr: &IncmDr,
				IncomeStatementCr: &IncmCr,
			}
			AdjustmentDetail = append(AdjustmentDetail, data)
		}
		if row[1] != "" && len(row) == 10 {
			if row[1] == "310401004" || row[1] == "310501002" || row[1] == "310502002" || row[1] == "310503002" || row[1] == "310402002" {
				return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+row[1])
			}
			coa := row[1]
			codeCoa := coa[:1]
			if codeCoa == "1" {
				if row[8] != "" || row[9] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "2" {
				if row[8] != "" || row[9] != ""  {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "3"{
				if row[8] != "" || row[9] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "4"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "5"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "6"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "7"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "8"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "9"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			findByCodeCoa, err := s.CoaRepository.FindWithCode(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			findByCodeCoas, err := s.CoaRepository.FindWithCodes(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			if len(*findByCodeCoa) == 0 {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Tidak Terdaftar di Master Coa "+row[1], "Coa Tidak Terdaftar di Master Coa "+row[1])
			}
			if row[9] == "" {
				row[9] = strings.Replace(strings.ToUpper(row[9]), "", "0", -1)
			}

			IncmCr, err := strconv.ParseFloat(row[9], 64)
			if err != nil {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Tolong Masukan Nominal yg Sesuai "+row[9], "Tolong Masukan Nominal yg Sesuai "+row[9])
			}
			BlnceCr := 0.0
			if err != nil {
				return nil, err
			}
			BlnceDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmDr := 0.0
			if err != nil {
				return nil, err
			}
			data := model.AdjustmentDetailEntity{
				CoaCode:           row[1],
				ReffNumber:        &row[2],
				Description:       &findByCodeCoas.Name,
				BalanceSheetDr:    &BlnceDr,
				BalanceSheetCr:    &BlnceCr,
				IncomeStatementDr: &IncmDr,
				IncomeStatementCr: &IncmCr,
			}
			AdjustmentDetail = append(AdjustmentDetail, data)
		}

	}
	return &AdjustmentDetail, nil
}
func (s *service) ExportJcte(ctx *abstraction.Context, filePath string, payload dto.AjeTemplateRequest) (*[]model.JcteDetailEntity, error) {

	var JcteDetail []model.JcteDetailEntity
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	sheet := "JCTE"
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	rows = rows[7:][:]
	for _, row := range rows {

		if len(row) == 0 {
			continue
		}
		if row[1] == "" {
			continue
		}

		if row[1] != "" && len(row) == 7 {
			if row[1] == "310401004" || row[1] == "310501002" || row[1] == "310502002" || row[1] == "310503002" || row[1] == "310402002" {
				return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+row[1])
			}
			coa := row[1]
			codeCoa := coa[:1]
			
			if codeCoa == "4"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "5"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "6"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "7"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "8"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "9"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			findByCodeCoa, err := s.CoaRepository.FindWithCode(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			findByCodeCoas, err := s.CoaRepository.FindWithCodes(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			if len(*findByCodeCoa) == 0 {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Tidak Terdaftar di Master Coa "+row[1], "Coa Tidak Terdaftar di Master Coa "+row[1])
			}
			if row[6] == "" {
				row[6] = strings.Replace(strings.ToUpper(row[6]), "", "0", -1)
			}

			BlnceDr, err := strconv.ParseFloat(row[6], 64)
			if err != nil {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Tolong Masukan Nominal yg Sesuai "+row[6], "Tolong Masukan Nominal yg Sesuai "+row[6])
			}
			BlnceCr := 0.0
			if err != nil {
				return nil, err
			}
			IncmDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmCr := 0.0
			if err != nil {
				return nil, err
			}
			data := model.JcteDetailEntity{
				// JcteID:      3,
				CoaCode:           row[1],
				ReffNumber:        &row[2],
				Description:       &findByCodeCoas.Name,
				BalanceSheetDr:    &BlnceDr,
				BalanceSheetCr:    &BlnceCr,
				IncomeStatementDr: &IncmDr,
				IncomeStatementCr: &IncmCr,
			}
			JcteDetail = append(JcteDetail, data)
		}
		if row[1] != "" && len(row) == 8 {
			if row[1] == "310401004" || row[1] == "310501002" || row[1] == "310502002" || row[1] == "310503002" || row[1] == "310402002" {
				return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+row[1])
			}
			coa := row[1]
			codeCoa := coa[:1]
			
			if codeCoa == "4"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "5"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "6"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "7"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "8"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "9"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			findByCodeCoa, err := s.CoaRepository.FindWithCode(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			findByCodeCoas, err := s.CoaRepository.FindWithCodes(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			if len(*findByCodeCoa) == 0 {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Tidak Terdaftar di Master Coa "+row[1], "Coa Tidak Terdaftar di Master Coa "+row[1])
			}
			if row[7] == "" {
				row[7] = strings.Replace(strings.ToUpper(row[7]), "", "0", -1)
			}

			BlnceCr, err := strconv.ParseFloat(row[7], 64)
			if err != nil {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Tolong Masukan Nominal yg Sesuai "+row[7], "Tolong Masukan Nominal yg Sesuai "+row[7])
			}
			BlnceDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmCr := 0.0
			if err != nil {
				return nil, err
			}
			data := model.JcteDetailEntity{
				// JcteID:      3,
				CoaCode:           row[1],
				ReffNumber:        &row[2],
				Description:       &findByCodeCoas.Name,
				BalanceSheetDr:    &BlnceDr,
				BalanceSheetCr:    &BlnceCr,
				IncomeStatementDr: &IncmDr,
				IncomeStatementCr: &IncmCr,
			}
			JcteDetail = append(JcteDetail, data)
		}
		if row[1] != "" && len(row) == 9 {
			if row[1] == "310401004" || row[1] == "310501002" || row[1] == "310502002" || row[1] == "310503002" || row[1] == "310402002" {
				return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+row[1])
			}
			coa := row[1]
			codeCoa := coa[:1]
			if codeCoa == "1"{
				if row[8] != ""  {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "2"{
				if row[8] != ""  {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "3"{
				if row[8] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "4"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "5"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "6"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "7"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "8"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "9"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			findByCodeCoa, err := s.CoaRepository.FindWithCode(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			findByCodeCoas, err := s.CoaRepository.FindWithCodes(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			if len(*findByCodeCoa) == 0 {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Tidak Terdaftar di Master Coa "+row[1], "Coa Tidak Terdaftar di Master Coa "+row[1])
			}
			if row[8] == "" {
				row[8] = strings.Replace(strings.ToUpper(row[8]), "", "0", -1)
			}

			IncmDr, err := strconv.ParseFloat(row[8], 64)
			if err != nil {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Tolong Masukan Nominal yg Sesuai "+row[8], "Tolong Masukan Nominal yg Sesuai "+row[8])
			}
			BlnceCr := 0.0
			if err != nil {
				return nil, err
			}
			BlnceDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmCr := 0.0
			if err != nil {
				return nil, err
			}
			data := model.JcteDetailEntity{
				// JcteID:      3,
				CoaCode:           row[1],
				ReffNumber:        &row[2],
				Description:       &findByCodeCoas.Name,
				BalanceSheetDr:    &BlnceDr,
				BalanceSheetCr:    &BlnceCr,
				IncomeStatementDr: &IncmDr,
				IncomeStatementCr: &IncmCr,
			}
			JcteDetail = append(JcteDetail, data)
		}
		if row[1] != "" && len(row) == 10 {
			if row[1] == "310401004" || row[1] == "310501002" || row[1] == "310502002" || row[1] == "310503002" || row[1] == "310402002" {
				return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+row[1])
			}
			coa := row[1]
			codeCoa := coa[:1]
			if codeCoa == "1" {
				if row[8] != "" || row[9] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "2" {
				if row[8] != "" || row[9] != ""  {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "3"{
				if row[8] != "" || row[9] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "4"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "5"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "6"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "7"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "8"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "9"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			findByCodeCoa, err := s.CoaRepository.FindWithCode(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			findByCodeCoas, err := s.CoaRepository.FindWithCodes(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			if len(*findByCodeCoa) == 0 {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Tidak Terdaftar di Master Coa "+row[1], "Coa Tidak Terdaftar di Master Coa "+row[1])
			}
			if row[9] == "" {
				row[9] = strings.Replace(strings.ToUpper(row[9]), "", "0", -1)
			}

			IncmCr, err := strconv.ParseFloat(row[9], 64)
			if err != nil {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Tolong Masukan Nominal yg Sesuai "+row[9], "Tolong Masukan Nominal yg Sesuai "+row[9])
			}
			BlnceCr := 0.0
			if err != nil {
				return nil, err
			}
			BlnceDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmDr := 0.0
			if err != nil {
				return nil, err
			}
			data := model.JcteDetailEntity{
				
				CoaCode:           row[1],
				ReffNumber:        &row[2],
				Description:       &findByCodeCoas.Name,
				BalanceSheetDr:    &BlnceDr,
				BalanceSheetCr:    &BlnceCr,
				IncomeStatementDr: &IncmDr,
				IncomeStatementCr: &IncmCr,
			}
			JcteDetail = append(JcteDetail, data)
		}

	}
	return &JcteDetail, nil
}
func (s *service) ExportJelim(ctx *abstraction.Context, filePath string, payload dto.AjeTemplateRequest) (*[]model.JelimDetailEntity, error) {

	var JelimDetail []model.JelimDetailEntity
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	sheet := "JELIM"
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	rows = rows[7:][:]
	for _, row := range rows {

		if len(row) == 0 {
			continue
		}
		if row[1] == "" {
			continue
		}

		if row[1] != "" && len(row) == 7 {
			if row[1] == "310401004" || row[1] == "310501002" || row[1] == "310502002" || row[1] == "310503002" || row[1] == "310402002" {
				return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+row[1])
			}
			coa := row[1]
			codeCoa := coa[:1]
			
			if codeCoa == "4"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "5"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "6"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "7"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "8"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "9"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			findByCodeCoa, err := s.CoaRepository.FindWithCode(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			findByCodeCoas, err := s.CoaRepository.FindWithCodes(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			if len(*findByCodeCoa) == 0 {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Tidak Terdaftar di Master Coa "+row[1], "Coa Tidak Terdaftar di Master Coa "+row[1])
			}
			if row[6] == "" {
				row[6] = strings.Replace(strings.ToUpper(row[6]), "", "0", -1)
			}

			BlnceDr, err := strconv.ParseFloat(row[6], 64)
			if err != nil {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Tolong Masukan Nominal yg Sesuai "+row[6], "Tolong Masukan Nominal yg Sesuai "+row[6])
			}
			BlnceCr := 0.0
			if err != nil {
				return nil, err
			}
			IncmDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmCr := 0.0
			if err != nil {
				return nil, err
			}
			data := model.JelimDetailEntity{
				// JelimID:      3,
				CoaCode:           row[1],
				ReffNumber:        &row[2],
				Description:       &findByCodeCoas.Name,
				BalanceSheetDr:    &BlnceDr,
				BalanceSheetCr:    &BlnceCr,
				IncomeStatementDr: &IncmDr,
				IncomeStatementCr: &IncmCr,
			}
			JelimDetail = append(JelimDetail, data)
		}
		if row[1] != "" && len(row) == 8 {
			if row[1] == "310401004" || row[1] == "310501002" || row[1] == "310502002" || row[1] == "310503002" || row[1] == "310402002" {
				return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+row[1])
			}
			coa := row[1]
			codeCoa := coa[:1]
			
			if codeCoa == "4"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "5"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "6"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "7"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "8"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "9"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			findByCodeCoa, err := s.CoaRepository.FindWithCode(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			findByCodeCoas, err := s.CoaRepository.FindWithCodes(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			if len(*findByCodeCoa) == 0 {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Tidak Terdaftar di Master Coa "+row[1], "Coa Tidak Terdaftar di Master Coa "+row[1])
			}
			if row[7] == "" {
				row[7] = strings.Replace(strings.ToUpper(row[7]), "", "0", -1)
			}

			BlnceCr, err := strconv.ParseFloat(row[7], 64)
			if err != nil {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Tolong Masukan Nominal yg Sesuai "+row[7], "Tolong Masukan Nominal yg Sesuai "+row[7])
			}
			BlnceDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmCr := 0.0
			if err != nil {
				return nil, err
			}
			data := model.JelimDetailEntity{
				// JelimID:      3,
				CoaCode:           row[1],
				ReffNumber:        &row[2],
				Description:       &findByCodeCoas.Name,
				BalanceSheetDr:    &BlnceDr,
				BalanceSheetCr:    &BlnceCr,
				IncomeStatementDr: &IncmDr,
				IncomeStatementCr: &IncmCr,
			}
			JelimDetail = append(JelimDetail, data)
		}
		if row[1] != "" && len(row) == 9 {
			if row[1] == "310401004" || row[1] == "310501002" || row[1] == "310502002" || row[1] == "310503002" || row[1] == "310402002" {
				return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+row[1])
			}
			coa := row[1]
			codeCoa := coa[:1]
			if codeCoa == "1"{
				if row[8] != ""  {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "2"{
				if row[8] != ""  {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "3"{
				if row[8] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "4"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "5"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "6"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "7"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "8"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "9"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			findByCodeCoa, err := s.CoaRepository.FindWithCode(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			findByCodeCoas, err := s.CoaRepository.FindWithCodes(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			if len(*findByCodeCoa) == 0 {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Tidak Terdaftar di Master Coa "+row[1], "Coa Tidak Terdaftar di Master Coa "+row[1])
			}
			if row[8] == "" {
				row[8] = strings.Replace(strings.ToUpper(row[8]), "", "0", -1)
			}

			IncmDr, err := strconv.ParseFloat(row[8], 64)
			if err != nil {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Tolong Masukan Nominal yg Sesuai "+row[8], "Tolong Masukan Nominal yg Sesuai "+row[8])
			}
			BlnceCr := 0.0
			if err != nil {
				return nil, err
			}
			BlnceDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmCr := 0.0
			if err != nil {
				return nil, err
			}
			data := model.JelimDetailEntity{
				// JelimID:      3,
				CoaCode:           row[1],
				ReffNumber:        &row[2],
				Description:       &findByCodeCoas.Name,
				BalanceSheetDr:    &BlnceDr,
				BalanceSheetCr:    &BlnceCr,
				IncomeStatementDr: &IncmDr,
				IncomeStatementCr: &IncmCr,
			}
			JelimDetail = append(JelimDetail, data)
		}
		if row[1] != "" && len(row) == 10 {
			if row[1] == "310401004" || row[1] == "310501002" || row[1] == "310502002" || row[1] == "310503002" || row[1] == "310402002" {
				return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+row[1])
			}
			coa := row[1]
			codeCoa := coa[:1]
			if codeCoa == "1" {
				if row[8] != "" || row[9] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "2" {
				if row[8] != "" || row[9] != ""  {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "3"{
				if row[8] != "" || row[9] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "4"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "5"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "6"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "7"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "8"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "9"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			findByCodeCoa, err := s.CoaRepository.FindWithCode(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			findByCodeCoas, err := s.CoaRepository.FindWithCodes(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			if len(*findByCodeCoa) == 0 {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Tidak Terdaftar di Master Coa "+row[1], "Coa Tidak Terdaftar di Master Coa "+row[1])
			}
			if row[9] == "" {
				row[9] = strings.Replace(strings.ToUpper(row[9]), "", "0", -1)
			}

			IncmCr, err := strconv.ParseFloat(row[9], 64)
			if err != nil {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Tolong Masukan Nominal yg Sesuai "+row[9], "Tolong Masukan Nominal yg Sesuai "+row[9])
			}
			BlnceCr := 0.0
			if err != nil {
				return nil, err
			}
			BlnceDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmDr := 0.0
			if err != nil {
				return nil, err
			}
			data := model.JelimDetailEntity{
				// JelimID:      	   3,
				CoaCode:           row[1],
				ReffNumber:        &row[2],
				Description:       &findByCodeCoas.Name,
				BalanceSheetDr:    &BlnceDr,
				BalanceSheetCr:    &BlnceCr,
				IncomeStatementDr: &IncmDr,
				IncomeStatementCr: &IncmCr,
			}
			JelimDetail = append(JelimDetail, data)
		}

	}
	return &JelimDetail, nil
}
func (s *service) ExportJpm(ctx *abstraction.Context, filePath string, payload dto.AjeTemplateRequest) (*[]model.JpmDetailEntity, error) {

	var JpmDetail []model.JpmDetailEntity
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	sheet := "JPM"
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	rows = rows[7:][:]
	for _, row := range rows {

		if len(row) == 0 {
			continue
		}
		if row[1] == "" {
			continue
		}

		if row[1] != "" && len(row) == 7 {
			if row[1] == "310401004" || row[1] == "310501002" || row[1] == "310502002" || row[1] == "310503002" || row[1] == "310402002" {
				return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+row[1])
			}
			coa := row[1]
			codeCoa := coa[:1]
			
			if codeCoa == "4"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "5"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "6"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "7"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "8"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "9"{
				if row[6] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			findByCodeCoa, err := s.CoaRepository.FindWithCode(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			findByCodeCoas, err := s.CoaRepository.FindWithCodes(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			if len(*findByCodeCoa) == 0 {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Tidak Terdaftar di Master Coa "+row[1], "Coa Tidak Terdaftar di Master Coa "+row[1])
			}
			if row[6] == "" {
				row[6] = strings.Replace(strings.ToUpper(row[6]), "", "0", -1)
			}

			BlnceDr, err := strconv.ParseFloat(row[6], 64)
			if err != nil {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Tolong Masukan Nominal yg Sesuai "+row[6], "Tolong Masukan Nominal yg Sesuai "+row[6])
			}
			BlnceCr := 0.0
			if err != nil {
				return nil, err
			}
			IncmDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmCr := 0.0
			if err != nil {
				return nil, err
			}
			data := model.JpmDetailEntity{
				
				CoaCode:           row[1],
				ReffNumber:        &row[2],
				Description:       &findByCodeCoas.Name,
				BalanceSheetDr:    &BlnceDr,
				BalanceSheetCr:    &BlnceCr,
				IncomeStatementDr: &IncmDr,
				IncomeStatementCr: &IncmCr,
			}
			JpmDetail = append(JpmDetail, data)
		}
		if row[1] != "" && len(row) == 8 {
			if row[1] == "310401004" || row[1] == "310501002" || row[1] == "310502002" || row[1] == "310503002" || row[1] == "310402002" {
				return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+row[1])
			}
			coa := row[1]
			codeCoa := coa[:1]
			
			if codeCoa == "4"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "5"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "6"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "7"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "8"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "9"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			findByCodeCoa, err := s.CoaRepository.FindWithCode(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			findByCodeCoas, err := s.CoaRepository.FindWithCodes(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			if len(*findByCodeCoa) == 0 {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Tidak Terdaftar di Master Coa "+row[1], "Coa Tidak Terdaftar di Master Coa "+row[1])
			}
			if row[7] == "" {
				row[7] = strings.Replace(strings.ToUpper(row[7]), "", "0", -1)
			}

			BlnceCr, err := strconv.ParseFloat(row[7], 64)
			if err != nil {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Tolong Masukan Nominal yg Sesuai "+row[7], "Tolong Masukan Nominal yg Sesuai "+row[7])
			}
			BlnceDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmCr := 0.0
			if err != nil {
				return nil, err
			}
			data := model.JpmDetailEntity{
				// JpmID:      3,
				CoaCode:           row[1],
				ReffNumber:        &row[2],
				Description:       &findByCodeCoas.Name,
				BalanceSheetDr:    &BlnceDr,
				BalanceSheetCr:    &BlnceCr,
				IncomeStatementDr: &IncmDr,
				IncomeStatementCr: &IncmCr,
			}
			JpmDetail = append(JpmDetail, data)
		}
		if row[1] != "" && len(row) == 9 {
			if row[1] == "310401004" || row[1] == "310501002" || row[1] == "310502002" || row[1] == "310503002" || row[1] == "310402002" {
				return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+row[1])
			}
			coa := row[1]
			codeCoa := coa[:1]
			if codeCoa == "1"{
				if row[8] != ""  {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "2"{
				if row[8] != ""  {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "3"{
				if row[8] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "4"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "5"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "6"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "7"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "8"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "9"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			findByCodeCoa, err := s.CoaRepository.FindWithCode(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			findByCodeCoas, err := s.CoaRepository.FindWithCodes(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			if len(*findByCodeCoa) == 0 {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Tidak Terdaftar di Master Coa "+row[1], "Coa Tidak Terdaftar di Master Coa "+row[1])
			}
			if row[8] == "" {
				row[8] = strings.Replace(strings.ToUpper(row[8]), "", "0", -1)
			}

			IncmDr, err := strconv.ParseFloat(row[8], 64)
			if err != nil {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Tolong Masukan Nominal yg Sesuai "+row[8], "Tolong Masukan Nominal yg Sesuai "+row[8])
			}
			BlnceCr := 0.0
			if err != nil {
				return nil, err
			}
			BlnceDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmCr := 0.0
			if err != nil {
				return nil, err
			}
			data := model.JpmDetailEntity{
				// JpmID:      3,
				CoaCode:           row[1],
				ReffNumber:        &row[2],
				Description:       &findByCodeCoas.Name,
				BalanceSheetDr:    &BlnceDr,
				BalanceSheetCr:    &BlnceCr,
				IncomeStatementDr: &IncmDr,
				IncomeStatementCr: &IncmCr,
			}
			JpmDetail = append(JpmDetail, data)
		}
		if row[1] != "" && len(row) == 10 {
			if row[1] == "310401004" || row[1] == "310501002" || row[1] == "310502002" || row[1] == "310503002" || row[1] == "310402002" {
				return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+row[1])
			}
			coa := row[1]
			codeCoa := coa[:1]
			if codeCoa == "1" {
				if row[8] != "" || row[9] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "2" {
				if row[8] != "" || row[9] != ""  {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "3"{
				if row[8] != "" || row[9] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Income Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "4"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "5"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "6"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "7"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "8"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			if codeCoa == "9"{
				if row[6] != "" || row[7] != "" {
					return nil,  response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Kolom Balance Tidak Dapat Diisi ")
				}
			}
			findByCodeCoa, err := s.CoaRepository.FindWithCode(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			findByCodeCoas, err := s.CoaRepository.FindWithCodes(ctx, &row[1])
			if err != nil {
				return nil, err
			}
			if len(*findByCodeCoa) == 0 {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Tidak Terdaftar di Master Coa "+row[1], "Coa Tidak Terdaftar di Master Coa "+row[1])
			}
			if row[9] == "" {
				row[9] = strings.Replace(strings.ToUpper(row[9]), "", "0", -1)
			}

			IncmCr, err := strconv.ParseFloat(row[9], 64)
			if err != nil {
				return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Tolong Masukan Nominal yg Sesuai "+row[9], "Tolong Masukan Nominal yg Sesuai "+row[9])
			}
			BlnceCr := 0.0
			if err != nil {
				return nil, err
			}
			BlnceDr := 0.0
			if err != nil {
				return nil, err
			}
			IncmDr := 0.0
			if err != nil {
				return nil, err
			}
			data := model.JpmDetailEntity{
				// JpmID:      	   3,
				CoaCode:           row[1],
				ReffNumber:        &row[2],
				Description:       &findByCodeCoas.Name,
				BalanceSheetDr:    &BlnceDr,
				BalanceSheetCr:    &BlnceCr,
				IncomeStatementDr: &IncmDr,
				IncomeStatementCr: &IncmCr,
			}
			JpmDetail = append(JpmDetail, data)
		}

	}
	return &JpmDetail, nil
}
