package database

import (
	"dompet-miniprojectalta/config"
	"dompet-miniprojectalta/models/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {

	username := config.ConfigValue("DB_USERNAME")
	password := config.ConfigValue("DB_PASSWORD")
	address := config.ConfigValue("DB_ADDRESS")
	dbName := config.ConfigValue("DB_NAME")

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
