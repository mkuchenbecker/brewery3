FROM golang

ARG app_env
ENV APP_ENV $app_env


COPY . /go/src/github.com/mkuchenbecker/brewery3/data
WORKDIR /go/src/github.com/mkuchenbecker/brewery3/data

RUN go get -u github.com/golang/dep/cmd/dep

# RUN dep ensure --vendor-only
RUN go get -d -v ./...

RUN go install -v .

CMD ["data"]

EXPOSE 9000:9009