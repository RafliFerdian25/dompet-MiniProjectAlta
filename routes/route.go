package routes

import (
	"dompet-miniprojectalta/controller/userController"
	"dompet-miniprojectalta/repository/userRepository"
	"dompet-miniprojectalta/service/userService"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func New(db *gorm.DB) *echo.Echo {
	// Repositories
	userRepository := userRepository.NewUserRepository(db)

	// Services
	userService := userService.NewUserService(userRepository)

	// Controllers
	userController := userController.UserController{
		UserService: userService,
	}

	app := echo.New()

	/* 
		API Routes
	*/
	// User Routes
	app.POST("/signup", userController.CreateUser)
	app.POST("/login", userController.LoginUser)

	return app
}