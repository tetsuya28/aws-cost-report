include .env

REGION := ap-northeast-1

.PHONY: run
run:
	export `cat .env` && go run main.go

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o bin/main main.go

.PHONY: init
init:
	terraform init -backend-config="bucket=$(BUCKET)" -backend-config="profile=$(PROFILE)" -reconfigure

.PHONY: plan
plan: build
	terraform fmt --recursive
	terraform plan -var="profile=$(PROFILE)" -var="region=$(REGION)"

.PHONY: apply
apply: build
	terraform fmt --recursive
	terraform apply -var="profile=$(PROFILE)" -var="region=$(REGION)"

