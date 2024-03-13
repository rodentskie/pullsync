package api

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func ApiGateway(ctx *pulumi.Context) error {
	conf := config.New(ctx, "")
	// region := conf.Require("region")
	// accountId := conf.Require("accountId")
	lambdaFunctionName := conf.Require("lambdaFunctionName")
	env := conf.Require("env")

	// API Gateway
	api, err := apigateway.NewRestApi(ctx, "api", &apigateway.RestApiArgs{
		Name: pulumi.String("slack-pr-api"),
	})
	if err != nil {
		return err
	}
	resource, err := apigateway.NewResource(ctx, "resource", &apigateway.ResourceArgs{
		PathPart: pulumi.String("resource"),
		ParentId: api.RootResourceId,
		RestApi:  api.ID(),
	})
	if err != nil {
		return err
	}

	lambdaFn, err := lambda.LookupFunction(ctx, &lambda.LookupFunctionArgs{
		FunctionName: lambdaFunctionName,
	}, nil)
	if err != nil {
		return err
	}

	methodGet, err := apigateway.NewMethod(ctx, "method_get", &apigateway.MethodArgs{
		RestApi:       api.ID(),
		ResourceId:    resource.ID(),
		HttpMethod:    pulumi.String("GET"),
		Authorization: pulumi.String("NONE"),
	})
	if err != nil {
		return err
	}

	methodPost, err := apigateway.NewMethod(ctx, "method_post", &apigateway.MethodArgs{
		RestApi:       api.ID(),
		ResourceId:    resource.ID(),
		HttpMethod:    pulumi.String("POST"),
		Authorization: pulumi.String("NONE"),
	})
	if err != nil {
		return err
	}

	_, err = apigateway.NewIntegration(ctx, "integration_get", &apigateway.IntegrationArgs{
		RestApi:               api.ID(),
		ResourceId:            resource.ID(),
		HttpMethod:            methodGet.HttpMethod,
		IntegrationHttpMethod: pulumi.String("POST"),
		Type:                  pulumi.String("AWS_PROXY"),
		Uri:                   pulumi.String(lambdaFn.InvokeArn),
	})

	if err != nil {
		return err
	}

	_, err = apigateway.NewIntegration(ctx, "integration_post", &apigateway.IntegrationArgs{
		RestApi:               api.ID(),
		ResourceId:            resource.ID(),
		HttpMethod:            methodPost.HttpMethod,
		IntegrationHttpMethod: pulumi.String("POST"),
		Type:                  pulumi.String("AWS_PROXY"),
		Uri:                   pulumi.String(lambdaFn.InvokeArn),
	})

	if err != nil {
		return err
	}

	_, err = lambda.NewPermission(ctx, "apigw_lambda", &lambda.PermissionArgs{
		StatementId: pulumi.String("AllowExecutionFromAPIGateway"),
		Action:      pulumi.String("lambda:InvokeFunction"),
		Function:    pulumi.String(lambdaFn.FunctionName),
		Principal:   pulumi.String("apigateway.amazonaws.com"),
		SourceArn: api.ExecutionArn.ApplyT(func(executionArn string) (string, error) {
			return fmt.Sprintf("%v/*", executionArn), nil
		}).(pulumi.StringOutput),
	})
	if err != nil {
		return err
	}

	_, err = apigateway.NewDeployment(ctx, "apigw_deploy", &apigateway.DeploymentArgs{
		RestApi:          api.ID(),
		StageName:        pulumi.String(env),
		Description:      pulumi.String("My API Gateway deployment"),
		StageDescription: pulumi.String(env + " stage"),
	})
	if err != nil {
		return err
	}

	_, err = cloudwatch.NewLogGroup(ctx, "apigw_cw", &cloudwatch.LogGroupArgs{
		Name: api.ID().ApplyT(func(id string) (string, error) {
			return fmt.Sprintf("API-Gateway-Execution-Logs_%v/%v", id, env), nil
		}).(pulumi.StringOutput),
		RetentionInDays: pulumi.Int(7),
	})
	if err != nil {
		return err
	}

	return nil
}
