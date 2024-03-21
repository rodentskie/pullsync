package lambda

import (
	"slack-pr-lambda/pulumimock"
	"testing"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

func TestLambdaFunction(t *testing.T) {
	config := map[string]string{
		"project:lambdaRoleName":     "testRoleName",
		"project:lambdaFunctionName": "testLambdaFunctionName",
		"project:slackToken":         "testToken",
		"project:slackChannel":       "testChannel",
		"project:env":                "test",
		"project:dbEndpoint":         "testEndpoint",
		"project:region":             "ap-southeast-2",
		"project:githubOwner":        "foo",
		"project:githubToken":        "bar",
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		role := &iam.Role{
			Arn: pulumi.Sprintf("%s", "fakeArn"),
		}

		err := LambdaFunction(ctx, role)
		assert.NoError(t, err)

		return nil
	}, pulumimock.WithMocksAndConfig("project", "stack", config, pulumimock.Mocks(0)))
	assert.NoError(t, err)
}
