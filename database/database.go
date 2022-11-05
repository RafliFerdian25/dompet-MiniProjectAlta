package database

import (
	"dompet-miniprojectalta/config"
	"dompet-miniprojectalta/models/model"
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	cfg := config.Cfg

	username := cfg.DB_USERNAME
	password := cfg.DB_PASSWORD
	address := cfg.DB_ADDRESS
	dbName := cfg.DB_NAME

	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		username,
		password,
		address,
		dbName,
	)
	return gorm.Open(mysql.Open(connectionString), &gorm.Config{})
}

func MigrateDB(db *gorm.DB) error {
	db.Migrator().DropTable(
		&model.Category{},
		&model.User{},
		model.SubCategory{},
		model.Account{},
		model.Debt{},
		model.TransactionAccount{},
		model.Transaction{},
	)
	err := db.AutoMigrate(
		model.User{},
		model.Category{},
		model.SubCategory{},
		model.Account{},
		model.Debt{},
		model.TransactionAccount{},
		model.Transaction{},
	)
	// seeder Category
	category := []string{"Debt and Loan", "Expense", "Income"}
	if err == nil && db.Migrator().HasTable(&model.Category{}) {
		if err := db.First(&model.Category{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			for i, name := range category {
				db.Create(&model.Category{
					Model: gorm.Model{
						ID: uint(i) + 1,
					},
					Name: name,
				})
			}
		}
	}
	// seeder SubCategory
	subCategory := []string{"Debt collection", "Repayment", "Loan", "Debt collection", "Food", "Transportation", "Entertainment","Shopping","Healt & Fitness","Salary","Gift","Selling","Award","Interest money", "Bonus"}
	if err == nil && db.Migrator().HasTable(&model.SubCategory{}) {
		if err := db.First(&model.SubCategory{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			for i, name := range subCategory {
				var categoryID uint
				if i < 4 {
					categoryID = 1
				} else if i < 9 {
					categoryID = 2
				} else {
					categoryID = 3
				}
				err = db.Create(&model.SubCategory{
					Model: gorm.Model{
						ID: uint(i) + 1,
					},
					CategoryID: categoryID,
					Name: name,
				}).Error
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	if err != nil {
		return err
	}
	return nil
}
