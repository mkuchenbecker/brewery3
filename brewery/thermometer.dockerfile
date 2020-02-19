FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

COPY . /go/src/github.com/mkuchenbecker/brewery3/brewery
WORKDIR /go/src/github.com/mkuchenbecker/brewery3/brewery

RUN go get -u github.com/golang/dep/cmd/dep

# Fetch dependencies.
# Using go get.
RUN go get ./...
RUN apk add ca-certificates

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -a -tags netgo -ldflags='-w -s -extldflags "-static"' -o /go/bin/thermometer ./thermometer

FROM scratch
# # Copy our static executable.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/bin/thermometer /go/bin/thermometer
COPY --from=builder /go/src/github.com/mkuchenbecker/brewery3/brewery/scripts/start_temp.sh /go/bin/start_temp.sh

ENTRYPOINT ["./go/bin/start_temp.sh && /go/bin/thermometer"]
EXPOSE 9100:9109
