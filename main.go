package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/kelseyhightower/envconfig"
	"github.com/slack-go/slack"
	"github.com/tetsuya28/aws_cost_report/external"
	"github.com/tetsuya28/aws_cost_report/testdata"
)

type Config struct {
	SlackToken   string `required:"true" envconfig:"SLACK_TOKEN"`
	SlackChannel string `required:"true" envconfig:"SLACK_CHANNEL"`
}

func main() {
	// Exec on Lambda or not
	if os.Getenv("_HANDLER") != "" {
		lambda.Start(handler)
	} else {
		handler()
	}
}

func handler() error {
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		panic(err)
	}
	slk := external.NewSlack(config.SlackToken)

	now := time.Now()
	yesterday := now.Add(-1 * 24 * time.Hour).Format("2006-01-02")

	var result *costexplorer.GetCostAndUsageOutput
	var err error
	if os.Getenv("_HANDLER") != "" {
		result, err = external.GetCost()
	} else {
		result, err = testdata.GetCostAndUsage()
	}
	if err != nil {
		panic(err)
	}

	totalCost := 0.0
	attachments := make([]slack.Attachment, 0)
	for i := range result.ResultsByTime {
		for _, service := range result.ResultsByTime[i].Groups {
			attachment := slack.Attachment{
				Color: "#00ff00",
				Fields: []slack.AttachmentField{
					{
						Title: "料金",
						Value: fmt.Sprintf("%s%s", *service.Metrics["BlendedCost"].Amount, *service.Metrics["BlendedCost"].Unit),
						Short: true,
					},
					{
						Title: "使用量",
						Value: *service.Metrics["UsageQuantity"].Amount,
						Short: true,
					},
				},
				AuthorName: *service.Keys[0],
				AuthorIcon: external.GetIconURL(*service.Keys[0]),
			}
			attachments = append(attachments, attachment)
			cost, err := strconv.ParseFloat(*service.Metrics["BlendedCost"].Amount, 10)
			if err != nil {
				continue
			}
			totalCost += cost
		}
	}
	text := fmt.Sprintf("%sのコスト一覧\n合計金額: $%.3f", yesterday, totalCost)
	option := slack.MsgOptionText(text, false)
	err = slk.PostMessage(config.SlackChannel, option, slack.MsgOptionAttachments(attachments...))
	if err != nil {
		panic(err)
	}
	return nil
}
