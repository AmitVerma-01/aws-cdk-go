package api

import "lambda/database"


type APIHandler struct {
	dbStore database.DynamoDBClient
}


func NewApiHandler(dbstore database.DynamoDBClient) APIHandler {
	return APIHandler{
		dbStore: dbstore,
	}	
}