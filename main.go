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
	"github.com/tetsuya28/aws-cost-report/i18y"
	"github.com/ucpr/mongo-streamer/pkg/log"
)

var (
	Language string
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
		log.Warn("failed to new config, err=%w", err)
		return err
	}

	Language = cfg.Language

	err = i18y.Init()
	if err != nil {
		log.Warn("failed to init i18y, err=%w", err)
		return err
	}

	slk := external.NewSlack(cfg.SlackToken)

	result, err := external.GetCost()
	if err != nil {
		log.Warn("failed to get cost, err=%w", err)
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

			// Sum total daily cost
			dailyCost.Total += c.CostAmount
		}

		cost[i] = dailyCost
	}

	fullName, err := external.GetAccountFullName(context.Background())
	if err != nil {
		log.Warn("failed to get account info, err=%w", err)
		return err
	}

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	text := i18y.Translate(Language, "title", fullName, yesterday.Format("2006-01-02"), cost[1].Total)
	option := slack.MsgOptionText(text, false)

	attachments := toAttachment(cost)
	err = slk.PostMessage(cfg.SlackChannel, option, slack.MsgOptionAttachments(attachments...))
	if err != nil {
		log.Warn("failed to post message to Slack, err=%w", err)
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
	// Just day before yesterday and yesterday
	if len(cost) != 2 {
		log.Warn("cost length is not 2")
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
				diffMark := ""
				// Set red color if diff is over 100%
				if diff == 100 {
					color = "#ffffff"
				} else if diff > 100 {
					color = "#ff0000"
					diffMark = "📈"
				} else {
					color = "#0000ff"
					diffMark = "📉"
				}

				priceDiffStatement = fmt.Sprintf(
					" ( %s %.1f%% )",
					diffMark,
					diff,
				)
			}
		}

		fields := []slack.AttachmentField{
			{
				Title: i18y.Translate(Language, "cost"),
				Value: fmt.Sprintf(
					"%.3f%s%s",
					detail.CostAmount,
					detail.CostUnit,
					priceDiffStatement,
				),
				Short: true,
			},
			{
				Title: i18y.Translate(Language, "usage"),
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
