package routes

import (
	"dompet-miniprojectalta/config"
	"dompet-miniprojectalta/controller/categoryController"
	"dompet-miniprojectalta/controller/userController"
	"dompet-miniprojectalta/repository/categoryRepository"
	"dompet-miniprojectalta/repository/userRepository"
	"dompet-miniprojectalta/service/categoryService"
	"dompet-miniprojectalta/service/userService"

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

	// Services
	userService := userService.NewUserService(userRepository)
	categoryService := categoryService.NewCategoryService(categoryRepository)

	// Controllers
	userController := userController.UserController{
		UserService: userService,
	}
	categoryController := categoryController.CategoryController{
		CategoryService: categoryService,
	}

	app := echo.New()

	/* 
		API Routes
	*/
	// User Routes
	app.POST("/signup", userController.CreateUser)
	app.POST("/login", userController.LoginUser)

	// Category Routes
	app.GET("/categories", categoryController.GetAllCategory, middleware.JWT([]byte(cfg.TOKEN_SECRET)))

	return app
}