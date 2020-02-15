.PHONY: bootstrap
bootstrap:
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -u github.com/kyoh86/richgo
	go get -v github.com/ramya-rao-a/go-outline
	go get -v github.com/mdempsky/gocode
	go get -v github.com/uudashr/gopkgs/cmd/gopkgs
	go get -v golang.org/x/tools/cmd/goimports

.PHONY: lint
lint:
	golangci-lint run

.PHONY: tests
tests:
	@echo "tests:"
	${GOPATH}/bin/richgo test \
	-timeout 10s \
	-cover \
	-v \
	-covermode=count \
	-coverprofile=coverage.out \
	-tags test \
	./...

.PHONY: test-ci
test-ci:
	@echo "tests:"
	go test -v -covermode=count -coverprofile=coverage.out ./...

.PHONY: install-golang-ci
lint-ci:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh


.PHONY: breweryproto
breweryproto:
	@echo "compiling protos:"
	protoc -I brewery/model \
	brewery/model/config.proto \
	brewery/model/switch.proto \
	brewery/model/thermometer.proto \
	--proto_path=. \
	--go_out=plugins=grpc:brewery/model/gomodel

.PHONY: dataproto
dataproto:
	@echo "compiling protos:"
	protoc -I data/protos/ \
	data.proto \
	--proto_path=. \
	--go_out=plugins=grpc:data/gomodel

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


.PHONY: start
start:
	kubectl create -f brewery.yml

.PHONY: stop
stop:
	kubectl delete -n default deployment brewery-deployment
	kubectl delete -n default service brewery-service


.PHONY: client
client:
	docker build -t local/client -f ./cli/Dockerfile ./cli
	docker tag local/client mkuchenbecker/brewery3:client-latest
	docker push mkuchenbecker/brewery3:client-latest


.PHONY: build
build:
	docker build -t local/brewery -f Dockerfile.brewery . \
	&& docker build -t local/element -f Dockerfile.element . \
	&& docker build -t local/thermometer -f Dockerfile.thermometer . \

.PHONY: tag
tag:
	docker tag local/brewery mkuchenbecker/brewery3:brewery-latest
	docker tag local/element mkuchenbecker/brewery3:element-latest
	docker tag local/thermometer mkuchenbecker/brewery3:thermometer-latest



