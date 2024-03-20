package dynamodb

import (
	"fmt"
	"slack-pr-lambda/types"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

func TestDynamoDbConnection(t *testing.T) {
	envVars := map[string]string{
		"DB_ENDPOINT": "http://localhost:8000",
		"REGION":      "us-east-1",
	}

	for key, value := range envVars {
		t.Setenv(key, value)
	}

	svc := DynamoDbConnection()

	assert.NotNil(t, svc)
	assert.Equal(t, envVars["DB_ENDPOINT"], aws.StringValue(&svc.Endpoint))
	assert.Equal(t, envVars["REGION"], aws.StringValue(svc.Config.Region))
}

func TestInsertItem(t *testing.T) {
	envVars := map[string]string{
		"TABLE_NAME": "PullRequests",
	}

	for key, value := range envVars {
		t.Setenv(key, value)
	}

	svc := DynamoDbConnection()

	t.Run("successful", func(t *testing.T) {
		item := &types.TablePullRequestData{
			ID:             fmt.Sprintf("%d", time.Now().UnixMilli()),
			PullRequestId:  int(time.Now().UnixMilli()),
			SlackTimeStamp: fmt.Sprintf("%d", time.Now().UnixMilli()),
		}

		err := InsertItem(svc, item)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		item := &types.TablePullRequestData{}

		err := InsertItem(svc, item)
		assert.Error(t, err)
	})

	if err := DeleteAllItem(svc); err != nil {
		t.Errorf("error delete all item %v", err)
	}
}

func TestGetSlackTimeStamp(t *testing.T) {
	envVars := map[string]string{
		"TABLE_NAME": "PullRequests",
	}

	for key, value := range envVars {
		t.Setenv(key, value)
	}

	svc := DynamoDbConnection()

	t.Run("successful", func(t *testing.T) {
		item := &types.TablePullRequestData{
			ID:             fmt.Sprintf("%d", time.Now().UnixMilli()),
			PullRequestId:  int(time.Now().UnixMilli()),
			SlackTimeStamp: fmt.Sprintf("%d", time.Now().UnixMilli()),
		}

		err := InsertItem(svc, item)
		assert.NoError(t, err)

		id, err := strconv.Atoi(item.ID)
		assert.NoError(t, err)

		timeStamp, err := GetSlackTimeStamp(svc, id, item.PullRequestId)
		assert.NotNil(t, timeStamp)
		assert.NoError(t, err)

	})

	t.Run("empty", func(t *testing.T) {
		item := &types.TablePullRequestData{
			ID:             fmt.Sprintf("%d", time.Now().UnixMilli()),
			PullRequestId:  int(time.Now().UnixMilli()),
			SlackTimeStamp: fmt.Sprintf("%d", time.Now().UnixMilli()),
		}

		err := InsertItem(svc, item)
		assert.NoError(t, err)

		timeStamp, err := GetSlackTimeStamp(svc, int(time.Now().UnixMilli()), item.PullRequestId)
		assert.Equal(t, timeStamp, "")
		assert.Error(t, err)

	})

	if err := DeleteAllItem(svc); err != nil {
		t.Errorf("error delete all item %v", err)
	}
}

func TestDeleteItem(t *testing.T) {
	envVars := map[string]string{
		"TABLE_NAME": "PullRequests",
	}

	for key, value := range envVars {
		t.Setenv(key, value)
	}

	svc := DynamoDbConnection()

	t.Run("successful", func(t *testing.T) {
		item := &types.TablePullRequestData{
			ID:             fmt.Sprintf("%d", time.Now().UnixMilli()),
			PullRequestId:  int(time.Now().UnixMilli()),
			SlackTimeStamp: fmt.Sprintf("%d", time.Now().UnixMilli()),
		}

		err := InsertItem(svc, item)
		assert.NoError(t, err)

		id, err := strconv.Atoi(item.ID)
		assert.NoError(t, err)

		err = DeleteItem(svc, id, item.PullRequestId)
		assert.NoError(t, err)

	})

	if err := DeleteAllItem(svc); err != nil {
		t.Errorf("error delete all item %v", err)
	}
}

func TestDeleteAllItem(t *testing.T) {
	envVars := map[string]string{
		"TABLE_NAME": "PullRequests",
	}

	for key, value := range envVars {
		t.Setenv(key, value)
	}

	svc := DynamoDbConnection()

	t.Run("successful", func(t *testing.T) {
		item := &types.TablePullRequestData{
			ID:             fmt.Sprintf("%d", time.Now().UnixMilli()),
			PullRequestId:  int(time.Now().UnixMilli()),
			SlackTimeStamp: fmt.Sprintf("%d", time.Now().UnixMilli()),
		}

		err := InsertItem(svc, item)
		assert.NoError(t, err)

		if err := DeleteAllItem(svc); err != nil {
			t.Errorf("error delete all item %v", err)
		}

	})

}
