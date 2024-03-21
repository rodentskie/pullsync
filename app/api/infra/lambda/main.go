package lambda

import (
	"encoding/base64"
	"time"

	"github.com/pulumi/pulumi-aws-apigateway/sdk/v2/go/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func LambdaFunction(ctx *pulumi.Context, role *iam.Role) error {
	conf := config.New(ctx, "")
	lambdaFunctionName := conf.Require("lambdaFunctionName")
	slackToken := conf.Require("slackToken")
	slackChannel := conf.Require("slackChannel")
	env := conf.Require("env")
	dbEndpoint := conf.Require("dbEndpoint")
	region := conf.Require("region")
	githubOwner := conf.Require("githubOwner")
	githubToken := conf.Require("githubToken")

	// built zip file
	fileName := "../bin/bootstrap.zip"

	now := time.Now()
	dateTimeString := now.Format("2006-01-02T15:04:05Z07:00")
	dateTimeBytes := []byte(dateTimeString)
	base64EncodedHash := base64.StdEncoding.EncodeToString(dateTimeBytes)

	lambdaFn, err := lambda.NewFunction(ctx, "test_lambda", &lambda.FunctionArgs{
		Code:           pulumi.NewFileArchive(fileName),
		Name:           pulumi.String(lambdaFunctionName),
		Role:           role.Arn,
		Handler:        pulumi.String("bootstrap"),
		SourceCodeHash: pulumi.String(base64EncodedHash),
		Runtime:        pulumi.String("provided.al2023"),
		Environment: &lambda.FunctionEnvironmentArgs{
			Variables: pulumi.StringMap{
				"ENV":           pulumi.String(env),
				"SLACK_TOKEN":   pulumi.String(slackToken),
				"SLACK_CHANNEL": pulumi.String(slackChannel),
				"DB_ENDPOINT":   pulumi.String(dbEndpoint),
				"REGION":        pulumi.String(region),
				"GITHUB_TOKEN":  pulumi.String(githubToken),
				"GITHUB_OWNER":  pulumi.String(githubOwner),
			},
		},
		Tags: pulumi.StringMap{
			"Name": pulumi.String(lambdaFunctionName),
		},
	})

	if err != nil {
		return err
	}

	methodGet := apigateway.MethodGET
	methodPost := apigateway.MethodPOST
	_, err = apigateway.NewRestAPI(ctx, "api_slack_pr", &apigateway.RestAPIArgs{
		Routes: []apigateway.RouteArgs{
			{
				Path: "/", Method: &methodGet, EventHandler: lambdaFn,
			},
			{
				Path: "/pull-request", Method: &methodPost, EventHandler: lambdaFn,
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
