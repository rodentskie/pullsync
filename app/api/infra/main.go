package main

import (
	"slack-pr-lambda/api/infra/dynamodb"
	"slack-pr-lambda/api/infra/lambda"
	lambdaiamrole "slack-pr-lambda/api/infra/lambda_iam_role"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		if err := lambdaiamrole.LambdaIamRole(ctx); err != nil {
			return err
		}
		if err := lambda.LambdaFunction(ctx); err != nil {
			return err
		}

		if err := dynamodb.DynamoDB(ctx); err != nil {
			return err
		}
		return nil
	})
}
