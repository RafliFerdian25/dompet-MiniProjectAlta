package routes

import (
	"dompet-miniprojectalta/config"
	"dompet-miniprojectalta/controller/accountController"
	"dompet-miniprojectalta/controller/categoryController"
	"dompet-miniprojectalta/controller/debtController"
	"dompet-miniprojectalta/controller/reportController"
	"dompet-miniprojectalta/controller/subCategoryController"
	"dompet-miniprojectalta/controller/transactionAccController"
	"dompet-miniprojectalta/controller/transactionController"
	"dompet-miniprojectalta/controller/userController"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/repository/accountRepository"
	"dompet-miniprojectalta/repository/categoryRepository"
	"dompet-miniprojectalta/repository/debtRepository"
	"dompet-miniprojectalta/repository/reportRepository"
	"dompet-miniprojectalta/repository/subCategoryRepository"
	"dompet-miniprojectalta/repository/transactionAccRepository"
	"dompet-miniprojectalta/repository/transactionRepository"
	"dompet-miniprojectalta/repository/userRepository"
	"dompet-miniprojectalta/service/accountService"
	"dompet-miniprojectalta/service/categoryService"
	"dompet-miniprojectalta/service/debtService"
	"dompet-miniprojectalta/service/reportService"
	"dompet-miniprojectalta/service/subCategoryService"
	"dompet-miniprojectalta/service/transactionAccService"
	"dompet-miniprojectalta/service/transactionService"
	"dompet-miniprojectalta/service/userService"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func New(db *gorm.DB) *echo.Echo {
	
	// Config
	cfg := config.Cfg
	
	// Repositories
	userRepository := userRepository.NewUserRepository(db)
	categoryRepository := categoryRepository.NewCategoryRepository(db)
	subcategoryRepository := subCategoryRepository.NewSubCategoryRepository(db)
	accountRepository := accountRepository.NewAccountRepository(db)
	transactionRepository := transactionRepository.NewTransactionRepository(db)
	debtRepository := debtRepository.NewDebtRepository(db)
	transAccRepository := transactionAccRepository.NewTransAccRepo(db)
	reportRepository := reportRepository.NewReportRepository(db)

	// Services
	userService := userService.NewUserService(userRepository)
	categoryService := categoryService.NewCategoryService(categoryRepository)
	subcategoryService := subCategoryService.NewSubCategoryService(subcategoryRepository)
	accountService := accountService.NewAccountService(accountRepository)
	transactionService := transactionService.NewTransactionService(transactionRepository, accountRepository, subcategoryRepository)
	debtService := debtService.NewDebtService(debtRepository, accountRepository, subcategoryRepository)
	transAccService := transactionAccService.NewTransAccService(transAccRepository, accountRepository)
	reportService := reportService.NewReportService(reportRepository)

	// Controllers
	userController := userController.UserController{
		UserService: userService,
	}
	categoryController := categoryController.CategoryController{
		CategoryService: categoryService,
	}
	subcategoryController := subCategoryController.SubCategoryController{
		SubCategoryService: subcategoryService,
	}
	accountController := accountController.AccountController{
		AccountService: accountService,
	}
	transactionController := transactionController.TransactionController{
		TransactionService: transactionService,
	}
	debtController := debtController.DebtController{
		DebtService: debtService,
	}
	transAccController := transactionAccController.TransactionAccController{
		TransAccService: transAccService,
	}
	reportController := reportController.ReportController{
		ReportService: reportService,
	}

	app := echo.New()

	app.Validator = &helper.CustomValidator{
		Validator: validator.New(),
	}

	/* 
		API Routes
	*/
	config := middleware.JWTConfig{
		Claims:     &helper.JWTCustomClaims{},
		SigningKey: []byte(cfg.TOKEN_SECRET),
	}

	// User Routes
	app.POST("/signup", userController.CreateUser)
	app.POST("/login", userController.LoginUser)

	// Category Routes
	appCategory := app.Group("/categories", middleware.JWTWithConfig(config))
	appCategory.GET("", categoryController.GetAllCategory)
	appCategory.GET("/:id", categoryController.GetCategoryByID)
	
	// SubCategory Routes
	appSubCategory := app.Group("/subcategories", middleware.JWTWithConfig(config))
	appSubCategory.POST("", subcategoryController.CreateSubCategory)
	appSubCategory.GET("", subcategoryController.GetSubCategoryByUser)
	appSubCategory.DELETE("/:id", subcategoryController.DeleteSubCategory)
	appSubCategory.PUT("/:id", subcategoryController.UpdateSubCategory)

	// Account Routes
	appAccount := app.Group("/accounts", middleware.JWTWithConfig(config))
	appAccount.POST("", accountController.CreateAccount)
	appAccount.PUT("/:id", accountController.UpdateAccount)
	appAccount.DELETE("/:id", accountController.DeleteAccount)
	appAccount.GET("", accountController.GetAccountByUser)

	// Transaction Routes
	appTransaction := app.Group("/transactions", middleware.JWTWithConfig(config))
	appTransaction.POST("",transactionController.CreateTransaction)
	appTransaction.PUT("/:id", transactionController.UpdateTransaction)
	appTransaction.DELETE("/:id", transactionController.DeleteTransaction)
	appTransaction.GET("", transactionController.GetTransaction)

	// Debt Routes
	appDebt := app.Group("/debts", middleware.JWTWithConfig(config))
	appDebt.POST("", debtController.CreateDebt)
	appDebt.DELETE("/:id", debtController.DeleteDebt)
	appDebt.GET("", debtController.GetDebt)

	//Transaction account Routes
	appTransAcc := app.Group("/transaction-accounts", middleware.JWTWithConfig(config)) 
	appTransAcc.POST("", transAccController.CreateTransactionAccount)
	appTransAcc.DELETE("/:id", transAccController.DeleteTransactionAccount)
	appTransAcc.GET("", transAccController.GetTransactionAccount)

	// Report Routes
	appReport := app.Group("/reports", middleware.JWTWithConfig(config))
	appReport.GET("/analytic", reportController.GetAnalyticPeriod)
	appReport.GET("/subcategory", reportController.GetReportbyCategory)
	appReport.GET("/cashflow", reportController.GetCashflow)

	return app
}
