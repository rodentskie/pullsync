package dynamodb

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func DynamoDB(ctx *pulumi.Context) error {
	conf := config.New(ctx, "")
	region := conf.Require("region")
	env := conf.Require("env")
	tableName := conf.Require("tableName")
	tableNameIndex := conf.Require("tableNameIndex")

	_, err := dynamodb.NewTable(ctx, "pr_table", &dynamodb.TableArgs{
		Name:          pulumi.String(tableName),
		BillingMode:   pulumi.String("PROVISIONED"),
		ReadCapacity:  pulumi.Int(5),
		WriteCapacity: pulumi.Int(5),
		HashKey:       pulumi.String("id"),
		RangeKey:      pulumi.String("pullRequestId"),
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("id"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("pullRequestId"),
				Type: pulumi.String("N"),
			},
		},
		GlobalSecondaryIndexes: dynamodb.TableGlobalSecondaryIndexArray{
			&dynamodb.TableGlobalSecondaryIndexArgs{
				Name:           pulumi.String(tableNameIndex),
				HashKey:        pulumi.String("pullRequestId"),
				WriteCapacity:  pulumi.Int(5),
				ReadCapacity:   pulumi.Int(5),
				ProjectionType: pulumi.String("ALL"),
			},
		},
		Tags: pulumi.StringMap{
			"Region":      pulumi.String(region),
			"Environment": pulumi.String(env),
			"TableName":   pulumi.String(tableName),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
