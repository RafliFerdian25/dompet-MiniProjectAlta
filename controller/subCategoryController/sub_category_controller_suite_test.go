package subCategoryController

import (
	"bytes"
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/dto"
	subCategoryMockService "dompet-miniprojectalta/service/subCategoryService/subCategoryMock"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type suiteSubCategory struct {
	suite.Suite
	subCategoryController *SubCategoryController
	mock                  *subCategoryMockService.SubCategoryMock
}

func (s *suiteSubCategory) SetupTest() {
	mock := &subCategoryMockService.SubCategoryMock{}
	s.mock = mock
	s.subCategoryController = &SubCategoryController{
		SubCategoryService: s.mock,
	}
}

func (s *suiteSubCategory) TestCreateSubCategory() {
	testCase := []struct {
		Name               string
		Method             string
		userId             uint
		Body               dto.SubCategoryDTO
		MockReturnError    error
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success create sub category",
			"POST",
			1,
			dto.SubCategoryDTO{
				CategoryID: 1,
				Name:       "test",
			},
			nil,
			http.StatusOK,
			"success create sub category",
		},
		{
			"fail bind data",
			"POST",
			1,
			dto.SubCategoryDTO{
				CategoryID: 1,
				Name:       "test",
			},
			nil,
			http.StatusInternalServerError,
			"fail bind data",
		},
		{
			"There is an empty field",
			"POST",
			1,
			dto.SubCategoryDTO{
				CategoryID: 1,
			},
			nil,
			http.StatusBadRequest,
			"There is an empty field",
		},
		{
			"fail create sub category",
			"POST",
			1,
			dto.SubCategoryDTO{
				Name:       "test",
				CategoryID: 1,
			},
			errors.New(constantError.ErrorNotAuthorized),
			http.StatusUnauthorized,
			"fail create sub category",
		},
		{
			"fail create sub category",
			"POST",
			1,
			dto.SubCategoryDTO{
				Name:       "test",
				CategoryID: 1,
			},
			errors.New("error"),
			http.StatusInternalServerError,
			"fail create sub category",
		},
	}
	for i, v := range testCase {
		body := dto.SubCategoryDTO{
			CategoryID: v.Body.CategoryID,
			Name:       v.Body.Name,
			UserID: &v.userId,
		}
		mockCall := s.mock.On("CreateSubCategory", body).Return(v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			res, _ := json.Marshal(v.Body)
			// Create request
			r := httptest.NewRequest(v.Method, "/subcategories", bytes.NewBuffer(res))
			if i != 1 {
				r.Header.Set("Content-Type", "application/json")
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			e.Validator = &helper.CustomValidator{
				Validator: validator.New(),
			}
			ctx := e.NewContext(r, w)
			ctx.SetPath("/subcategories")
			ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.subCategoryController.CreateSubCategory(ctx)
			s.NoError(err)
			s.Equal(v.ExpectedStatusCode, w.Code)

			var resp map[string]interface{}
			err = json.NewDecoder(w.Result().Body).Decode(&resp)
			s.NoError(err)

			s.Equal(v.ExpectedMesaage, resp["message"])
		})
		// remove mock
		mockCall.Unset()
	}
}

func (s *suiteSubCategory) TestGetSubCategoryByUser() {
	UserID := uint(1)

	testCase := []struct {
		Name               string
		Method             string
		userId             uint
		MockReturnBody     []dto.SubCategoryDTO
		MockReturnError    error
		HasReturnBody      bool
		ExpectedBody       []dto.SubCategoryDTO
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success get sub category by user",
			"GET",
			1,
			[]dto.SubCategoryDTO{
				{
					ID: 	   1,
					CategoryID: 2,
					Name:       "test",
					UserID: &UserID,
				},
				{
					ID: 	   2,
					CategoryID: 2,
					Name:       "test2",
					UserID: &UserID,
				},
			},
			nil,
			true,
			[]dto.SubCategoryDTO{
				{
					ID: 	   1,
					CategoryID: 2,
					Name:       "test",
					UserID: &UserID,
				},
				{
					ID: 	   2,
					CategoryID: 2,
					Name:       "test2",
					UserID: &UserID,
				},
			},
			http.StatusOK,
			"success get sub category by user",
		},
		{
			"fail get sub category by user",
			"GET",
			1,
			[]dto.SubCategoryDTO{},
			errors.New(constantError.ErrorNotAuthorized),
			false,
			[]dto.SubCategoryDTO{},
			http.StatusUnauthorized,
			"fail get sub category by user",
		},
		{
			"fail get sub category by user",
			"GET",
			1,
			[]dto.SubCategoryDTO{},
			errors.New("error"),
			false,
			[]dto.SubCategoryDTO{},
			http.StatusInternalServerError,
			"fail get sub category by user",
		},
	}
	for _, v := range testCase {
		mockCall := s.mock.On("GetSubCategoryByUser", v.userId).Return(v.MockReturnBody, v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			// Create request
			r := httptest.NewRequest(v.Method, "/subcategories", nil)
			// Create response recorder
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/subcategories")
			ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.subCategoryController.GetSubCategoryByUser(ctx)
			s.NoError(err)
			s.Equal(v.ExpectedStatusCode, w.Code)

			var resp map[string]interface{}
			err = json.NewDecoder(w.Result().Body).Decode(&resp)
			s.NoError(err)

			s.Equal(v.ExpectedMesaage, resp["message"])
			if v.HasReturnBody {
				s.Equal(v.ExpectedBody[0].Name, resp["sub_categories"].([]interface{})[0].(map[string]interface{})["name"])
				s.Equal(v.ExpectedBody[0].ID, uint(resp["sub_categories"].([]interface{})[0].(map[string]interface{})["id"].(float64)))
				s.Equal(v.ExpectedBody[1].Name, resp["sub_categories"].([]interface{})[1].(map[string]interface{})["name"])
				s.Equal(v.ExpectedBody[1].ID, uint(resp["sub_categories"].([]interface{})[1].(map[string]interface{})["id"].(float64)))
			}
		})
		// remove mock
		mockCall.Unset()
	}
}

func (s *suiteSubCategory) TestDeleteSubCategory() {
	testCase := []struct {
		Name               string
		Method             string
		userId             uint
		ParamId            string
		MockReturnError    error
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success delete sub category",
			"DELETE",
			1,
			"1",
			nil,
			http.StatusOK,
			"success delete sub category",
		},
		{
			"fail get id",
			"DELETE",
			1,
			"w",
			nil,
			http.StatusBadRequest,
			"fail get id",
		},
		{
			"fail delete sub category",
			"DELETE",
			1,
			"1",
			errors.New(constantError.ErrorNotAuthorized),
			http.StatusUnauthorized,
			"fail delete sub category",
		},
		{
			"fail delete sub category",
			"DELETE",
			1,
			"1",
			errors.New("error"),
			http.StatusInternalServerError,
			"fail delete sub category",
		},
	}
	for _, v := range testCase {
		id, _ := strconv.Atoi(v.ParamId)
		mockCall := s.mock.On("DeleteSubCategory", uint(id), v.userId).Return(v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			// Create request
			r := httptest.NewRequest(v.Method, "/debts", nil)
			// Create response recorder
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/debts")
			ctx.SetParamNames("id")
			ctx.SetParamValues(v.ParamId)
			ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.subCategoryController.DeleteSubCategory(ctx)
			s.NoError(err)
			s.Equal(v.ExpectedStatusCode, w.Code)

			var resp map[string]interface{}
			err = json.NewDecoder(w.Result().Body).Decode(&resp)
			s.NoError(err)

			s.Equal(v.ExpectedMesaage, resp["message"])
		})
		// remove mock
		mockCall.Unset()
	}
}

func (s *suiteSubCategory) TestUpdateSubCategory() {
	testCase := []struct {
		Name               string
		Method             string
		userId             uint
		Body               dto.SubCategoryDTO
		ParamID            string
		MockReturnError    error
		MockParamBody      dto.SubCategoryDTO
		HasReturnBody      bool
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success update sub category",
			"PUT",
			1,
			dto.SubCategoryDTO{
				CategoryID: 2,
				Name:       "test-update",
			},
			"1",
			nil,
			dto.SubCategoryDTO{
				CategoryID: 2,
				Name:       "test-update",
			},
			true,
			http.StatusOK,
			"success update sub category",
		},
		{
			"fail bind data",
			"PUT",
			1,
			dto.SubCategoryDTO{
				CategoryID: 2,
				Name:       "test-update",
			},
			"1",
			nil,
			dto.SubCategoryDTO{
				CategoryID: 2,
				Name:       "test-update",
			},
			true,
			http.StatusInternalServerError,
			"fail bind data",
		},
		{
			"fail get id",
			"PUT",
			1,
			dto.SubCategoryDTO{},
			"w",
			nil,
			dto.SubCategoryDTO{},
			true,
			http.StatusBadRequest,
			"fail get id",
		},
		{
			"fail update sub category",
			"PUT",
			1,
			dto.SubCategoryDTO{},
			"1",
			errors.New(constantError.ErrorNotAuthorized),
			dto.SubCategoryDTO{},
			true,
			http.StatusUnauthorized,
			"fail update sub category",
		},
		{
			"fail update sub category",
			"PUT",
			1,
			dto.SubCategoryDTO{},
			"1",
			errors.New("error"),
			dto.SubCategoryDTO{},
			true,
			http.StatusInternalServerError,
			"fail update sub category",
		},
	}
	for i, v := range testCase {
		id, _ := strconv.Atoi(v.ParamID)
		body := dto.SubCategoryDTO{
			ID: 	   uint(id),
			CategoryID: v.Body.CategoryID,
			Name:       v.Body.Name,
			UserID:    &v.userId,
		}
		mockCall := s.mock.On("UpdateSubCategory", body).Return(v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			res, _ := json.Marshal(v.Body)
			// Create request
			r := httptest.NewRequest(v.Method, "/subcategories", bytes.NewBuffer(res))
			if i != 1 {
				r.Header.Set("Content-Type", "application/json")
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/subcategories")
			ctx.SetParamNames("id")
			ctx.SetParamValues(v.ParamID)
			ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.subCategoryController.UpdateSubCategory(ctx)
			s.NoError(err)
			s.Equal(v.ExpectedStatusCode, w.Code)

			var resp map[string]interface{}
			err = json.NewDecoder(w.Result().Body).Decode(&resp)
			s.NoError(err)

			s.Equal(v.ExpectedMesaage, resp["message"])
		})
		// remove mock
		mockCall.Unset()
	}
}

func TestSuiteSubCategory(t *testing.T) {
	suite.Run(t, new(suiteSubCategory))
}
