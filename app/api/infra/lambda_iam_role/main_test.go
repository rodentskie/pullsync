package lambdaiamrole

import (
	"slack-pr-lambda/pulumimock"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

func TestLambdaIamRole(t *testing.T) {
	config := map[string]string{
		"project:lambdaBasicExecRoleArn":    "testBasicExecRoleArn",
		"project:lambdaDynamoDBExecRoleArn": "testDynamoDBExecRoleArn",
		"project:lambdaRoleName":            "testRoleName",
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		err := LambdaIamRole(ctx)
		assert.NoError(t, err)

		return nil
	}, pulumimock.WithMocksAndConfig("project", "stack", config, pulumimock.Mocks(0)))
	assert.NoError(t, err)
}
