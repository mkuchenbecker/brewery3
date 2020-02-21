FROM local/base:latest as builder

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -a -tags netgo -ldflags='-w -s -extldflags "-static"' -o /go/bin/controller ./controller
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -a -tags netgo -ldflags='-w -s -extldflags "-static"' -o /go/bin/cli ./cli

# RUN GRPC_GO_LOG_VERBOSITY_LEVEL=99
FROM golang:alpine as slate
# # Copy our static executable.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/bin/controller /go/bin/controller
COPY --from=builder /go/bin/cli /go/bin/cli

ENTRYPOINT ["/go/bin/controller"]
EXPOSE 9000:9009
