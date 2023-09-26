package factory

import (
	"worker-consol/internal/database"
	"worker-consol/internal/repository"

	"gorm.io/gorm"
)

type Factory struct {
	Db                                         *gorm.DB
	UserRepository                             repository.User
	CoaRepository                              repository.Coa
	CoaGroupRepository                         repository.CoaGroup
	CompanyRepository                          repository.Company
	TrialBalanceRepository                     repository.TrialBalance
	TrialBalanceDetailRepository               repository.TrialBalanceDetail
	FormatterRepository                        repository.Formatter
	FormatterDetailRepository                  repository.FormatterDetail
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
	JpmRepository                              repository.Jpm
	JpmDetailRepository                        repository.JpmDetail
	AjeRepository                              repository.Adjustment
	AjeDetailRepository                        repository.AdjustmentDetail
	JcteRepository                             repository.Jcte
	JcteDetailRepository                       repository.JcteDetail
	JelimRepository                            repository.Jelim
	JelimDetailRepository                      repository.JelimDetail
	NotificationRepository                     repository.Notification
	EmployeeBenefitRepository                  repository.EmployeeBenefit
	EmployeeBenefitDetailRepository            repository.EmployeeBenefitDetail
	ConsolidationRepository                    repository.Consolidation
	ConsolidationDetailRepository              repository.ConsolidationDetail
	ConsolidationBridgeRepository              repository.ConsolidationBridge
	ConsolidationBridgeDetailRepository        repository.ConsolidationBridgeDetail
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
	f.CoaRepository = repository.NewCoa(f.Db)
	f.CoaGroupRepository = repository.NewCoaGroup(f.Db)
	f.CompanyRepository = repository.NewCompany(f.Db)
	f.TrialBalanceRepository = repository.NewTrialBalance(f.Db)
	f.TrialBalanceDetailRepository = repository.NewTrialBalanceDetail(f.Db)
	f.FormatterRepository = repository.NewFormatter(f.Db)
	f.FormatterDetailRepository = repository.NewFormatterDetail(f.Db)
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
	f.JpmRepository = repository.NewJpm(f.Db)
	f.JpmDetailRepository = repository.NewJpmDetail(f.Db)
	f.AjeRepository = repository.NewAdjustment(f.Db)
	f.AjeDetailRepository = repository.NewAdjustmentDetail(f.Db)
	f.JcteRepository = repository.NewJcte(f.Db)
	f.JcteDetailRepository = repository.NewJcteDetail(f.Db)
	f.JelimRepository = repository.NewJelim(f.Db)
	f.JelimDetailRepository = repository.NewJelimDetail(f.Db)
	f.NotificationRepository = repository.NewNotification(f.Db)
	f.EmployeeBenefitRepository = repository.NewEmployeeBenefit(f.Db)
	f.EmployeeBenefitDetailRepository = repository.NewEmployeeBenefitDetail(f.Db)
	f.ConsolidationRepository = repository.NewConsolidation(f.Db)
	f.ConsolidationDetailRepository = repository.NewConsolidationDetail(f.Db)
	f.ConsolidationBridgeRepository = repository.NewConsolidationBridge(f.Db)
	f.ConsolidationBridgeDetailRepository = repository.NewConsolidationBridgeDetail(f.Db)

}
