package main

import (
	"dompet-miniprojectalta/config"
	"dompet-miniprojectalta/database"
	"dompet-miniprojectalta/routes"
	"log"
)

func main() {
	config.InitConfig()
	db, err := database.ConnectDB()
	if err != nil {
		panic(err)
	}
	err = database.MigrateDB(db)
	if err != nil {
		panic(err)
	}

	app := routes.New(db)

	apiPort := config.Cfg.API_PORT
	log.Fatal(app.Start(apiPort))
}