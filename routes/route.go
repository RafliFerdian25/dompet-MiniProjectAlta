package routes

import (
	"dompet-miniprojectalta/config"
	"dompet-miniprojectalta/controller/accountController"
	"dompet-miniprojectalta/controller/categoryController"
	"dompet-miniprojectalta/controller/subCategoryController"
	"dompet-miniprojectalta/controller/userController"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/repository/accountRepository"
	"dompet-miniprojectalta/repository/categoryRepository"
	"dompet-miniprojectalta/repository/subCategoryRepository"
	"dompet-miniprojectalta/repository/userRepository"
	"dompet-miniprojectalta/service/accountService"
	"dompet-miniprojectalta/service/categoryService"
	"dompet-miniprojectalta/service/subCategoryService"
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
	subcategoryRepository := subCategoryRepository.NewCategoryRepository(db)
	accountRepository := accountRepository.NewAccountRepository(db)

	// Services
	userService := userService.NewUserService(userRepository)
	categoryService := categoryService.NewCategoryService(categoryRepository)
	subcategoryService := subCategoryService.NewSubCategoryService(subcategoryRepository)
	accountService := accountService.NewAccountService(accountRepository)

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
	appCategory.GET("/:categoryId", categoryController.GetCategoryByID)
	
	// SubCategory Routes
	appSubCategory := app.Group("/subcategories", middleware.JWTWithConfig(config))
	appSubCategory.POST("", subcategoryController.CreateSubCategory)
	appSubCategory.GET("/userid", subcategoryController.GetSubCategoryByUser)
	appSubCategory.DELETE("/:id", subcategoryController.DeleteSubCategory)
	appSubCategory.PUT("/:id", subcategoryController.UpdateSubCategory)

	// Account Routes
	appAccount := app.Group("/accounts", middleware.JWTWithConfig(config))
	appAccount.POST("", accountController.CreateAccount)
	appAccount.PUT("/:id", accountController.UpdateAccount)
	// appAccount.GET("/userid", accountController.GetAccountByUser)

	return app
}