package categoryController

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/model"
	categoryMockService "dompet-miniprojectalta/service/categoryService/categoryMock"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type suiteCategory struct {
	suite.Suite
	categoryController *CategoryController
	mock *categoryMockService.CategoryMock
}

func (s *suiteCategory) SetupTest() {
	mock := &categoryMockService.CategoryMock{}
	s.mock = mock
	s.categoryController = &CategoryController{
		CategoryService: s.mock,
	}
}

func (s *suiteCategory) TestGetCategoryByID() {
	testCase := []struct {
		Name               string
		Method             string
		ParamID             string
		MockReturnBody     model.Category
		MockReturnError    error
		HasReturnBody      bool
		ExpectedBody       model.Category
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success get category by id",
			"GET",
			"1",
			model.Category{
				
					Model:   gorm.Model{ID: 1},
					Name: "test",
					SubCategories: []model.SubCategory{},
				
			},
			nil,
			true,
			model.Category{
				
					Model:   gorm.Model{ID: 1},
					Name: "test",
					SubCategories: []model.SubCategory{},
				
			},
			http.StatusOK,
			"success get category by id",
		},
		{
			"fail get id",
			"GET",
			"w",
			model.Category{},
			nil,
			false,
			model.Category{},
			http.StatusBadRequest,
			"fail get id",
		},
		{
			"fail get category by id",
			"GET",
			"1",
			model.Category{},
			errors.New(constantError.ErrorNotAuthorized),
			false,
			model.Category{
				
					Model:   gorm.Model{ID: 1},
					Name: "test",
					SubCategories: []model.SubCategory{},
				
			},
			http.StatusUnauthorized,
			"fail get category by id",
		},
		{
			"fail get category by id",
			"GET",
			"1",
			model.Category{},
			errors.New("error"),
			false,
			model.Category{
				
					Model:   gorm.Model{ID: 1},
					Name: "test",
					SubCategories: []model.SubCategory{},
				
			},
			http.StatusInternalServerError,
			"fail get category by id",
		},
	}
	for _, v := range testCase {
		id, _ := strconv.Atoi(v.ParamID)
		mockCall := s.mock.On("GetCategoryByID", uint(id)).Return(v.MockReturnBody, v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			// Create request
			r := httptest.NewRequest(v.Method, "/categories/"+v.ParamID, nil)
			// Create response recorder
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/categories/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues(v.ParamID)
			ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.categoryController.GetCategoryByID(ctx)
			s.NoError(err)
			s.Equal(v.ExpectedStatusCode, w.Code)

			var resp map[string]interface{}
			err = json.NewDecoder(w.Result().Body).Decode(&resp)
			s.NoError(err)

			s.Equal(v.ExpectedMesaage, resp["message"])

			if v.HasReturnBody {
				s.Equal(v.ExpectedBody.Name, resp["categories"].(map[string]interface{})["name"])
				s.Equal(v.ExpectedBody.Name, resp["categories"].(map[string]interface{})["name"])
			}
		})
		// remove mock
		mockCall.Unset()
	}
}

func (s *suiteCategory) TestGetAllCategory() {
	testCase := []struct {
		Name               string
		Method             string
		MockReturnBody     []model.Category
		MockReturnError    error
		HasReturnBody      bool
		ExpectedBody       []model.Category
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success get category",
			"GET",
			[]model.Category{
				{
					Model:   gorm.Model{ID: 1},
					Name: "test",
					SubCategories: []model.SubCategory{},
				},
				{
					Model:   gorm.Model{ID: 2},
					Name: "test2",
					SubCategories: []model.SubCategory{},
				},
			},
			nil,
			true,
			[]model.Category{
				{
					Model:   gorm.Model{ID: 1},
					Name: "test",
					SubCategories: []model.SubCategory{},
				},
				{
					Model:   gorm.Model{ID: 2},
					Name: "test2",
					SubCategories: []model.SubCategory{},
				},
			},
			http.StatusOK,
			"success get category",
		},
		{
			"fail get all category",
			"GET",
			[]model.Category{},
			errors.New(constantError.ErrorNotAuthorized),
			false,
			[]model.Category{
				{
					Model:   gorm.Model{ID: 1},
					Name: "test",
					SubCategories: []model.SubCategory{},
				},
				{
					Model:   gorm.Model{ID: 2},
					Name: "test2",
					SubCategories: []model.SubCategory{},
				},
			},
			http.StatusUnauthorized,
			"fail get all category",
		},
		{
			"fail get all category",
			"GET",
			[]model.Category{},
			errors.New("error"),
			false,
			[]model.Category{
				{
					Model:   gorm.Model{ID: 1},
					Name: "test",
					SubCategories: []model.SubCategory{},
				},
				{
					Model:   gorm.Model{ID: 2},
					Name: "test2",
					SubCategories: []model.SubCategory{},
				},
			},
			http.StatusInternalServerError,
			"fail get all category",
		},
	}
	for _, v := range testCase {
		mockCall := s.mock.On("GetAllCategory").Return(v.MockReturnBody, v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			// Create request
			r := httptest.NewRequest(v.Method, "/categories/", nil)
			// Create response recorder
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/categories")
			// ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.categoryController.GetAllCategory(ctx)
			s.NoError(err)
			s.Equal(v.ExpectedStatusCode, w.Code)

			var resp map[string]interface{}
			err = json.NewDecoder(w.Result().Body).Decode(&resp)
			s.NoError(err)

			s.Equal(v.ExpectedMesaage, resp["message"])

			if v.HasReturnBody {
				s.Equal(v.ExpectedBody[0].Name, resp["categories"].([]interface{})[0].(map[string]interface{})["name"])
				s.Equal(v.ExpectedBody[1].Name, resp["categories"].([]interface{})[1].(map[string]interface{})["name"])
			}
		})
		// remove mock
		mockCall.Unset()
	}
}

func TestSuiteCategory(t *testing.T) {
	suite.Run(t, new(suiteCategory))
}
