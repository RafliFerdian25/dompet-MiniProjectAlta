package subCategoryRepository

import (
	"database/sql/driver"
	"dompet-miniprojectalta/models/dto"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type suiteSubCategory struct {
	suite.Suite
	subCategoryRepository SubCategoryRepository
	mock                  sqlmock.Sqlmock
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func (s *suiteSubCategory) SetupSuite() {
	dbGormPalsu, mocking, err := sqlmock.New()
	s.NoError(err)

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      dbGormPalsu,
	}))

	s.mock = mocking

	NewSubCategoryRepository := NewSubCategoryRepository(dbGorm)
	s.subCategoryRepository = NewSubCategoryRepository
}

func (s *suiteSubCategory) TestGetSubCategoryById() {
	userId := uint(1)
	testCase := []struct {
		Name            string
		ParamId         uint
		MockReturnError error
		MockReturnBody  *sqlmock.Rows
		HasReturnBody   bool
		ExpectedBody    dto.SubCategoryDTO
	}{
		{
			"success",
			1,
			nil,
			sqlmock.NewRows([]string{"id", "category_id", "name", "user_id"}).
				AddRow(1, 1, "Food", 1),
			true,
			dto.SubCategoryDTO{
				ID:         1,
				CategoryID: 1,
				Name:       "Food",
				UserID:     &userId,
			},
		},
		{
			"failed get sub category by id",
			1,
			gorm.ErrRecordNotFound,
			sqlmock.NewRows([]string{"id", "category_id", "name", "user_id"}).
				AddRow(1, 1, "Food", 1),
			false,
			dto.SubCategoryDTO{},
		},
	}

	for _, v := range testCase {
		if v.MockReturnError != nil {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT `sub_categories`.`id`,`sub_categories`.`category_id`,`sub_categories`.`user_id`,`sub_categories`.`name` FROM `sub_categories` WHERE `sub_categories`.`id` = ? AND `sub_categories`.`deleted_at` IS NULL ORDER BY `sub_categories`.`id` LIMIT 1")).
				WithArgs(1).WillReturnError(v.MockReturnError)
		} else {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT `sub_categories`.`id`,`sub_categories`.`category_id`,`sub_categories`.`user_id`,`sub_categories`.`name` FROM `sub_categories` WHERE `sub_categories`.`id` = ? AND `sub_categories`.`deleted_at` IS NULL ORDER BY `sub_categories`.`id` LIMIT 1")).
				WithArgs(1).WillReturnRows(v.MockReturnBody)
		}
		s.T().Run(v.Name, func(t *testing.T) {
			user, err := s.subCategoryRepository.GetSubCategoryById(v.ParamId)
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

func (s *suiteSubCategory) TestGetSubCategoryByUser() {
	userId := uint(1)
	testCase := []struct {
		Name            string
		ParamUserId     uint
		MockReturnError error
		MockReturnBody  *sqlmock.Rows
		HasReturnBody   bool
		ExpectedBody    []dto.SubCategoryDTO
	}{
		{
			"success",
			1,
			nil,
			sqlmock.NewRows([]string{"id", "category_id", "name", "user_id"}).
				AddRow(1, 1, "Food", 1),
			true,
			[]dto.SubCategoryDTO{
				{
					ID:         1,
					CategoryID: 1,
					Name:       "Food",
					UserID:     &userId,
				},
			},
		},
		{
			"failed get sub category by id",
			1,
			gorm.ErrRecordNotFound,
			sqlmock.NewRows([]string{"id", "category_id", "name", "user_id"}).
				AddRow(1, 1, "Food", 1),
			false,
			[]dto.SubCategoryDTO{},
		},
	}

	for _, v := range testCase {
		if v.MockReturnError != nil {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT `sub_categories`.`id`,`sub_categories`.`category_id`,`sub_categories`.`user_id`,`sub_categories`.`name` FROM `sub_categories` WHERE (user_id = ? OR user_id IS NULL) AND `sub_categories`.`deleted_at` IS NULL")).
				WithArgs(1).WillReturnError(v.MockReturnError)
		} else {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT `sub_categories`.`id`,`sub_categories`.`category_id`,`sub_categories`.`user_id`,`sub_categories`.`name` FROM `sub_categories` WHERE (user_id = ? OR user_id IS NULL) AND `sub_categories`.`deleted_at` IS NULL")).
				WithArgs(1).WillReturnRows(v.MockReturnBody)
		}
		s.T().Run(v.Name, func(t *testing.T) {
			user, err := s.subCategoryRepository.GetSubCategoryByUser(v.ParamUserId)
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

func (s *suiteSubCategory) TearDownSuite() {
	s.mock = nil
}

func TestSuiteSubCategory(t *testing.T) {
	suite.Run(t, new(suiteSubCategory))
}
