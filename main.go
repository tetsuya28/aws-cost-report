package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/slack-go/slack"
	"github.com/tetsuya28/aws_cost_report/external"
)

type Config struct {
	SlackToken   string `required:"true" envconfig:"SLACK_TOKEN"`
	SlackChannel string `required:"true" envconfig:"SLACK_CHANNEL"`
}

func main() {
	lambda.Start(handler)
}

func handler() error {
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		panic(err)
	}
	slk := external.NewSlack(config.SlackToken)

	now := time.Now()
	yesterday := now.Add(-1 * 24 * time.Hour)

	result, err := external.GetCost(yesterday, now)
	if err != nil {
		return err
	}

	totalCost := 0.0
	attachments := make([]slack.Attachment, len(result.ResultsByTime), 0)
	for i := range result.ResultsByTime {
		for _, service := range result.ResultsByTime[i].Groups {
			value, ok := service.Metrics["BlendedCost"]
			if !ok || value == nil {
				continue
			}

			if value.Amount == nil {
				continue
			}

			usageQuantity, ok := service.Metrics["UsageQuantity"]
			if !ok || usageQuantity == nil {
				continue
			}

			attachment := slack.Attachment{
				Color: "#00ff00",
				Fields: []slack.AttachmentField{
					{
						Title: "料金",
						Value: fmt.Sprintf("%s%s", *value.Amount, *value.Unit),
						Short: true,
					},
					{
						Title: "使用量",
						Value: *usageQuantity.Amount,
						Short: true,
					},
				},
				AuthorName: *service.Keys[0],
				AuthorIcon: external.GetIconURL(*service.Keys[0]),
			}
			attachments = append(attachments, attachment)

			cost, err := strconv.ParseFloat(*value.Amount, 10)
			if err != nil {
				continue
			}
			totalCost += cost
		}
	}
	text := fmt.Sprintf("%sのコスト一覧\n合計金額: $%.3f", yesterday.Format("2006-01-02"), totalCost)
	option := slack.MsgOptionText(text, false)
	err = slk.PostMessage(config.SlackChannel, option, slack.MsgOptionAttachments(attachments...))
	if err != nil {
		return err
	}
	return nil
}
