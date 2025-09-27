package app

import (
	"lambda/api"
	"lambda/database"
)

type App struct {
	apiHandler api.APIHandler
}

func NewApp() App {

	db := database.NewDynamoDBClient()
	apiHandler := api.NewApiHandler(*db)

	return App{
		apiHandler: apiHandler,
	}
}