package dynamodb

import (
	"errors"
	"slack-pr-lambda/env"
	"slack-pr-lambda/types"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func DynamoDbConnection() *dynamodb.DynamoDB {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	db := env.GetEnv("DB_ENDPOINT", "http://localhost:8000")
	region := env.GetEnv("REGION", "us-east-1")

	svc := dynamodb.New(sess, &aws.Config{
		Endpoint: &db,
		Region:   &region,
	})

	return svc
}

func InsertItem(svc *dynamodb.DynamoDB, item *types.TablePullRequestData) error {
	tableName := env.GetEnv("TABLE_NAME", "PullRequests")

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	insert := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(insert)
	if err != nil {
		return err
	}

	return nil
}

func GetSlackTimeStamp(svc *dynamodb.DynamoDB, id int, pullRequestId int) (string, error) {
	tableName := env.GetEnv("TABLE_NAME", "PullRequests")

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(strconv.Itoa(id)),
			},
			"pullRequestId": {
				N: aws.String(strconv.Itoa(pullRequestId)),
			},
		},
	})
	if err != nil {
		return "", err
	}
	if result.Item == nil {
		return "", errors.New("no data found")
	}

	item := types.TablePullRequestData{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return "", err
	}

	return item.SlackTimeStamp, nil
}

func DeleteItem(svc *dynamodb.DynamoDB, id int, pullRequestId int) error {
	tableName := env.GetEnv("TABLE_NAME", "PullRequests")

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(strconv.Itoa(id)),
			},
			"pullRequestId": {
				N: aws.String(strconv.Itoa(pullRequestId)),
			},
		},
		TableName: aws.String(tableName),
	}

	if _, err := svc.DeleteItem(input); err != nil {
		return err
	}
	return nil
}
