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

.PHONY: test-ci
test-ci:
	@echo "tests:"
	go test -v -covermode=count -coverprofile=data/coverage.out ./data/...
	go test -v -covermode=count -coverprofile=brewery/coverage.out ./brewery/...

.PHONY: install-golang-ci
lint-ci:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh

.PHONY: fmt
fmt:
	@echo "fmt:"
	scripts/fmt

