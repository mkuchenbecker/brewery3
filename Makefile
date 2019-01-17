COVERAGE_FILE := coverage/cover.out
COVERAGE_ANALYSIS_FILE := coverage/cover.analysis
COVERAGE_ANALYSIS_FILE_XML := coverage/coverage.xml
COVERAGE_ANALYSIS_FILE_HTML := coverage/coverage.html

.PHONY: lint
lint: fmt
	golangci-lint run

.PHONY: tests
tests: fmt lint
	@echo "tests:"
	${GOPATH}/bin/richgo test -timeout 10s -cover -race -tags test ./...

.PHONY: coverage
coverage: generate
	@echo "coverage:"
	${GOPATH}/bin/richgo test -covermode=atomic -coverprofile ${COVERAGE_FILE} ./...
	go tool cover -func=${COVERAGE_FILE} -o ${COVERAGE_ANALYSIS_FILE}
	go tool cover -html=${COVERAGE_FILE} -o ${COVERAGE_ANALYSIS_FILE_HTML}
	gocover-cobertura < ${COVERAGE_FILE} > ${COVERAGE_ANALYSIS_FILE_XML}

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
	ThermometerClient \
	> brewery/model/gomock/gomock_models.go

.PHONY: structmockgen
structmockgen:
	@echo "generating mocks from structs:"
	go generate ./...

.PHONY: generate
generate: proto protomockgen structmockgen
