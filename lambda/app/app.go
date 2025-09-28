package app

import (
	"lambda/api"
	"lambda/database"
)

type App struct {
	ApiHandler api.APIHandler
}

func NewApp() App {

	db := database.NewDynamoDBClient()
	allStore := database.NewAllStore(db, db)
	apiHandler := api.NewApiHandler(allStore)

	return App{
		ApiHandler: apiHandler,
	}
}