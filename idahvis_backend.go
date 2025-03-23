
package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsses"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type IdahvisBackendStackProps struct {
	awscdk.StackProps
}

func NewIdahvisBackendStack(scope constructs.Construct, id string, props *IdahvisBackendStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// DynamoDB Table Configuration
	table := awsdynamodb.NewTable(stack, jsii.String("idahvisDatabase"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("email"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName: jsii.String("message"),
	})

	// Lambda Function Configuration
	myFunction := awslambda.NewFunction(stack, jsii.String("idahvisFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("lambda/function.zip"), nil),
		Handler: jsii.String("main"),
	})

	// Grant Lambda permission to read/write to the DynamoDB table
	table.GrantReadWriteData(myFunction)

	// API Gateway with CORS Configuration
	api := awsapigateway.NewRestApi(stack, jsii.String("idahvisRestApi"), &awsapigateway.RestApiProps{
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowHeaders: jsii.Strings(
				"Content-Type",
				"X-Amz-Date",
				"Authorization",
				"X-Api-Key",
				"X-Amz-Security-Token",
			),
			AllowMethods: jsii.Strings("GET", "POST", "PUT", "DELETE", "OPTIONS"),
			AllowOrigins: jsii.Strings("https://www.idahvisng.com"), // Set to your frontend origin
		},
		DeployOptions: &awsapigateway.StageOptions{
			LoggingLevel: awsapigateway.MethodLoggingLevel_ERROR,
		},
	})

	// Integrate Lambda with API Gateway
	integration := awsapigateway.NewLambdaIntegration(myFunction, nil)
	messageResource := api.Root().AddResource(jsii.String("message"), nil)
	messageResource.AddMethod(jsii.String("POST"), integration, nil)

	// SES Configuration for Email Identity Verification
  awsses.NewEmailIdentity(stack, jsii.String("idahvisDomainIdentity"), &awsses.EmailIdentityProps{
		Identity: awsses.Identity_Domain(jsii.String("idahvisng.com")),
	})
   // Output the DNS records that need to be added to your domain
    awscdk.NewCfnOutput(stack, jsii.String("DkimRecords"), &awscdk.CfnOutputProps{
    Value: jsii.String("Check the AWS SES console for DKIM records"),
    Description: jsii.String("DKIM records need to be added as CNAME records to your domain"),
  })

    awscdk.NewCfnOutput(stack, jsii.String("VerificationRecord"), &awscdk.CfnOutputProps{
    Value: jsii.String("Check the AWS SES console for verification record"),
    Description: jsii.String("Add the verification record as a TXT record to your domain"),
})
	// Define SES Policy for Lambda to send emails
	sesPolicy := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("ses:SendEmail", "ses:SendRawEmail"),
		Resources: jsii.Strings("*"),
	})
  
	// Attach SES Policy to Lambda
	myFunction.AddToRolePolicy(sesPolicy)
	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewIdahvisBackendStack(app, "IdahvisBackendStack", &IdahvisBackendStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to be deployed.
func env() *awscdk.Environment {
	// Uncomment and set your AWS account and region for production stacks
	// return &awscdk.Environment{
	// 	Account: jsii.String("123456789012"),
	// 	Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region
	// that are implied by the current CLI configuration. Recommended for dev stacks.
	// return &awscdk.Environment{
	// 	Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	// 	Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }

	// Default to environment-agnostic stack
	return nil
}

