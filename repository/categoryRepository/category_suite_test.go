package categoryRepository

import (
	"database/sql/driver"
	"dompet-miniprojectalta/models/dto"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type suiteCategory struct {
	suite.Suite
	categoryRepository CategoryRepository
	mock               sqlmock.Sqlmock
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func (s *suiteCategory) SetupSuite() {
	dbGormPalsu, mocking, err := sqlmock.New()
	s.NoError(err)

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      dbGormPalsu,
	}))

	s.mock = mocking

	NewCategoryRepository := NewCategoryRepository(dbGorm)
	s.categoryRepository = NewCategoryRepository
}

func (s *suiteCategory) TestGetCategoryByID() {
	userId := uint(1)
	testCase := []struct {
		Name                       string
		ParamId                    uint
		CategoryMockReturnError    error
		CategoryMockReturnBody     *sqlmock.Rows
		SubCategoryMockReturnError error
		SubCategoryMockReturnBody  *sqlmock.Rows
		HasReturnBody              bool
		ExpectedBody               dto.Category
	}{
		{
			"success",
			1,
			nil,
			sqlmock.NewRows([]string{"id", "name", "sub_categories"}).
				AddRow(1, "Expense", nil),
			nil,
			sqlmock.NewRows([]string{"id", "name", "category_id", "user_id"}).
				AddRow(1, "Food", 1, &userId),
			true,
			dto.Category{
				ID:   1,
				Name: "Expense",
				SubCategories: []dto.SubCategory{
					{
						ID:         1,
						CategoryID: 1,
						UserID:     &userId,
						Name:       "Food",
					},
				},
			},
		},
		{
			"categpry not found",
			1,
			nil,
			sqlmock.NewRows([]string{"id", "name", "sub_categories"}),
			nil,
			sqlmock.NewRows([]string{"id", "name", "category_id", "user_id"}),
			false,
			dto.Category{},
		},
		{
			"fail get category",
			1,
			errors.New("error"),
			sqlmock.NewRows([]string{"id", "name", "sub_categories"}).
				AddRow(1, "Expense", nil),
			errors.New("error"),
			sqlmock.NewRows([]string{"id", "name", "category_id", "user_id"}).
				AddRow(1, "Food", 1, &userId),
			false,
			dto.Category{},
		},
	}

	for _, v := range testCase {
		if v.CategoryMockReturnError != nil {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT `categories`.`id`,`categories`.`name` FROM `categories` WHERE id = ? AND `categories`.`deleted_at` IS NULL")).
				WithArgs(1).WillReturnError(v.CategoryMockReturnError)
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `sub_categories` WHERE `sub_categories`.`category_id` = ?")).
				WithArgs(1).WillReturnError(v.SubCategoryMockReturnError)
		} else {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT `categories`.`id`,`categories`.`name` FROM `categories` WHERE id = ? AND `categories`.`deleted_at` IS NULL")).
				WithArgs(1).WillReturnRows(v.CategoryMockReturnBody)
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `sub_categories` WHERE `sub_categories`.`category_id` = ?")).
				WithArgs(1).WillReturnRows(v.SubCategoryMockReturnBody)
		}
		s.T().Run(v.Name, func(t *testing.T) {
			user, err := s.categoryRepository.GetCategoryByID(v.ParamId)
			if v.HasReturnBody {
				s.Equal(v.ExpectedBody, user)
				s.NoError(err)
			} else {
				s.Error(err)
				s.Empty(user)
			}
		})
	}
}

func (s *suiteCategory) TestGetAllCategory() {
	userId := uint(1)
	testCase := []struct {
		Name                       string
		ParamId                    uint
		CategoryMockReturnError    error
		CategoryMockReturnBody     *sqlmock.Rows
		SubCategoryMockReturnError error
		SubCategoryMockReturnBody  *sqlmock.Rows
		HasReturnBody              bool
		ExpectedBody               []dto.Category
	}{
		{
			"success",
			1,
			nil,
			sqlmock.NewRows([]string{"id", "name", "sub_categories"}).
				AddRow(1, "Expense", nil),
			nil,
			sqlmock.NewRows([]string{"id", "name", "category_id", "user_id"}).
				AddRow(1, "Food", 1, &userId),
			true,
			[]dto.Category{
				{
					ID:   1,
					Name: "Expense",
					SubCategories: []dto.SubCategory{
						{
							ID:         1,
							CategoryID: 1,
							UserID:     &userId,
							Name:       "Food",
						},
					},
				},
			},
		},
		// {
		// 	"categpry not found",
		// 	1,
		// 	nil,
		// 	sqlmock.NewRows([]string{"id", "name", "sub_categories"}),
		// 	nil,
		// 	sqlmock.NewRows([]string{"id", "name", "category_id", "user_id"}),
		// 	false,
		// 	[]dto.Category{},
		// },
		// {
		// 	"fail get category",
		// 	1,
		// 	errors.New("error"),
		// 	sqlmock.NewRows([]string{"id", "name", "sub_categories"}).
		// 	AddRow(1, "Expense", nil),
		// 	errors.New("error"),
		// 	sqlmock.NewRows([]string{"id", "name", "category_id", "user_id"}).
		// 		AddRow(1, "Food", 1, &userId),
		// 	false,
		// 	[]dto.Category{},
		// },

	}

	for _, v := range testCase {
		if v.CategoryMockReturnError != nil {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT `categories`.`id`,`categories`.`name` FROM `categories` WHERE `categories`.`deleted_at` IS NULL")).
				WillReturnError(v.CategoryMockReturnError)
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `sub_categories` WHERE `sub_categories`.`category_id` = ?")).
				WillReturnError(v.SubCategoryMockReturnError)
		} else {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT `categories`.`id`,`categories`.`name` FROM `categories` WHERE `categories`.`deleted_at` IS NULL")).
				WillReturnRows(v.CategoryMockReturnBody)
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `sub_categories` WHERE `sub_categories`.`category_id` = ?")).
				WillReturnRows(v.SubCategoryMockReturnBody)
		}
		s.T().Run(v.Name, func(t *testing.T) {
			user, err := s.categoryRepository.GetAllCategory()
			if v.HasReturnBody {
				s.Equal(v.ExpectedBody, user)
				s.NoError(err)
			} else {
				s.Error(err)
				s.Empty(user)
			}
		})
	}
}

func (s *suiteCategory) TearDownSuite() {
	s.mock = nil
}

func TestSuiteCategorys(t *testing.T) {
	suite.Run(t, new(suiteCategory))
}
