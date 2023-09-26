package factory

import (
	"mcash-finance-console-core/internal/database"
	"mcash-finance-console-core/internal/repository"

	"gorm.io/gorm"
)

type Factory struct {
	Db                                         *gorm.DB
	UserRepository                             repository.User
	SampleRepository                           repository.Sample
	CoaRepository                              repository.Coa
	CoaDevRepository                           repository.CoaDev
	CoaGroupRepository                         repository.CoaGroup
	CoaTypeRepository                          repository.CoaType
	CompanyRepository                          repository.Company
	TrialBalanceRepository                     repository.TrialBalance
	TrialBalanceDetailRepository               repository.TrialBalanceDetail
	FormatterRepository                        repository.Formatter
	FormatterDetailRepository                  repository.FormatterDetail
	FormatterDetailDevRepository               repository.FormatterDetailDev
	AgingUtangPiutangRepository                repository.AgingUtangPiutang
	AgingUtangPiutangDetailRepository          repository.AgingUtangPiutangDetail
	MutasiPersediaanRepository                 repository.MutasiPersediaan
	MutasiPersediaanDetailRepository           repository.MutasiPersediaanDetail
	PembelianPenjualanBerelasiRepository       repository.PembelianPenjualanBerelasi
	PembelianPenjualanBerelasiDetailRepository repository.PembelianPenjualanBerelasiDetail
	MutasiFaRepository                         repository.MutasiFa
	MutasiFaDetailRepository                   repository.MutasiFaDetail
	MutasiIaRepository                         repository.MutasiIa
	MutasiIaDetailRepository                   repository.MutasiIaDetail
	MutasiRuaRepository                        repository.MutasiRua
	MutasiRuaDetailRepository                  repository.MutasiRuaDetail
	MutasiDtaRepository                        repository.MutasiDta
	MutasiDtaDetailRepository                  repository.MutasiDtaDetail
	InvestasiNonTbkRepository                  repository.InvestasiNonTbk
	InvestasiNonTbkDetailRepository            repository.InvestasiNonTbkDetail
	InvestasiTbkRepository                     repository.InvestasiTbk
	InvestasiTbkDetailRepository               repository.InvestasiTbkDetail
	ParameterRepository                        repository.Parameter
	FormatterBridgesRepository                 repository.FormatterBridges
	AdjustmentRepository                       repository.Adjustment
	AdjustmentDetailRepository                 repository.AdjustmentDetail
	JcteRepository                             repository.Jcte
	JcteDetailRepository                       repository.JcteDetail
	JpmRepository                              repository.Jpm
	JpmDetailRepository                        repository.JpmDetail
	JelimRepository                            repository.Jelim
	JelimDetailRepository                      repository.JelimDetail
	EmployeeBenefitRepository                  repository.EmployeeBenefit
	EmployeeBenefitDetailRepository            repository.EmployeeBenefitDetail
	AccessScopeRepository                      repository.AccessScope
	AccessScopeDetailRepository                repository.AccessScopeDetail
	RoleRepository                             repository.Role
	RolePermissionApiRepository                repository.RolePermissionApi
	RolePermissionRepository                   repository.RolePermission
	PermissionDefRepository                    repository.PermissionDef
	ImportedWorksheetRepository                repository.ImportedWorksheet
	ImportedWorksheetDetailRepository          repository.ImportedWorksheetDetail
	NotificationRepository                     repository.Notification
	ValidationRepository                       repository.Validation
	ValidationDetailRepository                 repository.ValidationDetail
	ConsolidationRepository                    repository.Consolidation
	ConsolidationDetailRepository              repository.ConsolidationDetail
	ApprovalValidationRepository               repository.ApprovalValidation
}

func NewFactory() *Factory {
	f := &Factory{}
	f.SetupDb()
	f.SetupRepository()

	return f
}

func (f *Factory) SetupDb() {
	db, err := database.Connection("SAMPLE1")
	if err != nil {
		panic("Failed setup db, connection is undefined")
	}
	f.Db = db
}

func (f *Factory) SetupRepository() {
	if f.Db == nil {
		panic("Failed setup repository, db is undefined")
	}

	f.UserRepository = repository.NewUser(f.Db)
	f.SampleRepository = repository.NewSample(f.Db)
	f.CoaRepository = repository.NewCoa(f.Db)
	f.CoaDevRepository = repository.NewCoaDev(f.Db)
	f.CoaGroupRepository = repository.NewCoaGroup(f.Db)
	f.CoaTypeRepository = repository.NewCoaType(f.Db)
	f.CompanyRepository = repository.NewCompany(f.Db)
	f.TrialBalanceRepository = repository.NewTrialBalance(f.Db)
	f.TrialBalanceDetailRepository = repository.NewTrialBalanceDetail(f.Db)
	f.FormatterRepository = repository.NewFormatter(f.Db)
	f.FormatterDetailRepository = repository.NewFormatterDetail(f.Db)
	f.FormatterDetailDevRepository = repository.NewFormatterDetailDev(f.Db)
	f.AgingUtangPiutangRepository = repository.NewAgingUtangPiutang(f.Db)
	f.AgingUtangPiutangDetailRepository = repository.NewAgingUtangPiutangDetail(f.Db)
	f.MutasiPersediaanRepository = repository.NewMutasiPersediaan(f.Db)
	f.MutasiPersediaanDetailRepository = repository.NewMutasiPersediaanDetail(f.Db)
	f.PembelianPenjualanBerelasiRepository = repository.NewPembelianPenjualanBerelasi(f.Db)
	f.PembelianPenjualanBerelasiDetailRepository = repository.NewPembelianPenjualanBerelasiDetail(f.Db)
	f.MutasiFaRepository = repository.NewMutasiFa(f.Db)
	f.MutasiFaDetailRepository = repository.NewMutasiFaDetail(f.Db)
	f.MutasiIaRepository = repository.NewMutasiIa(f.Db)
	f.MutasiIaDetailRepository = repository.NewMutasiIaDetail(f.Db)
	f.MutasiRuaRepository = repository.NewMutasiRua(f.Db)
	f.MutasiRuaDetailRepository = repository.NewMutasiRuaDetail(f.Db)
	f.MutasiDtaRepository = repository.NewMutasiDta(f.Db)
	f.MutasiDtaDetailRepository = repository.NewMutasiDtaDetail(f.Db)
	f.InvestasiNonTbkRepository = repository.NewInvestasiNonTbk(f.Db)
	f.InvestasiNonTbkDetailRepository = repository.NewInvestasiNonTbkDetail(f.Db)
	f.InvestasiTbkRepository = repository.NewInvestasiTbk(f.Db)
	f.InvestasiTbkDetailRepository = repository.NewInvestasiTbkDetail(f.Db)
	f.ParameterRepository = repository.NewParameter(f.Db)
	f.FormatterBridgesRepository = repository.NewFormatterBridges(f.Db)
	f.AdjustmentRepository = repository.NewAdjustment(f.Db)
	f.AdjustmentDetailRepository = repository.NewAdjustmentDetail(f.Db)
	f.ImportedWorksheetRepository = repository.NewImportedWorksheet(f.Db)
	f.JpmRepository = repository.NewJpm(f.Db)
	f.JpmDetailRepository = repository.NewJpmDetail(f.Db)
	f.JcteRepository = repository.NewJcte(f.Db)
	f.JcteDetailRepository = repository.NewJcteDetail(f.Db)
	f.JelimRepository = repository.NewJelim(f.Db)
	f.JelimDetailRepository = repository.NewJelimDetail(f.Db)
	f.EmployeeBenefitRepository = repository.NewEmployeeBenefit(f.Db)
	f.EmployeeBenefitDetailRepository = repository.NewEmployeeBenefitDetail(f.Db)
	f.AccessScopeRepository = repository.NewAccessScope(f.Db)
	f.AccessScopeDetailRepository = repository.NewAccessScopeDetail(f.Db)
	f.RoleRepository = repository.NewRole(f.Db)
	f.RolePermissionApiRepository = repository.NewRolePermissionApi(f.Db)
	f.RolePermissionRepository = repository.NewRolePermission(f.Db)
	f.PermissionDefRepository = repository.NewPermissionDef(f.Db)
	f.ImportedWorksheetDetailRepository = repository.NewImportedWorksheetDetail(f.Db)
	f.NotificationRepository = repository.NewNotification(f.Db)
	f.ValidationRepository = repository.NewValidation(f.Db)
	f.ValidationDetailRepository = repository.NewValidationDetail(f.Db)
	f.ConsolidationRepository = repository.NewConsolidation(f.Db)
	f.ConsolidationDetailRepository = repository.NewConsolidationDetail(f.Db)
	f.ApprovalValidationRepository = repository.NewApprovalValidation(f.Db)
}
