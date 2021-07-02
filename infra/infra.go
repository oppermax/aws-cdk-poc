package main

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awss3"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

const (
	sourceBucketName      = "max-test-replication-source-bucket"
	destinationBucketName = "max-test-replication-destination-bucket"
)

func getAccountId() *string {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return nil
	}

	svc := sts.New(sess)
	var params *sts.GetCallerIdentityInput
	resp, err := svc.GetCallerIdentity(params)
	if err != nil {
		fmt.Println("failed to get caller identity", err)
	}

	fmt.Printf("Deploying into account: %s\n", *resp.Account)
	return resp.Account
}

type InfraStackProps struct {
	awscdk.StackProps
}

type replicationTester struct {
	destinationBucket awss3.CfnBucket
	sourceBucket      awss3.CfnBucket
	replicatioRole    awsiam.Role
	stack             awscdk.Stack
}

func (r *replicationTester) createReplicationRule() awss3.CfnBucket_ReplicationRuleProperty {
	dest := awss3.CfnBucket_ReplicationDestinationProperty{
		Bucket:                   r.destinationBucket.AttrArn(),
		AccessControlTranslation: awss3.CfnBucket_AccessControlTranslationProperty{Owner: jsii.String("Destination")},
		Metrics: awss3.CfnBucket_MetricsProperty{
			Status:         jsii.String("Enabled"),
			EventThreshold: awss3.CfnBucket_ReplicationTimeValueProperty{Minutes: jsii.Number(15)},
		},
		ReplicationTime: awss3.CfnBucket_ReplicationTimeProperty{
			Status: jsii.String("Enabled"),
			Time:   awss3.CfnBucket_ReplicationTimeValueProperty{Minutes: jsii.Number(15)},
		},
		Account: getAccountId(),
	}

	return awss3.CfnBucket_ReplicationRuleProperty{
		Destination: dest,
		Status:      jsii.String("Enabled"),
		Id:          jsii.String("max-test-replication-rule"),
		Filter: awss3.CfnBucket_ReplicationRuleFilterProperty{
			Prefix:    jsii.String(""),
		},
		DeleteMarkerReplication: awss3.CfnBucket_DeleteMarkerReplicationProperty{Status: jsii.String("Disabled")},
		Priority:                jsii.Number(1),
	}
}

func (r *replicationTester) CreateSourceComponents() {

	r.sourceBucket = awss3.NewCfnBucket(r.stack, jsii.String("MaxTestReplicationSourceBucket"), &awss3.CfnBucketProps{
		BucketName:              jsii.String(sourceBucketName),
		VersioningConfiguration: awss3.CfnBucket_VersioningConfigurationProperty{Status: jsii.String("Enabled")},
		ReplicationConfiguration: awss3.CfnBucket_ReplicationConfigurationProperty{
			Role:  r.replicatioRole.RoleArn(),
			Rules: []awss3.CfnBucket_ReplicationRuleProperty{r.createReplicationRule()},
		},
	})

}

func (r *replicationTester) CreateDestinationComponents() {

	r.destinationBucket = awss3.NewCfnBucket(r.stack, jsii.String("MaxTestReplicationDestinationBucket"), &awss3.CfnBucketProps{
		BucketName:              jsii.String(destinationBucketName),
		VersioningConfiguration: awss3.CfnBucket_VersioningConfigurationProperty{Status: jsii.String("Enabled")},
	})

}

func (r *replicationTester) addPoliciesToReplicationRole() {
	policy := awsiam.NewPolicyDocument(nil)

	policy.AddStatements(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   &[]*string{jsii.String("s3:ReplicateObject"), jsii.String("s3:ReplicateDelete"), jsii.String("s3:ReplicateTags")},
		Effect:    awsiam.Effect_ALLOW,
		Resources: &[]*string{jsii.String(*r.destinationBucket.AttrArn() + "/*")},
	}))

	policy.AddStatements(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   &[]*string{jsii.String("s3:GetReplicationConfiguration"), jsii.String("s3:ListBucket")},
		Effect:    awsiam.Effect_ALLOW,
		Resources: &[]*string{r.sourceBucket.AttrArn()},
	}))

	policy.AddStatements(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   &[]*string{jsii.String("s3:GetObjectVersionForReplication"), jsii.String("s3:GetObjectVersionAcl"), jsii.String("s3:GetObjectVersionTagging")},
		Effect:    awsiam.Effect_ALLOW,
		Resources: &[]*string{jsii.String(*r.sourceBucket.AttrArn() + "/*")},
	}))

	policy.AddStatements(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   &[]*string{jsii.String("s3:ObjectOwnerOverrideToBucketOwner")},
		Effect:    awsiam.Effect_ALLOW,
		Resources: &[]*string{jsii.String(*r.destinationBucket.AttrArn() + "/*")},
	}))

	r.replicatioRole.AttachInlinePolicy(awsiam.NewPolicy(r.stack, jsii.String("max-test-replication-policy"), &awsiam.PolicyProps{
		Document: policy,
	}))
}

//func (r *replicationTester) addStatementsToBucketPolicy() {
//	policy := awsiam.NewPolicyDocument(nil)
//
//	policy.AddStatements(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
//		Actions:   &[]*string{jsii.String("s3:ReplicateObject"), jsii.String("s3:ReplicateDelete"), jsii.String("s3:ReplicateTags")},
//		Effect:    awsiam.Effect_ALLOW,
//		Resources: &[]*string{r.destinationBucket.AttrArn()},
//	}))
//
//	policy.AddStatements(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
//		Actions:   &[]*string{jsii.String("s3:GetReplicationConfiguration"), jsii.String("s3:ListBucket")},
//		Effect:    awsiam.Effect_ALLOW,
//		Resources: &[]*string{r.sourceBucket.AttrArn()},
//	}))
//
//	policy.AddStatements(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
//		Actions:   &[]*string{jsii.String("s3:GetObjectVersionForReplication"), jsii.String("s3:GetObjectVersionAcl"), jsii.String("s3:GetObjectVersionTagging")},
//		Effect:    awsiam.Effect_ALLOW,
//		Resources: &[]*string{jsii.String(*r.destinationBucket.AttrArn() + "/*")},
//	}))
//
//	r.replicatioRole.AttachInlinePolicy(awsiam.NewPolicy(r.stack, jsii.String("max-test-replication-policy"), &awsiam.PolicyProps{
//		Document: policy,
//	}))
//}


func (r *replicationTester) CreateIamRoleComponents() {

	r.replicatioRole = awsiam.NewRole(r.stack, jsii.String("max-test-replication-role"), &awsiam.RoleProps{
		AssumedBy:      awsiam.NewServicePrincipal(jsii.String("s3.amazonaws.com"), nil),
		Description:    jsii.String("test role for replication oppermax"),
		InlinePolicies: nil,
		RoleName:       jsii.String("max-test-replication-role"),
		Path: jsii.String("/"),
	})

}

func (r *replicationTester) NewReplicationTestStack(scope constructs.Construct, id string, props *InfraStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	r.stack = awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here

	r.CreateDestinationComponents()
	r.CreateIamRoleComponents()
	r.CreateSourceComponents()
	r.addPoliciesToReplicationRole()

	return r.stack
}

func main() {
	app := awscdk.NewApp(nil)

	r := replicationTester{}

	r.NewReplicationTestStack(app, "max-test-replication-stack", &InfraStackProps{awscdk.StackProps{
		Env:       env(),
		StackName: jsii.String("max-test-replication-stack"),
	}})

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
