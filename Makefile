.PHONY: bootstrap
bootstrap:
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -u github.com/kyoh86/richgo
	go get -v github.com/ramya-rao-a/go-outline

.PHONY: lint
lint: fmt
	golangci-lint run

.PHONY: tests
tests: fmt lint
	@echo "tests:"
	${GOPATH}/bin/richgo test \
	-timeout 10s \
	-cover \
	-v \
	-covermode=count \
	-coverprofile=coverage.out \
	-tags test \
	./brewery/...

.PHONY: test-ci
test-ci:
	@echo "tests:"
	go test -v -covermode=count -coverprofile=coverage.out ./...

.PHONY: install-golang-ci
lint-ci:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh

.PHONY: fmt
fmt:
	@echo "fmt:"
	scripts/fmt

.PHONY: proto
proto:
	@echo "compiling protos:"
	protoc -I brewery/model \
	brewery/model/config.proto \
	brewery/model/switch.proto \
	brewery/model/thermometer.proto \
	--proto_path=. \
	--go_out=plugins=grpc:brewery/model/gomodel

.PHONY: protomockgen
protomockgen:
	@echo "generating mocks from protos:"
	mockgen github.com/mkuchenbecker/brewery3/brewery/model/gomodel \
	SwitchClient,\
	ThermometerClient,\
	BreweryClient \
	> brewery/model/gomock/gomock_models.go

.PHONY: structmockgen
structmockgen:
	@echo "generating mocks from structs:"
	go generate ./...

.PHONY: generate
generate: proto protomockgen structmockgen


.PHONY: lite-launch
lite-launch:
	kubectl create -f brewery.yml

.PHONY: launch
launch: build
	eval $(minikube docker-env)
	kubectl create -f brewery.yml

.PHONY: stop
stop:
	./scripts/stop.local

.PHONY: up
up: build
	kubectl set image deployment/brewery-deployment brewery=local/brewery:latest
	kubectl set image deployment/brewery-deployment element=local/element:latest
	kubectl set image deployment/brewery-deployment element=local/thermometer:latest


.PHONY: build
build:
	docker build -t local/brewery -f Dockerfile.brewery . \
	&& docker build -t local/element -f Dockerfile.element . \
	&& docker build -t local/thermometer -f Dockerfile.thermometer .
