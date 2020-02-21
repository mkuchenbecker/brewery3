FROM local/base:latest as builder

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -a -tags netgo -ldflags='-w -s -extldflags "-static"' -o /go/bin/thermometer ./thermometer

FROM golang:alpine as slate
# # Copy our static executable.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/bin/thermometer /go/bin/thermometer

ENTRYPOINT ["/go/bin/thermometer"]
EXPOSE 9100:9109
