package imports

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"math"
	"mcash-finance-console-core/configs"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/pkg/kafka"
	"mcash-finance-console-core/pkg/util/response"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

type handler struct {
	service *service
}

var err error

func NewHandler(f *factory.Factory) *handler {
	service := NewService(f)
	return &handler{service}
}

func (h *handler) UploadTemplate(c echo.Context) error {

	cc := c.(*abstraction.Context)
	payload := new(dto.ImportReUploadTemplateRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}
	if payload.Template == "allworksheet" {
		file, err := c.FormFile(payload.Template)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
		}
		start := time.Now()
		src, err := file.Open()
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
		}
		defer src.Close()
		duration := time.Since(start)
		fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")

		filepath := fmt.Sprintf("templates/%s", file.Filename)

		TrialBalance := path.Join(configs.App().StoragePath(), filepath)

		dst, err := os.Create(TrialBalance)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
	}
	

	return response.SuccessResponse("Upload Template Sukses").Send(cc)
}

func (h *handler) UploadJurnal(c echo.Context) error {

	cc := c.(*abstraction.Context)
	payload := new(dto.AjeTemplateRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}
	file, err := c.FormFile("JURNAL")
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	start := time.Now()
	src, err := file.Open()
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	defer src.Close()
	duration := time.Since(start)
	fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")

	filepath := fmt.Sprintf("templates/%s", file.Filename)

	Aje := path.Join(configs.App().StoragePath(), filepath)

	dst, err := os.Create(Aje)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	
	if payload.Jurnal == "AJE" {
		AJE, err := h.service.ExportAje(cc, filepath, *payload)
		if err != nil {
			return response.ErrorResponse(err).Send(c)
		}
		return response.SuccessResponse(AJE).Send(cc)
	}
	if payload.Jurnal == "JPM" {
		JPM, err := h.service.ExportJpm(cc, filepath, *payload)
		if err != nil {
			return response.ErrorResponse(err).Send(c)
		}
		return response.SuccessResponse(JPM).Send(cc)
	}
	if payload.Jurnal == "JCTE" {
		JCTE, err := h.service.ExportJcte(cc, filepath, *payload)
		if err != nil {
			return response.ErrorResponse(err).Send(c)
		}
		return response.SuccessResponse(JCTE).Send(cc)
	}
	if payload.Jurnal == "JELIM" {
		JELIM, err := h.service.ExportJelim(cc, filepath, *payload)
		if err != nil {
			return response.ErrorResponse(err).Send(c)
		}
		return response.SuccessResponse(JELIM).Send(cc)
	}
	return response.SuccessResponse("Sukses").Send(cc)
}


// Bulk Upload
// @Summary Bulk Upload
// @Description Bulk Upload
// @Tags Bulk Upload
// @Accept mpfd
// @Produce application/json
// @Security BearerAuth
// @param company_id formData int true "company_id formdata"
// @param period formData string true "period formdata"
// @param trial_balance formData file true "file formdata"
// @param mutasi_fa formData file true "file formdata"
// @param aging_utang_piutang formData file true "file formdata"
// @param investasi_tbk formData file true "file formdata"
// @param investasi_non_tbk formData file true "file formdata"
// @param mutasi_dta formData file true "file formdata"
// @param mutasi_ia formData file true "file formdata"
// @param mutasi_rua formData file true "file formdata"
// @param mutasi_persediaan formData file true "file formdata"
// @param pembelian_penjualan_berelasi formData file true "file formdata"
// @param employee_benefit formData file true "file formdata"
// @Success 200 {object} response.successResponse
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /import/importasync [post]
func (h *handler) ImportAsync(c echo.Context) error {

	cc := c.(*abstraction.Context)
	payload := new(dto.ImportedWorksheetCreateRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	
	importedworksheet, err := h.service.Create(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	//Trial Balance
	file, err := c.FormFile("all_worksheet")
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	start := time.Now()
	src, err := file.Open()
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	defer src.Close()
	duration := time.Since(start)
	fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")

	filepath := fmt.Sprintf("uploaded/%d-%d-%s-%s", importedworksheet.CompanyID, importedworksheet.Versions, importedworksheet.Period, file.Filename)

	TrialBalance := path.Join(configs.App().StoragePath(), filepath)
	FNTrialBalance := file.Filename
	dst, err := os.Create(TrialBalance)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	map1 := kafka.JsonDataImport{
		TrialBalance:        TrialBalance,
		FNTrialBalance:      FNTrialBalance,
		CompanyID:           payload.CompanyID,
		UserID:              cc.Auth.ID,
		Version:             importedworksheet.Versions,
		ImportedWorkSheetID: importedworksheet.ID,
		Period:              importedworksheet.Period,
	}
	jsonStr, err := json.Marshal(map1)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	kafka.NewService("IMPORT").SendMessage("IMPORT", string(jsonStr))

	return response.SuccessResponse(importedworksheet).Send(cc)
}
func (h *handler) Import(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.ImportReUploadRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}
	importedworksheet, err := h.service.FindByID(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	importedworksheetDetailsucces, err := h.service.FindByIDDetail(cc, &importedworksheet.ID)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	importedworksheetDetailfailed, err := h.service.FindByIDDetailS(cc, &importedworksheet.ID)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	Map := kafka.JsonDataReUpload{}
	Map.CompanyID = importedworksheet.CompanyID
	Map.UserID = cc.Auth.ID
	Map.Version = importedworksheet.Versions
	Map.ImportedWorkSheetID = importedworksheet.ID

	datasucces := importedworksheetDetailsucces.Data
	datafailed := importedworksheetDetailfailed.Data
	reupload := "reupload"
	if len(datasucces) == 11 || len(datasucces) < 11 {
		for _, iw := range importedworksheet.ImportedWorksheetDetail {
			file, err := c.FormFile("ALL-WORKSHEET")
			if err != nil {
				continue
			}

			start := time.Now()
			src, err := file.Open()
			if err != nil {
				return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
			}
			defer src.Close()
			duration := time.Since(start)
			fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")

			filepath := fmt.Sprintf("uploaded/%s-%d-%d-%s", reupload, importedworksheet.CompanyID, importedworksheet.Versions, file.Filename)

			AllWorksheet := path.Join(configs.App().StoragePath(), filepath)
			FNAllWorksheet := file.Filename
			dst, err := os.Create(AllWorksheet)
			if err != nil {
				return err
			}
			defer dst.Close()

			if _, err = io.Copy(dst, src); err != nil {
				return err
			}
			datePeriod, err := time.Parse(time.RFC3339, importedworksheet.Period)
			if err != nil {
				return err
			}
			period := datePeriod.Format("2006-01-02")
			criteriaTB := model.TrialBalanceFilterModel{}
			criteriaTB.Period = &period
			criteriaTB.CompanyID = &importedworksheet.CompanyID
			criteriaTB.Versions = &importedworksheet.Versions
			tbId, err := h.service.FindByVCTrialBalance(cc, &criteriaTB)
			if err != nil {
				return response.ErrorResponse(err).Send(c)
			}
			criteriaAupId := model.AgingUtangPiutangFilterModel{}
			criteriaAupId.Period = &period
			criteriaAupId.CompanyID = &importedworksheet.CompanyID
			criteriaAupId.Versions = &importedworksheet.Versions
			AupId, err := h.service.FindByVCAgingUtangPiutang(cc, &criteriaAupId)
			if err != nil {
				return response.ErrorResponse(err).Send(c)
			}
			criteriaItId := model.InvestasiTbkFilterModel{}
			criteriaItId.Period = &period
			criteriaItId.CompanyID = &importedworksheet.CompanyID
			criteriaItId.Versions = &importedworksheet.Versions
			ItId, err := h.service.FindByVCInvestasiTbk(cc, &criteriaItId)
			if err != nil {
				return response.ErrorResponse(err).Send(c)
			}
			criteriaIntId := model.InvestasiNonTbkFilterModel{}
			criteriaIntId.Period = &period
			criteriaIntId.CompanyID = &importedworksheet.CompanyID
			criteriaIntId.Versions = &importedworksheet.Versions
			IntId, err := h.service.FindByVCInvestasiNonTbk(cc, &criteriaIntId)
			if err != nil {
				return response.ErrorResponse(err).Send(c)
			}
			criteriaMfaId := model.MutasiFaFilterModel{}
			criteriaMfaId.Period = &period
			criteriaMfaId.CompanyID = &importedworksheet.CompanyID
			criteriaMfaId.Versions = &importedworksheet.Versions
			MfaId, err := h.service.FindByVCMutasiFa(cc, &criteriaMfaId)
			if err != nil {
				return response.ErrorResponse(err).Send(c)
			}
			criteriaMdtaId := model.MutasiDtaFilterModel{}
			criteriaMdtaId.Period = &period
			criteriaMdtaId.CompanyID = &importedworksheet.CompanyID
			criteriaMdtaId.Versions = &importedworksheet.Versions
			MdtaId, err := h.service.FindByVCMutasiDta(cc, &criteriaMdtaId)
			if err != nil {
				return response.ErrorResponse(err).Send(c)
			}
			criteriaMiaId := model.MutasiIaFilterModel{}
			criteriaMiaId.Period = &period
			criteriaMiaId.CompanyID = &importedworksheet.CompanyID
			criteriaMiaId.Versions = &importedworksheet.Versions
			MiaId, err := h.service.FindByVCMutasiIa(cc, &criteriaMiaId)
			if err != nil {
				return response.ErrorResponse(err).Send(c)
			}
			criteriaMruaId := model.MutasiRuaFilterModel{}
			criteriaMruaId.Period = &period
			criteriaMruaId.CompanyID = &importedworksheet.CompanyID
			criteriaMruaId.Versions = &importedworksheet.Versions
			MruaId, err := h.service.FindByVCMutasiRua(cc, &criteriaMruaId)
			if err != nil {
				return response.ErrorResponse(err).Send(c)
			}
			criteriaMpId := model.MutasiPersediaanFilterModel{}
			criteriaMpId.Period = &period
			criteriaMpId.CompanyID = &importedworksheet.CompanyID
			criteriaMpId.Versions = &importedworksheet.Versions
			MpId, err := h.service.FindByVCMutasiPersediaan(cc, &criteriaMpId)
			if err != nil {
				return response.ErrorResponse(err).Send(c)
			}
			criteriaPpb := model.PembelianPenjualanBerelasiFilterModel{}
			criteriaPpb.Period = &period
			criteriaPpb.CompanyID = &importedworksheet.CompanyID
			criteriaPpb.Versions = &importedworksheet.Versions
			PpbId, err := h.service.FindByVCPembelianPenjualanBerelasi(cc, &criteriaPpb)
			if err != nil {
				return response.ErrorResponse(err).Send(c)
			}

			criteriaEmpId := model.EmployeeBenefitFilterModel{}
			criteriaEmpId.Period = &period
			criteriaEmpId.CompanyID = &importedworksheet.CompanyID
			criteriaEmpId.Versions = &importedworksheet.Versions
			EmpId, err := h.service.FindByVCEmployeeBenefit(cc, &criteriaEmpId)
			if err != nil {
				return response.ErrorResponse(err).Send(c)
			}
			Map.AllWorksheet = AllWorksheet

			Map.IDTrialBalance = tbId.ID
			Map.IDAgingUtangPiutang = AupId.ID
			Map.IDMutasiDta = MdtaId.ID
			Map.IDMutasiFA = MfaId.ID
			Map.IDMutasiIa = MiaId.ID
			Map.IDMutasiPersediaan = MpId.ID
			Map.IDMutasiRua = MruaId.ID
			Map.IDPembelianPenjualanBerelasi = PpbId.ID
			Map.IDInvestasiTbk = ItId.ID
			Map.IDInvestasiNonTbk = IntId.ID
			Map.IDEmployeeBenefit = EmpId.ID

			Map.FNAllWorksheet = FNAllWorksheet

			if (iw.Name) == "Trial Balance" {
				Map.IDWorksheetDetailTrialBalance = iw.ID
			}
			if (iw.Name) == "Aging Utang Piutang" {
				Map.IDWorksheetDetailAgingUtangPiutang = iw.ID
			}
			if (iw.Name) == "Mutasi DTA" {
				Map.IDWorksheetDetailMutasiDta = iw.ID
			}
			if (iw.Name) == "Mutasi FA" {
				Map.IDWorksheetDetailMutasiFA = iw.ID
			}
			if (iw.Name) == "Mutasi IA" {
				Map.IDWorksheetDetailMutasiIa = iw.ID
			}
			if (iw.Name) == "Mutasi Persediaan" {
				Map.IDWorksheetDetailMutasiPersediaan = iw.ID
			}
			if (iw.Name) == "Mutasi RUA" {
				Map.IDWorksheetDetailMutasiRua = iw.ID
			}
			if (iw.Name) == "Investasi TBK" {
				Map.IDWorksheetDetailInvestasiTbk = iw.ID
			}
			if (iw.Name) == "Investasi Non TBK" {
				Map.IDWorksheetDetailInvestasiNonTbk = iw.ID
			}
			if (iw.Name) == "Pembelian & Penjualan Berelasi" {
				Map.IDWorksheetDetailPembelianPenjualanBerelasi = iw.ID
			}
			if (iw.Name) == "Employee Benefit" {
				Map.IDWorksheetDetailEmployeeBenefit = iw.ID
			}

		}
	}

	if len(datasucces) < 11 {
		for _, v := range datafailed {
			ImportedWorksheetDetail := model.ImportedWorksheetDetailEntityModel{
				ImportedWorksheetDetailEntity: v.ImportedWorksheetDetailEntity,
			}
			if (ImportedWorksheetDetail.Name) == "Trial Balance" {

				file, err := c.FormFile(ImportedWorksheetDetail.Code)
				if err != nil {
					continue
				}

				start := time.Now()
				src, err := file.Open()
				if err != nil {
					return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
				}
				defer src.Close()
				duration := time.Since(start)
				fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")

				filepath := fmt.Sprintf("uploaded/%s-%d-%d-%s", reupload, importedworksheet.CompanyID, importedworksheet.Versions, file.Filename)

				TrialBalance := path.Join(configs.App().StoragePath(), filepath)
				FNTrialBalance := file.Filename
				dst, err := os.Create(TrialBalance)
				if err != nil {
					return err
				}
				defer dst.Close()

				if _, err = io.Copy(dst, src); err != nil {
					return err
				}
				datePeriod, err := time.Parse(time.RFC3339, importedworksheet.Period)
				if err != nil {
					return err
				}
				period := datePeriod.Format("2006-01-02")
				criteriaTB := model.TrialBalanceFilterModel{}
				criteriaTB.Period = &period
				criteriaTB.CompanyID = &importedworksheet.CompanyID
				criteriaTB.Versions = &importedworksheet.Versions
				tbId, err := h.service.FindByVCTrialBalance(cc, &criteriaTB)
				if err != nil {
					return response.ErrorResponse(err).Send(c)
				}
				Map.TrialBalance = TrialBalance
				Map.IDTrialBalance = tbId.ID
				Map.IDWorksheetDetailTrialBalance = v.ID
				Map.FNTrialBalance = FNTrialBalance
			}
			if (ImportedWorksheetDetail.Name) == "Aging Utang Piutang" {

				file, err := c.FormFile(ImportedWorksheetDetail.Code)
				if err != nil {
					continue
				}

				start := time.Now()
				src, err := file.Open()
				if err != nil {
					return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
				}
				defer src.Close()
				duration := time.Since(start)
				fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")

				filepath := fmt.Sprintf("uploaded/%s-%d-%d-%s", reupload, importedworksheet.CompanyID, importedworksheet.Versions, file.Filename)

				AgingUtangPiutang := path.Join(configs.App().StoragePath(), filepath)

				dst, err := os.Create(AgingUtangPiutang)
				if err != nil {
					return err
				}
				defer dst.Close()

				if _, err = io.Copy(dst, src); err != nil {
					return err
				}
				datePeriod, err := time.Parse(time.RFC3339, importedworksheet.Period)
				if err != nil {
					return err
				}
				period := datePeriod.Format("2006-01-02")
				criteriaTB := model.AgingUtangPiutangFilterModel{}
				criteriaTB.Period = &period
				criteriaTB.CompanyID = &importedworksheet.CompanyID
				criteriaTB.Versions = &importedworksheet.Versions
				AupId, err := h.service.FindByVCAgingUtangPiutang(cc, &criteriaTB)
				if err != nil {
					return response.ErrorResponse(err).Send(c)
				}
				Map.AgingUtangPiutang = AgingUtangPiutang
				Map.IDAgingUtangPiutang = AupId.ID
				Map.IDWorksheetDetailAgingUtangPiutang = v.ID
			}
			if (ImportedWorksheetDetail.Name) == "Investasi TBK" {

				file, err := c.FormFile(ImportedWorksheetDetail.Code)
				if err != nil {
					continue
				}

				start := time.Now()
				src, err := file.Open()
				if err != nil {
					return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
				}
				defer src.Close()
				duration := time.Since(start)
				fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")

				filepath := fmt.Sprintf("uploaded/%s-%d-%d-%s", reupload, importedworksheet.CompanyID, importedworksheet.Versions, file.Filename)

				InvestasiTbk := path.Join(configs.App().StoragePath(), filepath)

				dst, err := os.Create(InvestasiTbk)
				if err != nil {
					return err
				}
				defer dst.Close()

				if _, err = io.Copy(dst, src); err != nil {
					return err
				}
				datePeriod, err := time.Parse(time.RFC3339, importedworksheet.Period)
				if err != nil {
					return err
				}
				period := datePeriod.Format("2006-01-02")
				criteriaTB := model.InvestasiTbkFilterModel{}
				criteriaTB.Period = &period
				criteriaTB.CompanyID = &importedworksheet.CompanyID
				criteriaTB.Versions = &importedworksheet.Versions
				ItId, err := h.service.FindByVCInvestasiTbk(cc, &criteriaTB)
				if err != nil {
					return response.ErrorResponse(err).Send(c)
				}
				Map.InvestasiTbk = InvestasiTbk
				Map.IDInvestasiTbk = ItId.ID
				Map.IDWorksheetDetailInvestasiTbk = v.ID
			}
			if (ImportedWorksheetDetail.Name) == "Investasi Non TBK" {

				file, err := c.FormFile(ImportedWorksheetDetail.Code)
				if err != nil {
					continue
				}

				start := time.Now()
				src, err := file.Open()
				if err != nil {
					return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
				}
				defer src.Close()
				duration := time.Since(start)
				fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")

				filepath := fmt.Sprintf("uploaded/%s-%d-%d-%s", reupload, importedworksheet.CompanyID, importedworksheet.Versions, file.Filename)

				InvestasiNonTbk := path.Join(configs.App().StoragePath(), filepath)

				dst, err := os.Create(InvestasiNonTbk)
				if err != nil {
					return err
				}
				defer dst.Close()

				if _, err = io.Copy(dst, src); err != nil {
					return err
				}
				datePeriod, err := time.Parse(time.RFC3339, importedworksheet.Period)
				if err != nil {
					return err
				}
				period := datePeriod.Format("2006-01-02")
				criteriaTB := model.InvestasiNonTbkFilterModel{}
				criteriaTB.Period = &period
				criteriaTB.CompanyID = &importedworksheet.CompanyID
				criteriaTB.Versions = &importedworksheet.Versions
				IntId, err := h.service.FindByVCInvestasiNonTbk(cc, &criteriaTB)
				if err != nil {
					return response.ErrorResponse(err).Send(c)
				}
				Map.InvestasiNonTbk = InvestasiNonTbk
				Map.IDInvestasiNonTbk = IntId.ID
				Map.IDWorksheetDetailInvestasiNonTbk = v.ID
			}
			if (ImportedWorksheetDetail.Name) == "Mutasi FA" {

				file, err := c.FormFile(ImportedWorksheetDetail.Code)
				if err != nil {
					continue
				}

				start := time.Now()
				src, err := file.Open()
				if err != nil {
					return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
				}
				defer src.Close()
				duration := time.Since(start)
				fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")

				filepath := fmt.Sprintf("uploaded/%s-%d-%d-%s", reupload, importedworksheet.CompanyID, importedworksheet.Versions, file.Filename)

				MutasiFA := path.Join(configs.App().StoragePath(), filepath)

				dst, err := os.Create(MutasiFA)
				if err != nil {
					return err
				}
				defer dst.Close()

				if _, err = io.Copy(dst, src); err != nil {
					return err
				}
				datePeriod, err := time.Parse(time.RFC3339, importedworksheet.Period)
				if err != nil {
					return err
				}
				period := datePeriod.Format("2006-01-02")
				criteriaTB := model.MutasiFaFilterModel{}
				criteriaTB.Period = &period
				criteriaTB.CompanyID = &importedworksheet.CompanyID
				criteriaTB.Versions = &importedworksheet.Versions
				MfaId, err := h.service.FindByVCMutasiFa(cc, &criteriaTB)
				if err != nil {
					return response.ErrorResponse(err).Send(c)
				}
				Map.MutasiFA = MutasiFA
				Map.IDMutasiFA = MfaId.ID
				Map.IDWorksheetDetailMutasiFA = v.ID
			}
			if (ImportedWorksheetDetail.Name) == "Mutasi DTA" {

				file, err := c.FormFile(ImportedWorksheetDetail.Code)
				if err != nil {
					continue
				}

				start := time.Now()
				src, err := file.Open()
				if err != nil {
					return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
				}
				defer src.Close()
				duration := time.Since(start)
				fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")

				filepath := fmt.Sprintf("uploaded/%s-%d-%d-%s", reupload, importedworksheet.CompanyID, importedworksheet.Versions, file.Filename)

				MutasiDta := path.Join(configs.App().StoragePath(), filepath)

				dst, err := os.Create(MutasiDta)
				if err != nil {
					return err
				}
				defer dst.Close()

				if _, err = io.Copy(dst, src); err != nil {
					return err
				}
				datePeriod, err := time.Parse(time.RFC3339, importedworksheet.Period)
				if err != nil {
					return err
				}
				period := datePeriod.Format("2006-01-02")
				criteriaTB := model.MutasiDtaFilterModel{}
				criteriaTB.Period = &period
				criteriaTB.CompanyID = &importedworksheet.CompanyID
				criteriaTB.Versions = &importedworksheet.Versions
				MdtaId, err := h.service.FindByVCMutasiDta(cc, &criteriaTB)
				if err != nil {
					return response.ErrorResponse(err).Send(c)
				}
				Map.MutasiDta = MutasiDta
				Map.IDMutasiDta = MdtaId.ID
				Map.IDWorksheetDetailMutasiDta = v.ID
			}
			if (ImportedWorksheetDetail.Name) == "Mutasi IA" {

				file, err := c.FormFile(ImportedWorksheetDetail.Code)
				if err != nil {
					continue
				}

				start := time.Now()
				src, err := file.Open()
				if err != nil {
					return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
				}
				defer src.Close()
				duration := time.Since(start)
				fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")

				filepath := fmt.Sprintf("uploaded/%s-%d-%d-%s", reupload, importedworksheet.CompanyID, importedworksheet.Versions, file.Filename)

				MutasiIa := path.Join(configs.App().StoragePath(), filepath)

				dst, err := os.Create(MutasiIa)
				if err != nil {
					return err
				}
				defer dst.Close()

				if _, err = io.Copy(dst, src); err != nil {
					return err
				}
				datePeriod, err := time.Parse(time.RFC3339, importedworksheet.Period)
				if err != nil {
					return err
				}
				period := datePeriod.Format("2006-01-02")
				criteriaTB := model.MutasiIaFilterModel{}
				criteriaTB.Period = &period
				criteriaTB.CompanyID = &importedworksheet.CompanyID
				criteriaTB.Versions = &importedworksheet.Versions
				MdtaId, err := h.service.FindByVCMutasiIa(cc, &criteriaTB)
				if err != nil {
					return response.ErrorResponse(err).Send(c)
				}
				Map.MutasiIa = MutasiIa
				Map.IDMutasiIa = MdtaId.ID
				Map.IDWorksheetDetailMutasiIa = v.ID
			}
			if (ImportedWorksheetDetail.Name) == "Mutasi RUA" {

				file, err := c.FormFile(ImportedWorksheetDetail.Code)
				if err != nil {
					continue
				}

				start := time.Now()
				src, err := file.Open()
				if err != nil {
					return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
				}
				defer src.Close()
				duration := time.Since(start)
				fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")
				filepath := fmt.Sprintf("uploaded/%s-%d-%d-%s", reupload, importedworksheet.CompanyID, importedworksheet.Versions, file.Filename)

				MutasiRua := path.Join(configs.App().StoragePath(), filepath)

				dst, err := os.Create(MutasiRua)
				if err != nil {
					return err
				}
				defer dst.Close()

				if _, err = io.Copy(dst, src); err != nil {
					return err
				}
				datePeriod, err := time.Parse(time.RFC3339, importedworksheet.Period)
				if err != nil {
					return err
				}
				period := datePeriod.Format("2006-01-02")
				criteriaTB := model.MutasiRuaFilterModel{}
				criteriaTB.Period = &period
				criteriaTB.CompanyID = &importedworksheet.CompanyID
				criteriaTB.Versions = &importedworksheet.Versions
				MruaId, err := h.service.FindByVCMutasiRua(cc, &criteriaTB)
				if err != nil {
					return response.ErrorResponse(err).Send(c)
				}
				Map.MutasiRua = MutasiRua
				Map.IDMutasiRua = MruaId.ID
				Map.IDWorksheetDetailMutasiRua = v.ID
			}
			if (ImportedWorksheetDetail.Name) == "Mutasi Persediaan" {

				file, err := c.FormFile(ImportedWorksheetDetail.Code)
				if err != nil {
					continue
				}

				start := time.Now()
				src, err := file.Open()
				if err != nil {
					return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
				}
				defer src.Close()
				duration := time.Since(start)
				fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")

				filepath := fmt.Sprintf("uploaded/%s-%d-%d-%s", reupload, importedworksheet.CompanyID, importedworksheet.Versions, file.Filename)

				MutasiPersediaan := path.Join(configs.App().StoragePath(), filepath)

				dst, err := os.Create(MutasiPersediaan)
				if err != nil {
					return err
				}
				defer dst.Close()

				if _, err = io.Copy(dst, src); err != nil {
					return err
				}

				if _, err = io.Copy(dst, src); err != nil {
					return err
				}
				datePeriod, err := time.Parse(time.RFC3339, importedworksheet.Period)
				if err != nil {
					return err
				}
				period := datePeriod.Format("2006-01-02")
				criteriaTB := model.MutasiPersediaanFilterModel{}
				criteriaTB.Period = &period
				criteriaTB.CompanyID = &importedworksheet.CompanyID
				criteriaTB.Versions = &importedworksheet.Versions
				MpId, err := h.service.FindByVCMutasiPersediaan(cc, &criteriaTB)
				if err != nil {
					return response.ErrorResponse(err).Send(c)
				}
				Map.MutasiPersediaan = MutasiPersediaan
				Map.IDMutasiPersediaan = MpId.ID
				Map.IDWorksheetDetailMutasiPersediaan = v.ID
			}
			if (ImportedWorksheetDetail.Name) == "Pembelian & Penjualan Berelasi" {

				file, err := c.FormFile(ImportedWorksheetDetail.Code)
				if err != nil {
					continue
				}

				start := time.Now()
				src, err := file.Open()
				if err != nil {
					return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
				}
				defer src.Close()
				duration := time.Since(start)
				fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")

				filepath := fmt.Sprintf("uploaded/%s-%d-%d-%s", reupload, importedworksheet.CompanyID, importedworksheet.Versions, file.Filename)

				PembelianPenjualanBerelasi := path.Join(configs.App().StoragePath(), filepath)

				dst, err := os.Create(PembelianPenjualanBerelasi)
				if err != nil {
					return err
				}
				defer dst.Close()

				if _, err = io.Copy(dst, src); err != nil {
					return err
				}
				datePeriod, err := time.Parse(time.RFC3339, importedworksheet.Period)
				if err != nil {
					return err
				}
				period := datePeriod.Format("2006-01-02")
				criteriaTB := model.PembelianPenjualanBerelasiFilterModel{}
				criteriaTB.Period = &period
				criteriaTB.CompanyID = &importedworksheet.CompanyID
				criteriaTB.Versions = &importedworksheet.Versions
				Ppb, err := h.service.FindByVCPembelianPenjualanBerelasi(cc, &criteriaTB)
				if err != nil {
					return response.ErrorResponse(err).Send(c)
				}
				Map.PembelianPenjualanBerelasi = PembelianPenjualanBerelasi
				Map.IDPembelianPenjualanBerelasi = Ppb.ID
				Map.IDWorksheetDetailPembelianPenjualanBerelasi = v.ID
			}
			if (ImportedWorksheetDetail.Name) == "Employee Benefit" {

				file, err := c.FormFile(ImportedWorksheetDetail.Code)
				if err != nil {
					continue
				}

				start := time.Now()
				src, err := file.Open()
				if err != nil {
					return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
				}
				defer src.Close()
				duration := time.Since(start)
				fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")
				filepath := fmt.Sprintf("uploaded/%s-%d-%d-%s", reupload, importedworksheet.CompanyID, importedworksheet.Versions, file.Filename)

				EmployeeBenefit := path.Join(configs.App().StoragePath(), filepath)

				dst, err := os.Create(EmployeeBenefit)
				if err != nil {
					return err
				}
				defer dst.Close()

				if _, err = io.Copy(dst, src); err != nil {
					return err
				}
				datePeriod, err := time.Parse(time.RFC3339, importedworksheet.Period)
				if err != nil {
					return err
				}
				period := datePeriod.Format("2006-01-02")
				criteriaTB := model.EmployeeBenefitFilterModel{}
				criteriaTB.Period = &period
				criteriaTB.CompanyID = &importedworksheet.CompanyID
				criteriaTB.Versions = &importedworksheet.Versions
				Ppb, err := h.service.FindByVCEmployeeBenefit(cc, &criteriaTB)
				if err != nil {
					return response.ErrorResponse(err).Send(c)
				}
				Map.EmployeeBenefit = EmployeeBenefit
				Map.IDEmployeeBenefit = Ppb.ID
				Map.IDWorksheetDetailEmployeeBenefit = v.ID
			}
		}
	}

	jsonStr, err := json.Marshal(Map)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	kafka.NewService("REUPLOAD").SendMessage("REUPLOAD", string(jsonStr))
	return response.SuccessResponse("REUPLOAD SUKSES").Send(cc)
}
func (h *handler) DownloadAll(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.ImportedWorksheetGetByIDRequest)

	if err = c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err = c.Validate(payload); err != nil {
		response := response.ErrorBuilder(&response.ErrorConstant.Validation, err)
		return response.Send(c)
	}

	fmt.Printf("%+v", payload)

	result, err := h.service.DownloadAll(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	tmpExtFile := "allworksheet.zip"
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	tmpFileName := fmt.Sprintf("uploaded/%s-%s%s", "downloadallworksheet", base64.StdEncoding.EncodeToString([]byte(timestamp)), tmpExtFile)

	downloadAll := path.Join(configs.App().StoragePath(), tmpFileName)

	// buat file zip kosong
	zipfile, err := os.Create(downloadAll)
	if err != nil {
		panic(err)
	}
	defer zipfile.Close()

	// buat writer untuk menulis ke file zip
	zipwriter := zip.NewWriter(zipfile)
	defer zipwriter.Close()

	// loop untuk membaca file excel dan menambahkan kontennya ke file zip
	for _, file := range result.FileName {

		xlsxFile, err := os.Open(file)
		if err != nil {
			return err
		}

		// tambahkan file excel ke file zip
		writer, err := zipwriter.Create(file)
		if err != nil {
			return err
		}
		if _, err := io.Copy(writer, xlsxFile); err != nil {
			return err
		}
		xlsxFile.Close()
	}

	xlsxFile, err := os.Open(downloadAll)
	if err != nil {
		return err
	}
	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", "allworksheet"))
	return c.Stream(http.StatusOK, "application/zip", xlsxFile)
}
func (h *handler) ImportJurnalAsync(c echo.Context) error {

	cc := c.(*abstraction.Context)
	payload := new(dto.ImportJurnal)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	map1 := kafka.JsonDataJurnal{
		TbID: payload.TbID,
		DataJurnal: payload.DataJurnal,
	}
	jsonStr, err := json.Marshal(map1)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	kafka.NewService("JURNAL").SendMessage("AJE", string(jsonStr))

	return response.SuccessResponse("SEND TO WORKER").Send(cc)
}
// Download
// @Summary Bulk Upload
// @Description Bulk Upload
// @Tags Bulk Upload
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "id path"
// @Success 200 {file} succes
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /import/{id} [get]
func (h *handler) Download(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.ImportedWorksheetDetailGetByIDRequest)

	if err = c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err = c.Validate(payload); err != nil {
		response := response.ErrorBuilder(&response.ErrorConstant.Validation, err)
		return response.Send(c)
	}

	fmt.Printf("%+v", payload)

	result, err := h.service.Download(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	f, err := os.Open(result.Note)
	if err != nil {
		return err
	}

	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", result.FileName))
	return c.Stream(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", f)
}

// Download
// @Summary Bulk Upload
// @Description Bulk Upload
// @Tags Bulk Upload
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "id path"
// @Success 200 {file} succes
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /import/{id} [get]
func (h *handler) DownloadTemplate(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.GetTemplateRequest)

	if err = c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err = c.Validate(payload); err != nil {
		response := response.ErrorBuilder(&response.ErrorConstant.Validation, err)
		return response.Send(c)
	}

	fmt.Printf("%+v", payload)


	result, err := h.service.ExportAllTemplate(cc)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	

	
	f, err := os.Open(*result)
	if err != nil {
		return err
	}

	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", f))
	return c.Stream(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", f)
}
