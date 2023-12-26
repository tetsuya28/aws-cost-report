package external

import (
	"time"

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
	case "Amazon EC2 Container Registry (ECR)":
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
	default:
		return ""
	}
}

func GetCost(startDay, endDay time.Time) (*costexplorer.GetCostAndUsageOutput, error) {
	start := startDay.Format("2006-01-02")
	end := endDay.Format("2006-01-02")
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
			Start: aws.String(start),
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
