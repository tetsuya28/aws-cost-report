TAG ?= v0.0.0
SLACK_TOKEN ?=
SLACK_CHANNEL ?=

.PHONY: run
run:
	SLACK_TOKEN=$(SLACK_TOKEN) SLACK_CHANNEL=$(SLACK_CHANNEL) go run main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: integration-test
integration-test:
	go clean -testcache
	SLACK_TOKEN=$(SLACK_TOKEN) SLACK_CHANNEL=$(SLACK_CHANNEL) go test -v ./... -tags=integration

.PHONY: upload
upload: build
	zip -j ./bin/main.zip ./bin/main
	aws s3 cp ./bin/main.zip s3://tetsuya28-aws-cost-report/$(TAG)/main.zip

.PHONY: terraform/docs
terraform/docs:
	cd module && terraform-docs markdown table . > README.md
