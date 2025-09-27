package main

import (
	"errors"
	"fmt"
	"lambda/app"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(event MyEvent) (string, error) {
	if event.Name == "" {
		return "" , errors.New("Name is required")
	}
	return fmt.Sprintf("Hello, %s!\n", event.Name), nil
}

func main() {
	app.NewApp()
	lambda.Start(HandleRequest)
}