package main

import (
	"fmt"
	"lambda/app"
	"lambda/middleware"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(event MyEvent) (string, error) {
	if event.Name == "" {
		return "" , fmt.Errorf("Name is required")
	}
	return fmt.Sprintf("Hello, %s!\n", event.Name), nil
}

func main() {
	myApp := app.NewApp()
	// lambda.Start(myApp.ApiHandler.RegisterUser)
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
		switch request.Path {
			case "/register":
				return myApp.ApiHandler.RegisterUser(request)
			case "/login":
				return myApp.ApiHandler.LoginUser(request)
			case "/protected":
				return middleware.ValidateJWTMiddleware(middleware.Protectedhandler)(request)
			default: 
				return events.APIGatewayProxyResponse{ 
					Body:  "Not Found",
					StatusCode: http.StatusNotFound,
				}, nil
		}
	})
}