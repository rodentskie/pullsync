package dynamodb

import (
	"slack-pr-lambda/pulumimock"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

func TestDynamoDB(t *testing.T) {
	config := map[string]string{
		"project:region":         "ap-southeast-2",
		"project:env":            "test",
		"project:tableName":      "testTable",
		"project:tableNameIndex": "testTableIndex",
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		err := DynamoDB(ctx)
		assert.NoError(t, err)

		return nil
	}, pulumimock.WithMocksAndConfig("project", "stack", config, pulumimock.Mocks(0)))
	assert.NoError(t, err)
}
