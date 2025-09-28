package database

import (
	"fmt"
	types "lambda/type"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const Table_Name = "userTable"

type UserStore interface {
	DoesUserExist(username string) (bool, error)
	InsertUser(user types.User) error
	GetUser(username string) (types.User , error)
}

type DynamoDBClient struct {
	databasestore *dynamodb.DynamoDB
}

func NewDynamoDBClient() *DynamoDBClient {
	dbSession := session.Must(session.NewSession())
	db := dynamodb.New(dbSession)
	return &DynamoDBClient{
		databasestore: db,
	}
}


func (u DynamoDBClient) DoesUserExist(username string) (bool, error) {
	result, err := u.databasestore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(Table_Name),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	}) 
	if err != nil {
		return false, err
	}
	return result.Item != nil, nil
}

func (u DynamoDBClient) InsertUser (user types.User) error {
	item := &dynamodb.PutItemInput{
		TableName: aws.String(Table_Name),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
			"password": {
				S: aws.String(user.PasswordHash),
			},
		},
	}

	_, err := u.databasestore.PutItem(item)

	if err != nil {
		return err
	}
	return nil
}

func (u DynamoDBClient) GetUser(username string) (types.User , error) {
	var user types.User

	result, err := u.databasestore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(Table_Name),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return user, err
	}

	if result.Item == nil {
		return user, fmt.Errorf("user not found")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)

	if err != nil {
		return user, err
	}

	return user, nil
}