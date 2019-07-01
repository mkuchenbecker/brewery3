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
	${GOPATH}/bin/go test \
	-timeout 10s \
	-cover \
	-v \
	-covermode=count \
	-coverprofile=coverage.out \
	-tags test \
	./brewery/...

.PHONY: build
build:
	go build -o cli.bin entry/cli/main.go
	go build -o server.bin entry/server/main.go

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
