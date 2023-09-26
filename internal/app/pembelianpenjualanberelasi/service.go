package pembelianpenjualanberelasi

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
	"os"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type service struct {
	Repository                 repository.PembelianPenjualanBerelasi
	PPBerelasiDetailRepository repository.PembelianPenjualanBerelasiDetail
	ParameterRepository        repository.Parameter
	FormatterRepository        repository.Formatter
	CompanyRepository          repository.Company
	Db                         *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiGetRequest) (*dto.PembelianPenjualanBerelasiGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiGetByIDRequest) (*dto.PembelianPenjualanBerelasiGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiCreateRequest) (*dto.PembelianPenjualanBerelasiCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiUpdateRequest) (*dto.PembelianPenjualanBerelasiUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiDeleteRequest) (*dto.PembelianPenjualanBerelasiDeleteResponse, error)
	GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.PembelianPenjualanBerelasiRepository
	ppBerelasiDetailRepository := f.PembelianPenjualanBerelasiDetailRepository
	parameterRepository := f.ParameterRepository
	formatterRepository := f.FormatterRepository
	companyRepository := f.CompanyRepository
	db := f.Db
	return &service{
		Repository:                 repository,
		PPBerelasiDetailRepository: ppBerelasiDetailRepository,
		ParameterRepository:        parameterRepository,
		FormatterRepository:        formatterRepository,
		CompanyRepository:          companyRepository,
		Db:                         db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiGetRequest) (*dto.PembelianPenjualanBerelasiGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.PembelianPenjualanBerelasiFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.PembelianPenjualanBerelasiGetResponse{}, helper.ErrorHandler(err)
	}

	result := &dto.PembelianPenjualanBerelasiGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiGetByIDRequest) (*dto.PembelianPenjualanBerelasiGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.PembelianPenjualanBerelasiGetByIDResponse{}, helper.ErrorHandler(err)
	}
	allowed := helper.CompanyValidation(ctx.Auth.ID, data.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}
	result := &dto.PembelianPenjualanBerelasiGetByIDResponse{
		PembelianPenjualanBerelasiEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiCreateRequest) (*dto.PembelianPenjualanBerelasiCreateResponse, error) {
	var data model.PembelianPenjualanBerelasiEntityModel

	allowed := helper.CompanyValidation(ctx.Auth.ID, payload.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.PembelianPenjualanBerelasiEntity = payload.PembelianPenjualanBerelasiEntity
		data.PembelianPenjualanBerelasiEntity.Status = 1

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.PembelianPenjualanBerelasiCreateResponse{}, err
	}

	result := &dto.PembelianPenjualanBerelasiCreateResponse{
		PembelianPenjualanBerelasiEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiUpdateRequest) (*dto.PembelianPenjualanBerelasiUpdateResponse, error) {
	var data model.PembelianPenjualanBerelasiEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if existing.Status != 1 {
			return response.ErrorBuilder(&response.ErrorConstant.DataValidated, errors.New("Cannot Update Data"))
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, existing.CompanyID)
		allowed2 := helper.CompanyValidation(ctx.Auth.ID, payload.CompanyID)
		if !allowed || !allowed2 {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
		}

		data.Context = ctx
		data.PembelianPenjualanBerelasiEntity = payload.PembelianPenjualanBerelasiEntity

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.PembelianPenjualanBerelasiUpdateResponse{}, err
	}
	result := &dto.PembelianPenjualanBerelasiUpdateResponse{
		PembelianPenjualanBerelasiEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiDeleteRequest) (*dto.PembelianPenjualanBerelasiDeleteResponse, error) {
	var data model.PembelianPenjualanBerelasiEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, existing.CompanyID)
		if !allowed {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
		}

		data.Context = ctx
		result, err := s.Repository.Delete(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.PembelianPenjualanBerelasiDeleteResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.PembelianPenjualanBerelasiDeleteResponse{
		// PembelianPenjualanBerelasiEntityModel: data,
	}
	return result, nil
}

func (s *service) Export(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiExportRequest) (*dto.PembelianPenjualanBerelasiExportResponse, error) {
	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	err := f.SetColWidth(sheet, "D", "D", 66)
	if err != nil {
		return &dto.PembelianPenjualanBerelasiExportResponse{}, err
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
		return &dto.PembelianPenjualanBerelasiExportResponse{}, err
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
		return &dto.PembelianPenjualanBerelasiExportResponse{}, err
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
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return &dto.PembelianPenjualanBerelasiExportResponse{}, err
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
		return &dto.PembelianPenjualanBerelasiExportResponse{}, err
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
		return &dto.PembelianPenjualanBerelasiExportResponse{}, err
	}

	f.SetCellStyle(sheet, "B4", "F4", styleHeader)
	f.SetCellValue(sheet, "C2", "List pembelian dan penjualan berelasi")
	f.SetCellValue(sheet, "B4", "NO")
	f.SetCellValue(sheet, "C4", "CODE")
	f.SetCellValue(sheet, "D4", "COMPANY")
	f.SetCellValue(sheet, "E4", "PEMBELIAN")
	f.SetCellValue(sheet, "F4", "PENJUALAN")

	pembelianpenjualanberelasi, err := s.Repository.FindByID(ctx, &payload.PembelianPenjualanBerelasiID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, pembelianpenjualanberelasi.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	criteria := model.PembelianPenjualanBerelasiFilterModel{}
	criteria.CompanyID = &pembelianpenjualanberelasi.CompanyID
	criteria.Period = &pembelianpenjualanberelasi.Period
	criteria.Versions = &pembelianpenjualanberelasi.Versions
	data, err := s.Repository.Export(ctx, &criteria)
	if err != nil {
		return nil, err
	}

	row, rowStart := 5, 5
	for i, v := range data.PembelianPenjualanBerelasiDetail {
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), (i + 1))
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Company.Code)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), v.Company.Name)
		if v.BoughtAmount != nil && *v.BoughtAmount != 0 {
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *v.BoughtAmount)
		}
		if v.SalesAmount != nil && *v.SalesAmount != 0 {
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), *v.SalesAmount)
		}
		row++
	}

	f.SetCellStyle(sheet, "B5", fmt.Sprintf("D%d", row), styleLabel)
	f.SetCellStyle(sheet, "E5", fmt.Sprintf("F%d", row), styleCurrency)
	f.SetCellStyle(sheet, fmt.Sprintf("B%d", row+1), fmt.Sprintf("D%d", row+1), styleLabelTotal)
	f.SetCellStyle(sheet, fmt.Sprintf("E%d", row+1), fmt.Sprintf("F%d", row+1), styleCurrencyTotal)

	f.SetCellValue(sheet, fmt.Sprintf("D%d", row+1), "Total")
	f.SetCellFormula(sheet, fmt.Sprintf("E%d", row+1), fmt.Sprintf("=SUM(E%d:E%d)", rowStart, row))
	f.SetCellFormula(sheet, fmt.Sprintf("F%d", row+1), fmt.Sprintf("=SUM(F%d:F%d)", rowStart, row))
	tmpFolder := fmt.Sprintf("assets/%d", ctx.Auth.ID)
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
	datePeriod, err := time.Parse(time.RFC3339, pembelianpenjualanberelasi.Period)
	if err != nil {
		return nil, err
	}
	period := datePeriod.Format("2006-01-02")
	fileName := fmt.Sprintf("MutasiRua_%s.xlsx", period)
	fileLoc := fmt.Sprintf("assets/%d/%s", ctx.Auth.ID, fileName)
	err = f.SaveAs(fileLoc)
	if err != nil {
		return nil, err
	}

	result := &dto.PembelianPenjualanBerelasiExportResponse{
		FileName: fileName,
		Path:     fileLoc,
	}
	return result, nil
}

func (s *service) GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error) {
	filter := model.PembelianPenjualanBerelasiFilterModel{
		CompanyCustomFilter: model.CompanyCustomFilter{
			CompanyID:          payload.CompanyID,
			ArrCompanyID:       payload.ArrCompanyID,
			ArrCompanyString:   payload.ArrCompanyString,
			ArrCompanyOperator: payload.ArrCompanyOperator,
		},
	}
	filter.Status = payload.Status
	filter.Period = payload.Period
	data, err := s.Repository.GetVersion(ctx, &filter)
	if err != nil {
		return &dto.GetVersionResponse{}, err
	}
	result := &dto.GetVersionResponse{
		Data: *data,
	}
	return result, nil
}
