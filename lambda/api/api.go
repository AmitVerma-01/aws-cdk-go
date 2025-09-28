package api

import (
	"encoding/json"
	"fmt"
	"lambda/database"
	types "lambda/type"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type APIHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbstore database.UserStore) APIHandler {
	return APIHandler{
		dbStore: dbstore,
	}
}

func (api APIHandler) RegisterUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registerUser types.RegisterUser

	err := json.Unmarshal([]byte(request.Body), &registerUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Data",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	if registerUser.Password == "" || registerUser.Username == "" {
		return events.APIGatewayProxyResponse{
			Body:       "username or password is empty",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	userExists, err := api.dbStore.DoesUserExist(registerUser.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "failed to check if user exists",
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	if userExists {
		return events.APIGatewayProxyResponse{
			Body:       "user already exists",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	user , err := types.NewUser(registerUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "failed to create user",
			StatusCode: http.StatusInternalServerError,
		}, nil
	}


	err = api.dbStore.InsertUser(user)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "failed to insert user",
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       "User registered successfully",
		StatusCode: http.StatusOK,
	}, nil
}


func (api APIHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	type LoginUser struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var loginUser LoginUser
	
	err := json.Unmarshal([]byte(request.Body), &loginUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Data",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	user , err := api.dbStore.GetUser(loginUser.Username)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "failed to get user",
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	if !types.ValidatedPassword(user.PasswordHash, loginUser.Password) {
		return events.APIGatewayProxyResponse{
			Body : "Invalid Password",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	accessToken := types.CreateToken(user)
	if accessToken == "" {
		return events.APIGatewayProxyResponse{
			Body:       "failed to create access token",
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	successMessaege := fmt.Sprintf("{\"accessToken\": \"%s\"}", accessToken) 

	return events.APIGatewayProxyResponse{
		Body:       successMessaege,
		StatusCode: http.StatusOK,
	}, nil
}