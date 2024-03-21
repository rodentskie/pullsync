package lambdaiamrole

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func LambdaIamRole(ctx *pulumi.Context) (*iam.Role, error) {
	conf := config.New(ctx, "")
	lambdaBasicExecRoleArn := conf.Require("lambdaBasicExecRoleArn")
	lambdaDynamoDBExecRoleArn := conf.Require("lambdaDynamoDBExecRoleArn")
	lambdaRoleName := conf.Require("lambdaRoleName")

	assumeRole, err := iam.GetPolicyDocument(ctx, &iam.GetPolicyDocumentArgs{
		Statements: []iam.GetPolicyDocumentStatement{
			{
				Effect: pulumi.StringRef("Allow"),
				Principals: []iam.GetPolicyDocumentStatementPrincipal{
					{
						Type: "Service",
						Identifiers: []string{
							"lambda.amazonaws.com",
						},
					},
				},
				Actions: []string{
					"sts:AssumeRole",
				},
			},
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	allow := "Allow"
	inlinePolicy, err := iam.GetPolicyDocument(ctx, &iam.GetPolicyDocumentArgs{
		Statements: []iam.GetPolicyDocumentStatement{
			{
				Actions: []string{
					"dynamodb:BatchGetItem",
					"dynamodb:GetItem",
					"dynamodb:Query",
					"dynamodb:Scan",
					"dynamodb:BatchWriteItem",
					"dynamodb:PutItem",
					"dynamodb:UpdateItem",
					"dynamodb:DeleteItem",
				},
				Resources: []string{
					"*",
				},
				Effect: &allow,
			},
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	role, err := iam.NewRole(ctx, "slack_pr_lambda", &iam.RoleArgs{
		Name:             pulumi.String(lambdaRoleName),
		AssumeRolePolicy: pulumi.String(assumeRole.Json),
		ManagedPolicyArns: pulumi.StringArray{
			pulumi.String(lambdaBasicExecRoleArn),
			pulumi.String(lambdaDynamoDBExecRoleArn),
		},
		InlinePolicies: iam.RoleInlinePolicyArray{
			&iam.RoleInlinePolicyArgs{
				Name:   pulumi.String("lambda-dynamodb-crud"),
				Policy: pulumi.String(inlinePolicy.Json),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return role, nil
}
