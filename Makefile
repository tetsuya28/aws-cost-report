TAG ?= v0.0.0

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o bin/main main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: upload
upload: build
	zip -j ./bin/main.zip ./bin/main
	aws s3 cp ./bin/main.zip s3://tetsuya28-aws-cost-report/$(TAG)/main.zip

.PHONY: terraform/docs
terraform/docs:
	cd module && terraform-docs markdown table . > README.md
