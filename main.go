package main

import (
	"dompet-miniprojectalta/config"
	"dompet-miniprojectalta/database"
	"dompet-miniprojectalta/routes"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		panic(err)
	}
	err = database.MigrateDB(db)
	if err != nil {
		panic(err)
	}

	app := routes.New(db)


	apiPort := config.ConfigValue("API_PORT")
	app.Start(apiPort)
}