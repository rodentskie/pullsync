package lambdaiamrole

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func LambdaIamRole(ctx *pulumi.Context) error {
	conf := config.New(ctx, "")
	lambdaBasicExecRoleArn := conf.Require("lambdaBasicExecRoleArn")
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
		return err
	}

	_, err = iam.NewRole(ctx, "slack_pr_lambda", &iam.RoleArgs{
		Name:             pulumi.String(lambdaRoleName),
		AssumeRolePolicy: pulumi.String(assumeRole.Json),
		ManagedPolicyArns: pulumi.StringArray{
			pulumi.String(lambdaBasicExecRoleArn),
		},
	})

	if err != nil {
		return err
	}

	return nil
}
