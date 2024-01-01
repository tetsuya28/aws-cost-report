package main

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/slack-go/slack"
	"github.com/tetsuya28/aws-cost-report/config"
	"github.com/tetsuya28/aws-cost-report/external"
	"github.com/ucpr/mongo-streamer/pkg/log"
)

type DailyCost struct {
	Total    float64
	Services map[string]ServiceDetail
}

type ServiceDetail struct {
	CostAmount  float64
	CostUnit    string
	UsageAmount float64
	UsageUnit   string
}

func main() {
	lambda.Start(handler)
}

func handler() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	slk := external.NewSlack(cfg.SlackToken)

	result, err := external.GetCost()
	if err != nil {
		return err
	}

	cost := make([]DailyCost, len(result.ResultsByTime))
	for i := range result.ResultsByTime {
		dailyCost := DailyCost{
			Services: make(map[string]ServiceDetail),
		}

		for _, service := range result.ResultsByTime[i].Groups {
			serviceName := ""
			// Set service name
			if service.Keys != nil {
				serviceName = *service.Keys[0]
			}

			if serviceName == "" {
				log.Warn("service name is empty")
				continue
			}

			c, err := toCost(service)
			if err != nil {
				log.Warn("failed to convert cost: %v", err)
				continue
			}

			dailyCost.Services[serviceName] = c

			// 日次合計を計算する
			dailyCost.Total += c.CostAmount
		}

		cost[i] = dailyCost
	}

	fullName, err := external.GetAccountFullName(context.Background())
	if err != nil {
		return err
	}

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	text := fmt.Sprintf("%s の %s コスト\n合計金額: $%.3f", fullName, yesterday.Format("2006-01-02"), cost[0].Total)
	option := slack.MsgOptionText(text, false)

	attachments := toAttachment(cost)
	err = slk.PostMessage(cfg.SlackChannel, option, slack.MsgOptionAttachments(attachments...))
	if err != nil {
		return err
	}

	return nil
}

func toCost(result *costexplorer.Group) (ServiceDetail, error) {
	if result == nil {
		return ServiceDetail{}, nil
	}

	if result.Metrics == nil {
		return ServiceDetail{}, nil
	}

	if result.Metrics["BlendedCost"] == nil {
		return ServiceDetail{}, nil
	}

	costAmount, err := strconv.ParseFloat(*result.Metrics["BlendedCost"].Amount, 10)
	if err != nil {
		return ServiceDetail{}, nil
	}

	costUnit := ""
	if result.Metrics["BlendedCost"].Unit != nil {
		costUnit = *result.Metrics["BlendedCost"].Unit
	}

	if result.Metrics["UsageQuantity"] == nil {
		return ServiceDetail{}, nil
	}

	usageUnit := ""
	if result.Metrics["UsageQuantity"].Unit != nil && *result.Metrics["UsageQuantity"].Unit != "N/A" {
		usageUnit = *result.Metrics["UsageQuantity"].Unit
	}

	usageAmount, err := strconv.ParseFloat(*result.Metrics["UsageQuantity"].Amount, 10)
	if err != nil {
		return ServiceDetail{}, nil
	}

	return ServiceDetail{
		CostAmount:  costAmount,
		CostUnit:    costUnit,
		UsageAmount: usageAmount,
		UsageUnit:   usageUnit,
	}, nil
}

func toAttachment(cost []DailyCost) []slack.Attachment {
	// 一昨日、昨日のコスト比較なので 2 つのみ
	// [0] : 一昨日、 [1] : 昨日
	if len(cost) != 2 {
		return nil
	}

	attachments := make([]slack.Attachment, len(cost[1].Services))
	for name, detail := range cost[1].Services {
		color := "#00ff00"

		priceDiffStatement := ""
		before, ok := cost[0].Services[name]
		if ok {
			diff := (detail.CostAmount / before.CostAmount) * 100

			if !math.IsNaN(diff) {
				priceDiffStatement += " ( 前日比 : "

				// 前日よりも高くなってたら赤色にする
				if diff == 100 {
					color = "#ffffff"
					priceDiffStatement += ""
				} else if diff > 100 {
					color = "#ff0000"
					priceDiffStatement += "📈 "
				} else {
					color = "#0000ff"
					priceDiffStatement += "📉 "
				}

				priceDiffStatement += fmt.Sprintf(
					"%.1f%% )",
					diff,
				)
			}
		}

		fields := []slack.AttachmentField{
			{
				Title: "料金",
				Value: fmt.Sprintf(
					"%.3f%s%s",
					detail.CostAmount,
					detail.CostUnit,
					priceDiffStatement,
				),
				Short: true,
			},
			{
				Title: "使用量",
				Value: fmt.Sprintf("%.3f%s", detail.UsageAmount, detail.UsageUnit),
				Short: true,
			},
		}

		attachment := slack.Attachment{
			Color:      color,
			Fields:     fields,
			AuthorName: name,
			AuthorIcon: external.GetIconURL(name),
		}
		attachments = append(attachments, attachment)
	}

	return attachments
}
