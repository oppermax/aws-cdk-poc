package main

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awss3"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type InfraStackProps struct {
	awscdk.StackProps
}

func createReplicationRule(destinationBucket awss3.CfnBucket) awss3.CfnBucket_ReplicationRuleProperty {
	return awss3.CfnBucket_ReplicationRuleProperty{
		Destination:             awss3.CfnBucket_ReplicationDestinationProperty{
			Bucket:                   destinationBucket.BucketName(),
			AccessControlTranslation: jsii.String("Destination"),
			Metrics:                  awss3.CfnBucket_MetricsProperty{
				Status:         jsii.String("Enabled"),
			},
			ReplicationTime:          awss3.CfnBucket_ReplicationTimeValueProperty{Minutes: jsii.Number(15)},
		},
		Status:                  jsii.String("Enabled"),
		DeleteMarkerReplication: awss3.CfnBucket_DeleteMarkerReplicationProperty{Status: jsii.String("Enabled")},
		Id:                      jsii.String("max-test-replication-rule"),
	}
}

func NewReplicationSourceStack(scope constructs.Construct, id string, props *InfraStackProps, destinationBucket awss3.CfnBucket) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here

	bucket := awss3.NewCfnBucket(stack, jsii.String("MaxTestReplicationSource"), &awss3.CfnBucketProps{
		BucketName:                       jsii.String("MaxTestReplicationSource"),
	})

	bucket.SetReplicationConfiguration(createReplicationRule(destinationBucket))

	return stack
}


func NewReplicationDestinationStack(scope constructs.Construct, id string, props *InfraStackProps) (awscdk.Stack, awss3.CfnBucket) {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here

	bucket := awss3.NewCfnBucket(stack, jsii.String("MaxTestReplicationDestination"), &awss3.CfnBucketProps{
		BucketName:                       jsii.String("MaxTestReplicationDestination"),
	})



	return stack, bucket
}
func main() {
	app := awscdk.NewApp(nil)

	_, destBucket  := NewReplicationDestinationStack(app, "MaxTestReplicationSourceStack", &InfraStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	NewReplicationSourceStack(app, "MaxTestReplicationDestinationStack", &InfraStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	}, destBucket)

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
