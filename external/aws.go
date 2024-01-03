package external

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/account"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func GetIconURL(service string) string {
	switch service {
	case "AWS Cost Explorer":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/CloudFinancialManagement/CostExplorer.png?raw=true"
	case "AWS Key Management Service":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/SecurityIdentityCompliance/KeyManagementService.png?raw=true"
	case "AWS Lambda":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/Compute/Lambda.png?raw=true"
	case "AWS X-Ray":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/DeveloperTools/XRay.png?raw=true"
	case "Amazon API Gateway":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/ApplicationIntegration/APIGateway.png?raw=true"
	case "Amazon Simple Email Service":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/BusinessApplications/SimpleEmailService.png?raw=true"
	case "Amazon DynamoDB":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/Database/DynamoDB.png?raw=true"
	case "Amazon EC2 Container Registry (ECR)", "Amazon Elastic Container Registry Public":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/Containers/ElasticContainerRegistry.png?raw=true"
	case "Amazon Elastic Container Service":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/Containers/ElasticContainerService.png?raw=true"
	case "Amazon Elastic Load Balancing":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/NetworkingContentDelivery/ElasticLoadBalancingApplicationLoadBalancer.png?raw=true"
	case "Amazon Relational Database Service":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/Database/Aurora.png?raw=true"
	case "Amazon Route 53":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/NetworkingContentDelivery/Route53.png?raw=true"
	case "Amazon Simple Storage Service":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/Storage/SimpleStorageService.png?raw=true"
	case "AmazonCloudWatch":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/ManagementGovernance/CloudWatch.png?raw=true"
	case "Amazon CloudFront":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/NetworkingContentDelivery/CloudFront.png?raw=true"
	case "AWS Amplify":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/FrontEndWebMobile/AmplifyAWSAmplifyStudio.png?raw=true"
	case "AWS Glue":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/Analytics/Glue.png?raw=true"
	case "Amazon Simple Notification Service":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/ApplicationIntegration/SimpleNotificationService.png?raw=true"
	case "AWS Secrets Manager":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/SecurityIdentityCompliance/SecretsManager.png?raw=true"
	case "Amazon Virtual Private Cloud":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/Groups/VPC.png?raw=true"
	case "AWS WAF":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/SecurityIdentityCompliance/WAF.png?raw=true"
	case "EC2 - Other":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/Compute/EC2.png?raw=true"
	case "AWS Step Functions":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/ApplicationIntegration/StepFunctions.png?raw=true"
	case "Amazon Elastic Compute Cloud - Compute":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/Compute/EC2.png?raw=true"
	case "AWS CloudTrail":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/ManagementGovernance/CloudTrail.png?raw=true"
	case "Amazon Simple Queue Service":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/ApplicationIntegration/SimpleQueueService.png?raw=true"
	case "Amazon GuardDuty":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/SecurityIdentityCompliance/GuardDuty.png?raw=true"
	case "Amazon Elastic Container Service for Kubernetes":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/Containers/EKSCloud.png?raw=true"
	case "Amazon Cognito":
		return "https://github.com/awslabs/aws-icons-for-plantuml/blob/main/dist/SecurityIdentityCompliance/Cognito.png?raw=true"
	case "Tax":
		return ""
	default:
		return ""
	}
}

func GetCost() (*costexplorer.GetCostAndUsageOutput, error) {
	now := time.Now()
	end := now.Format("2006-01-02")
	twoDaysBefore := now.AddDate(0, 0, -3).Format("2006-01-02")

	granularity := "DAILY"
	metrics := []string{
		"AmortizedCost",
		"BlendedCost",
		"UnblendedCost",
		"UsageQuantity",
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return nil, err
	}

	svc := costexplorer.New(sess)
	result, err := svc.GetCostAndUsage(&costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(twoDaysBefore),
			End:   aws.String(end),
		},
		Granularity: aws.String(granularity),
		GroupBy: []*costexplorer.GroupDefinition{
			{
				Type: aws.String("DIMENSION"),
				Key:  aws.String("SERVICE"),
			},
		},
		Metrics: aws.StringSlice(metrics),
		Filter: &costexplorer.Expression{
			Not: &costexplorer.Expression{
				Dimensions: &costexplorer.DimensionValues{
					Key:    aws.String("RECORD_TYPE"),
					Values: aws.StringSlice([]string{"Credit"}),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetAccountFullName(ctx context.Context) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", err
	}

	client := account.NewFromConfig(cfg)

	output, err := client.GetContactInformation(ctx, &account.GetContactInformationInput{})
	if err != nil {
		return "", err
	}

	return *output.ContactInformation.FullName, nil
}
