package database

import (
	"fmt"
	"lambda/type"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

const Post_Table_Name = "postTable"

type PostStore interface {
	InsertPost(post types.Post) error
	GetPost(postId string) (types.Post, error)
	GetAllPosts() ([]types.Post, error)
	DeletePost(postId string) error
	UpdatePost(post types.Post) error
}

func (u DynamoDBClient) InsertPost(post types.Post) error {
	postId := uuid.New().String()
	Items := &dynamodb.PutItemInput{
		TableName: aws.String(Post_Table_Name),
		Item: map[string]*dynamodb.AttributeValue{
			"postId" : {
				S: aws.String(postId),
			},
			"postContent" : {
				S: aws.String(post.PostContent),
			},
			"username" : {
				S: aws.String(*post.Username),
			},
			"createdAt" : {
				S: aws.String(post.CreatedAt.Format(time.RFC3339)),
			},
			"updatedAt" : {
				S: aws.String(post.UpdatedAt.Format(time.RFC3339)),
			},
		},
	}

	_, err := u.databasestore.PutItem(Items)
	if err != nil {
		return err
	}
	return nil
}


func (u DynamoDBClient) GetPost(postId string) (types.Post, error) {
	var post types.Post

	result, err := u.databasestore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(Post_Table_Name),
		Key: map[string]*dynamodb.AttributeValue{
			"postId": {
				S: aws.String(postId),
			},
		},
	})

	if err != nil {
		return post, err
	}

	if result.Item == nil {
		return post, fmt.Errorf("post not found")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &post)
	if err != nil {
		return post, err
	}
	return post, nil
}


func (u DynamoDBClient) GetAllPosts() ([]types.Post, error) {
	var posts []types.Post

	result, err := u.databasestore.Scan(&dynamodb.ScanInput{
		TableName: aws.String(Post_Table_Name),
	})

	if err != nil {
		return posts, err
	}

	for _, item := range result.Items {
		var post types.Post
		err = dynamodbattribute.UnmarshalMap(item, &post)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}


func (u DynamoDBClient) DeletePost(postId string) error {
	_, err := u.GetPost(postId)
	if err != nil {
		return err
	}

	deleteItem := &dynamodb.DeleteItemInput{
		TableName: aws.String(Post_Table_Name),
		Key: map[string]*dynamodb.AttributeValue{
			"postId": {
				S: aws.String(postId),
			},
		},
	}

	_, err = u.databasestore.DeleteItem(deleteItem)
	if err != nil {
		return err
	}
	return nil
}

func (u DynamoDBClient) UpdatePost(post types.Post) error {
	_, err := u.GetPost(*post.PostId)
	if err != nil {
		return err
	}

	updateItem := &dynamodb.UpdateItemInput{
		TableName: aws.String(Post_Table_Name),
		Key: map[string]*dynamodb.AttributeValue{
			"postId": {
				S: aws.String(*post.PostId),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":postContent": {
				S: aws.String(post.PostContent),
			},
			":updatedAt": {
				S: aws.String(post.UpdatedAt.Format(time.RFC3339)),
			},
		},
		UpdateExpression: aws.String("set postContent = :postContent, updatedAt = :updatedAt"),
	}

	_, err = u.databasestore.UpdateItem(updateItem)
	if err != nil {
		return err
	}
	return nil
}

		