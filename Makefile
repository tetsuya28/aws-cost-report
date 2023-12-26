DOCKER_TAG ?= $(shell git rev-parse --short HEAD)
REGION := ap-northeast-1
REPOSITORY = public.ecr.aws/a8y2r9d0/aws-cost-usage

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o bin/main main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: /docker/build
docker/build:
	docker build --platform=linux/amd64 -t $(REPOSITORY):$(DOCKER_TAG) .

.PHONY: /docker/push
docker/push: docker/build
	docker push $(REPOSITORY):$(DOCKER_TAG)

.PHONY: terraform/docs
terraform/docs:
	cd module && terraform-docs markdown table . > README.md

.PHONY: terraform/init
terraform/init:
	terraform init -backend-config="bucket=$(BUCKET)" -backend-config="profile=$(PROFILE)" -reconfigure

.PHONY: terraform/plan
terraform/plan: build
	terraform fmt --recursive
	terraform plan -var="profile=$(PROFILE)" -var="region=$(REGION)"

.PHONY: terraform/apply
terraform/apply: build
	terraform fmt --recursive
	terraform apply -var="profile=$(PROFILE)" -var="region=$(REGION)"
