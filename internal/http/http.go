package http

import (
	"fmt"
	"mcash-finance-console-core/configs"
	docs "mcash-finance-console-core/docs"
	"mcash-finance-console-core/internal/app/accessscope"
	"mcash-finance-console-core/internal/app/accessscopedetail"
	"mcash-finance-console-core/internal/app/adjustment"
	"mcash-finance-console-core/internal/app/adjustmentdetail"
	"mcash-finance-console-core/internal/app/agingutangpiutang"
	"mcash-finance-console-core/internal/app/agingutangpiutangdetail"
	"mcash-finance-console-core/internal/app/auth"
	"mcash-finance-console-core/internal/app/coa"
	"mcash-finance-console-core/internal/app/coagroup"
	"mcash-finance-console-core/internal/app/coatype"
	"mcash-finance-console-core/internal/app/company"
	"mcash-finance-console-core/internal/app/parameter"

	// "mcash-finance-console-core/internal/app/consolidation"
	"mcash-finance-console-core/internal/app/approvalvalidation"
	"mcash-finance-console-core/internal/app/consolidation"
	"mcash-finance-console-core/internal/app/consolidationdetail"
	"mcash-finance-console-core/internal/app/employeebenefit"
	"mcash-finance-console-core/internal/app/employeebenefitdetail"
	"mcash-finance-console-core/internal/app/export"
	"mcash-finance-console-core/internal/app/importedworksheet"
	"mcash-finance-console-core/internal/app/importedworksheetdetail"
	"mcash-finance-console-core/internal/app/imports"
	"mcash-finance-console-core/internal/app/investasinontbk"
	"mcash-finance-console-core/internal/app/investasinontbkdetail"
	"mcash-finance-console-core/internal/app/investasitbk"
	"mcash-finance-console-core/internal/app/investasitbkdetail"
	"mcash-finance-console-core/internal/app/jcte"
	"mcash-finance-console-core/internal/app/jctedetail"
	"mcash-finance-console-core/internal/app/jelim"
	"mcash-finance-console-core/internal/app/jelimdetail"
	"mcash-finance-console-core/internal/app/jpm"
	"mcash-finance-console-core/internal/app/jpmdetail"
	"mcash-finance-console-core/internal/app/mutasidta"
	"mcash-finance-console-core/internal/app/mutasidtadetail"
	"mcash-finance-console-core/internal/app/mutasifa"
	"mcash-finance-console-core/internal/app/mutasifadetail"
	"mcash-finance-console-core/internal/app/mutasiia"
	"mcash-finance-console-core/internal/app/mutasiiadetail"
	"mcash-finance-console-core/internal/app/mutasipersediaan"
	"mcash-finance-console-core/internal/app/mutasipersediaandetail"
	"mcash-finance-console-core/internal/app/mutasirua"
	"mcash-finance-console-core/internal/app/mutasiruadetail"
	"mcash-finance-console-core/internal/app/notification"
	"mcash-finance-console-core/internal/app/pembelianpenjualanberelasi"
	"mcash-finance-console-core/internal/app/pembelianpenjualanberelasidetail"
	"mcash-finance-console-core/internal/app/permissiondef"
	"mcash-finance-console-core/internal/app/role"
	"mcash-finance-console-core/internal/app/rolepermission"
	"mcash-finance-console-core/internal/app/rolepermissionapi"
	"mcash-finance-console-core/internal/app/sample"
	"mcash-finance-console-core/internal/app/trialbalance"
	"mcash-finance-console-core/internal/app/trialbalancedetail"
	"mcash-finance-console-core/internal/app/user"
	"mcash-finance-console-core/internal/app/validation"
	"mcash-finance-console-core/internal/app/formatterdetail"
	"mcash-finance-console-core/internal/factory"
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Init(e *echo.Echo, f *factory.Factory) {
	// index
	e.GET("/", func(c echo.Context) error {
		message := fmt.Sprintf("Welcome to %s version %s", configs.App().Name(), configs.App().Version())
		return c.String(http.StatusOK, message)
	})

	// doc
	docs.SwaggerInfo.Title = configs.App().Name()
	docs.SwaggerInfo.Version = configs.App().Version()
	docs.SwaggerInfo.Host = configs.App().Host()
	docs.SwaggerInfo.Schemes = configs.App().Schemes()
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// routes
	auth.NewHandler(f).Route(e.Group("/auth"))
	sample.NewHandler(f).Route(e.Group("/samples"))
	coa.NewHandler(f).Route(e.Group("/coa"))
	coagroup.NewHandler(f).Route(e.Group("/coa-group"))
	coatype.NewHandler(f).Route(e.Group("/coa-type"))
	company.NewHandler(f).Route(e.Group("/company"))
	trialbalance.NewHandler(f).Route(e.Group("/trial-balance"))
	trialbalancedetail.NewHandler(f).Route(e.Group("/trial-balance-detail"))
	agingutangpiutang.NewHandler(f).Route(e.Group("/aging-utang-piutang"))
	agingutangpiutangdetail.NewHandler(f).Route(e.Group("/aging-utang-piutang-detail"))
	mutasipersediaan.NewHandler(f).Route(e.Group("/mutasi-persediaan"))
	mutasipersediaandetail.NewHandler(f).Route(e.Group("/mutasi-persediaan-detail"))
	pembelianpenjualanberelasi.NewHandler(f).Route(e.Group("/pembelian-penjualan-berelasi"))
	pembelianpenjualanberelasidetail.NewHandler(f).Route(e.Group("/pembelian-penjualan-berelasi-detail"))
	mutasifa.NewHandler(f).Route(e.Group("/mutasi-fa"))
	mutasifadetail.NewHandler(f).Route(e.Group("/mutasi-fa-detail"))
	mutasiia.NewHandler(f).Route(e.Group("/mutasi-ia"))
	mutasiiadetail.NewHandler(f).Route(e.Group("/mutasi-ia-detail"))
	mutasirua.NewHandler(f).Route(e.Group("/mutasi-rua"))
	mutasiruadetail.NewHandler(f).Route(e.Group("/mutasi-rua-detail"))
	mutasidta.NewHandler(f).Route(e.Group("/mutasi-dta"))
	mutasidtadetail.NewHandler(f).Route(e.Group("/mutasi-dta-detail"))
	investasitbk.NewHandler(f).Route(e.Group("/investasi-tbk"))
	investasitbkdetail.NewHandler(f).Route(e.Group("/investasi-tbk-detail"))
	investasinontbk.NewHandler(f).Route(e.Group("/investasi-non-tbk"))
	investasinontbkdetail.NewHandler(f).Route(e.Group("/investasi-non-tbk-detail"))
	export.NewHandler(f).Route(e.Group("/export"))
	//static
	e.Static("/assets", "assets")
	e.Static("/data/templates", "templates")
	e.Static("/data/uploaded", "uploaded")

	adjustment.NewHandler(f).Route(e.Group("/adjustment"))
	adjustmentdetail.NewHandler(f).Route(e.Group("/adjustment-detail"))
	jcte.NewHandler(f).Route(e.Group("/jcte"))
	jctedetail.NewHandler(f).Route(e.Group("/jcte-detail"))
	jpm.NewHandler(f).Route(e.Group("/jpm"))
	jpmdetail.NewHandler(f).Route(e.Group("/jpm-detail"))
	jelim.NewHandler(f).Route(e.Group("/jelim"))
	jelimdetail.NewHandler(f).Route(e.Group("/jelim-detail"))
	imports.NewHandler(f).Route(e.Group("/import"))
	employeebenefit.NewHandler(f).Route(e.Group("/employee-benefit"))
	employeebenefitdetail.NewHandler(f).Route(e.Group("/employee-benefit-detail"))
	accessscope.NewHandler(f).Route(e.Group("/access-scope"))
	accessscopedetail.NewHandler(f).Route(e.Group("/access-scope-detail"))
	rolepermissionapi.NewHandler(f).Route(e.Group("/role-permission-api"))
	role.NewHandler(f).Route(e.Group("/role"))
	rolepermission.NewHandler(f).Route(e.Group("/role-permission"))
	permissiondef.NewHandler(f).Route(e.Group("/permission-def"))
	user.NewHandler(f).Route(e.Group("/user"))
	importedworksheet.NewHandler(f).Route(e.Group("/imported-worksheet"))
	importedworksheetdetail.NewHandler(f).Route(e.Group("/imported-worksheet-detail"))
	validation.NewHandler(f).Route(e.Group("/validation"))
	notification.NewHandler(f).Route(e.Group("/notification"))
	consolidation.NewHandler(f).Route(e.Group("/consolidation"))
	consolidationdetail.NewHandler(f).Route(e.Group("/consolidation-detail"))
	parameter.NewHandler(f).Route(e.Group("/parameter"))
	approvalvalidation.NewHandler(f).Route(e.Group("/approval-validation"))
	formatterdetail.NewHandler(f).Route(e.Group("/formatter-detail"))
}
