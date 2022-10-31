package routes

import (
	"dompet-miniprojectalta/config"
	"dompet-miniprojectalta/controller/categoryController"
	"dompet-miniprojectalta/controller/subCategoryController"
	"dompet-miniprojectalta/controller/userController"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/repository/categoryRepository"
	"dompet-miniprojectalta/repository/subCategoryRepository"
	"dompet-miniprojectalta/repository/userRepository"
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

	// Services
	userService := userService.NewUserService(userRepository)
	categoryService := categoryService.NewCategoryService(categoryRepository)
	subcategoryService := subCategoryService.NewSubCategoryService(subcategoryRepository)

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
	app.GET("/categories", categoryController.GetAllCategory, middleware.JWTWithConfig(config))

	// SubCategory Routes
	app.POST("/subcategories", subcategoryController.CreateSubCategory, middleware.JWTWithConfig(config))

	return app
}