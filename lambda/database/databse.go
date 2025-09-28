package database

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const User_Table_Name = "userTable"
	
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

type AllStore struct {
	UserStore UserStore
	PostStore PostStore
}

func NewAllStore(UserStore UserStore, PostStore PostStore) AllStore {
	return AllStore{
		UserStore: UserStore,
		PostStore: PostStore,
	}
}

// func (allstore AllStore) GetUserStore() UserStore {
// 	return allstore.UserStore
// }

// func (allstore AllStore) GetPostStore() PostStore {
// 	return allstore.PostStore
// }