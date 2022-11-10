package userController

import (
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"
	"dompet-miniprojectalta/service/userService"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	UserService userService.UserService
}

func (u *UserController) CreateUser(c echo.Context) error {
	var user dto.UserDTO
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(500, echo.Map{
			"message": "fail bind data",
			"error":   err.Error(),
		})
	}

	// validate data user
	err = c.Validate(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "There is an empty field",
			"error":   err.Error(),
		})
	}

	err = u.UserService.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail create user",
			"error":   err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"message": "success create user",
	})
}

func (u *UserController) LoginUser(c echo.Context) error {
	var user model.User
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail bind data",
			"error":   err.Error(),
		})
	}

	user, err = u.UserService.LoginUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail login",
			"error":   err.Error(),
		})
	}
	
	token, errToken := helper.CreateToken(user.ID, user.Name)

	if errToken != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail create token",
			"error":   errToken,
		})
	}

	userResponse := dto.UserResponseDTO{
		Name: user.Name,
		Email: user.Email,
		Token: token,
	}

	return c.JSON(200, echo.Map{
		"message": "success login",
		"user":   userResponse,
	})
}