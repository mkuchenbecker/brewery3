FROM golang

ARG app_env
ENV APP_ENV $app_env

COPY ./ /go/src/github.com/mkuchenbecker/brewery3
WORKDIR /go/src/github.com/mkuchenbecker/brewery3

RUN go get -u github.com/golang/dep/cmd/dep \
  && go get github.com/kyoh86/richgo \
  && go get -u github.com/golangci/golangci-lint/cmd/golangci-lint \
  && go get github.com/mattn/goveralls

RUN dep ensure --vendor-only

RUN go build ./entry/server

CMD server

EXPOSE 8080