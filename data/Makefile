.PHONY: bootstrap
bootstrap:
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -u github.com/kyoh86/richgo
	go get -v github.com/ramya-rao-a/go-outline
	go get -v github.com/mdempsky/gocode
	go get -v github.com/uudashr/gopkgs/cmd/gopkgs
	go get -v golang.org/x/tools/cmd/goimports

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
	./...

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
	../scripts/fmt


.PHONY: dataproto
dataproto:
	@echo "compiling protos:"
	protoc -I data/protos/ \
	data.proto \
	--proto_path=. \
	--go_out=plugins=grpc:data/gomodel


.PHONY: build
build:
	docker build -t datastore -f Dockerfile .

.PHONY: tag
tag:
	docker tag datastore mkuchenbecker/brewery3:datastore-latest

.PHONY: publish
publish: build tag
	docker push mkuchenbecker/brewery3:datastore-latest

.PHONY: apply
apply:
	kubectl apply -f k8s.yml

.PHONY: delete
delete:
	kubectl delete -n default deployment datastore-deployment
	kubectl delete -n default service datastore-service

