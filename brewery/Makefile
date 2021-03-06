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


.PHONY: proto
proto:
	@echo "compiling protos:"
	protoc -I model \
	model/config.proto \
	model/switch.proto \
	model/thermometer.proto \
	--proto_path=. \
	--go_out=plugins=grpc:model/gomodel


.PHONY: protomockgen
protomockgen:
	@echo "generating mocks from protos:"
	mockgen github.com/mkuchenbecker/brewery3/brewery/model/gomodel \
	SwitchClient,\
	ThermometerClient,\
	BreweryClient \
	> model/gomock/gomock_models.go

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
	docker build -t local/client -f ./client.dockerfile .
	docker tag local/client mkuchenbecker/brewery3:client-latest
	docker push mkuchenbecker/brewery3:client-latest


.PHONY: base
base:
	docker build --no-cache -t local/base -f base.dockerfile . \
	&& docker tag local/base mkuchenbecker/brewery3:base-latest \
	&& docker push mkuchenbecker/brewery3:base-latest

.PHONY: build
build: controller element thermometer
	echo "success"

.PHONY: controller
controller:
	docker build --no-cache -t mkuchenbecker/brewery3:controller-latest -f controller.dockerfile . \

.PHONY: thermometer
thermometer:
	docker build --no-cache -t mkuchenbecker/brewery3:thermometer-latest -f thermometer.dockerfile . \

.PHONY: element
element:
	docker build --no-cache -t mkuchenbecker/brewery3:element-latest -f element.dockerfile .

.PHONY: cli
cli:
	 docker build  -t local/cli -f client.dockerfile .

.PHONY: tag
tag:
	docker tag local/controller mkuchenbecker/brewery3:controller-latest
	docker tag local/element mkuchenbecker/brewery3:element-latest
	docker tag local/thermometer mkuchenbecker/brewery3:thermometer-latest

.PHONY: publish
publish:
	docker push mkuchenbecker/brewery3:controller-latest
	docker push mkuchenbecker/brewery3:element-latest
	docker push mkuchenbecker/brewery3:thermometer-latest

.PHONY: probe
probe:
	ssh -Y pi@raspi-1 'sudo modprobe w1-gpio && sudo modprobe w1-therm && ls /sys/bus/w1/devices'

.PHONY: gpio-owner
gpio-owner:
	ssh -Y pi@raspi-1 'sudo chown root.gpio /dev/gpiomem && sudo chmod g+rw /dev/gpiomem'

