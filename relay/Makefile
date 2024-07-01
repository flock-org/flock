SW_VERSION ?= latest
IMAGE_ORG ?= mcnet

IMAGE_TAG_BASE ?= quay.io/$(IMAGE_ORG)/flock_client_func
IMG ?= $(IMAGE_TAG_BASE):$(SW_VERSION)
build:
	@echo "Start flock relay build"
	go build -o ./bin/relay ./cmd/relay/relay.go
	go build -o ./bin/fr-adm ./cmd/admin/admin.go
	go build -o ./bin/client_func ./client/

lint:  ; $(info running linters...)
	@golangci-lint run --config=./.golangci.yaml ./...

docker-build:
	docker build --progress=plain --rm --tag client_func .

build-image:
	docker build --build-arg SW_VERSION="$(SW_VERSION)" -t ${IMG} .

push-image:
	docker push ${IMG}