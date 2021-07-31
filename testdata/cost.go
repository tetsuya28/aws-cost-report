package testdata

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

type sampleMetric struct {
	Amount string
	Unit   string
}

type sampleData struct {
	ServiceName   string
	BlendedCost   sampleMetric
	UsageQuantity sampleMetric
}

func GetCostAndUsage() (*costexplorer.GetCostAndUsageOutput, error) {
	// 	costExplorerResponse := `
	// [
	//     {
	//       Keys: ["AWS Glue"],
	//       Metrics: {
	//         AmortizedCost: {
	//           Amount: "0",
	//           Unit: "USD"
	//         },
	//         BlendedCost: {
	//           Amount: "0",
	//           Unit: "USD"
	//         },
	//         UnblendedCost: {
	//           Amount: "0",
	//           Unit: "USD"
	//         },
	//         UsageQuantity: {
	//           Amount: "0.008064516",
	//           Unit: "N/A"
	//         }
	//       }
	//     },
	//     {
	//       Keys: ["Amazon DynamoDB"],
	//       Metrics: {
	//         UsageQuantity: {
	//           Amount: "120.0000203899",
	//           Unit: "N/A"
	//         },
	//         AmortizedCost: {
	//           Amount: "0",
	//           Unit: "USD"
	//         },
	//         BlendedCost: {
	//           Amount: "0",
	//           Unit: "USD"
	//         },
	//         UnblendedCost: {
	//           Amount: "0",
	//           Unit: "USD"
	//         }
	//       }
	//     },
	//     {
	//       Keys: ["Amazon EC2 Container Registry (ECR)"],
	//       Metrics: {
	//         BlendedCost: {
	//           Amount: "0.0001027965",
	//           Unit: "USD"
	//         },
	//         UnblendedCost: {
	//           Amount: "0.0001027965",
	//           Unit: "USD"
	//         },
	//         UsageQuantity: {
	//           Amount: "0.001027965",
	//           Unit: "N/A"
	//         },
	//         AmortizedCost: {
	//           Amount: "0.0001027965",
	//           Unit: "USD"
	//         }
	//       }
	//     },
	//     {
	//       Keys: ["Amazon Elastic Container Service"],
	//       Metrics: {
	//         UnblendedCost: {
	//           Amount: "0.03081",
	//           Unit: "USD"
	//         },
	//         UsageQuantity: {
	//           Amount: "1.5",
	//           Unit: "N/A"
	//         },
	//         AmortizedCost: {
	//           Amount: "0.03081",
	//           Unit: "USD"
	//         },
	//         BlendedCost: {
	//           Amount: "0.03081",
	//           Unit: "USD"
	//         }
	//       }
	//     },
	//     {
	//       Keys: ["Amazon Elastic Load Balancing"],
	//       Metrics: {
	//         AmortizedCost: {
	//           Amount: "0.0729046222",
	//           Unit: "USD"
	//         },
	//         BlendedCost: {
	//           Amount: "0.0729046222",
	//           Unit: "USD"
	//         },
	//         UnblendedCost: {
	//           Amount: "0.0729046222",
	//           Unit: "USD"
	//         },
	//         UsageQuantity: {
	//           Amount: "3.0009415198",
	//           Unit: "N/A"
	//         }
	//       }
	//     },
	//     {
	//       Keys: ["Amazon Relational Database Service"],
	//       Metrics: {
	//         BlendedCost: {
	//           Amount: "0.0557096774",
	//           Unit: "USD"
	//         },
	//         UnblendedCost: {
	//           Amount: "0.0557096774",
	//           Unit: "USD"
	//         },
	//         UsageQuantity: {
	//           Amount: "2.0268915207",
	//           Unit: "N/A"
	//         },
	//         AmortizedCost: {
	//           Amount: "0.0557096774",
	//           Unit: "USD"
	//         }
	//       }
	//     },
	//     {
	//       Keys: ["Amazon Route 53"],
	//       Metrics: {
	//         UnblendedCost: {
	//           Amount: "0.0000252",
	//           Unit: "USD"
	//         },
	//         UsageQuantity: {
	//           Amount: "67",
	//           Unit: "N/A"
	//         },
	//         AmortizedCost: {
	//           Amount: "0.0000252",
	//           Unit: "USD"
	//         },
	//         BlendedCost: {
	//           Amount: "0.0000252",
	//           Unit: "USD"
	//         }
	//       }
	//     },
	//     {
	//       Keys: ["Amazon Simple Storage Service"],
	//       Metrics: {
	//         AmortizedCost: {
	//           Amount: "0.0003722487",
	//           Unit: "USD"
	//         },
	//         BlendedCost: {
	//           Amount: "0.0003722487",
	//           Unit: "USD"
	//         },
	//         UnblendedCost: {
	//           Amount: "0.0003722487",
	//           Unit: "USD"
	//         },
	//         UsageQuantity: {
	//           Amount: "542.0001115065",
	//           Unit: "N/A"
	//         }
	//       }
	//     },
	//     {
	//       Keys: ["AmazonCloudWatch"],
	//       Metrics: {
	//         BlendedCost: {
	//           Amount: "0.0282258064",
	//           Unit: "USD"
	//         },
	//         UnblendedCost: {
	//           Amount: "0.0282258064",
	//           Unit: "USD"
	//         },
	//         UsageQuantity: {
	//           Amount: "4896.0989286767",
	//           Unit: "N/A"
	//         },
	//         AmortizedCost: {
	//           Amount: "0.0282258064",
	//           Unit: "USD"
	//         }
	//       }
	//     }
	// ]
	// `
	data := []sampleData{
		sampleData{
			ServiceName: "AWS Glue",
			BlendedCost: sampleMetric{
				Amount: "0",
				Unit:   "USD",
			},
			UsageQuantity: sampleMetric{
				Amount: "0",
				Unit:   "N/A",
			},
		},
		sampleData{
			ServiceName: "Amazon DynamoDB",
			BlendedCost: sampleMetric{
				Amount: "0",
				Unit:   "USD",
			},
			UsageQuantity: sampleMetric{
				Amount: "120.0000203899",
				Unit:   "N/A",
			},
		},
		sampleData{
			ServiceName: "Amazon EC2 Container Registry (ECR)",
			BlendedCost: sampleMetric{
				Amount: "0.0001027965",
				Unit:   "USD",
			},
			UsageQuantity: sampleMetric{
				Amount: "0.001027965",
				Unit:   "N/A",
			},
		},
		sampleData{
			ServiceName: "Amazon Elastic Container Service",
			BlendedCost: sampleMetric{
				Amount: "0.03081",
				Unit:   "USD",
			},
			UsageQuantity: sampleMetric{
				Amount: "1.5",
				Unit:   "N/A",
			},
		},
		sampleData{
			ServiceName: "Amazon Elastic Load Balancing",
			BlendedCost: sampleMetric{
				Amount: "0.0729046222",
				Unit:   "USD",
			},
			UsageQuantity: sampleMetric{
				Amount: "3.0009415198",
				Unit:   "N/A",
			},
		},
		sampleData{
			ServiceName: "Amazon Relational Database Service",
			BlendedCost: sampleMetric{
				Amount: "0.0557096774",
				Unit:   "USD",
			},
			UsageQuantity: sampleMetric{
				Amount: "2.0268915207",
				Unit:   "N/A",
			},
		},
		sampleData{
			ServiceName: "Amazon Route 53",
			BlendedCost: sampleMetric{
				Amount: "0.0000252",
				Unit:   "USD",
			},
			UsageQuantity: sampleMetric{
				Amount: "67",
				Unit:   "N/A",
			},
		},
		sampleData{
			ServiceName: "Amazon Simple Storage Service",
			BlendedCost: sampleMetric{
				Amount: "0.0003722487",
				Unit:   "USD",
			},
			UsageQuantity: sampleMetric{
				Amount: "542.0001115065",
				Unit:   "N/A",
			},
		},
		sampleData{
			ServiceName: "AmazonCloudWatch",
			BlendedCost: sampleMetric{
				Amount: "0.0282258064",
				Unit:   "USD",
			},
			UsageQuantity: sampleMetric{
				Amount: "4896.0989286767",
				Unit:   "N/A",
			},
		},
	}

	groups := make([]*costexplorer.Group, 0)
	for _, g := range data {
		group := costexplorer.Group{
			Keys: []*string{aws.String(g.ServiceName)},
			Metrics: map[string]*costexplorer.MetricValue{
				"BlendedCost": &costexplorer.MetricValue{
					Amount: aws.String(g.BlendedCost.Amount),
					Unit:   aws.String(g.BlendedCost.Unit),
				},
				"UsageQuantity": &costexplorer.MetricValue{
					Amount: aws.String(g.UsageQuantity.Amount),
					Unit:   aws.String(g.UsageQuantity.Unit),
				},
			},
		}
		groups = append(groups, &group)
	}

	start := "2021-07-25"
	end := "2021-07-26"
	results := make([]*costexplorer.ResultByTime, 0)
	result := costexplorer.ResultByTime{
		Groups: groups,
		TimePeriod: &costexplorer.DateInterval{
			Start: &start,
			End:   &end,
		},
		Total: nil,
	}
	results = append(results, &result)
	output := costexplorer.GetCostAndUsageOutput{
		ResultsByTime: results,
	}

	return &output, nil
}
