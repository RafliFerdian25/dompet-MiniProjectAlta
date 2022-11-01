package database

import (
	"dompet-miniprojectalta/config"
	"dompet-miniprojectalta/models/model"
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
	// db.Migrator().DropTable(
	// 	&model.Category{},
	// 	&model.User{},
	// 	model.SubCategory{},
	// 	model.Account{},
	// 	model.Debt{},
	// 	model.TransactionAccount{},
	// 	model.Transaction{},
	// )
	return db.AutoMigrate(
		model.User{},
		model.Category{},
		model.SubCategory{},
		model.Account{},
		model.Debt{},
		model.TransactionAccount{},
		model.Transaction{},
	)
}
