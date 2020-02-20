FROM local/base:latest as builder

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -a -tags netgo -ldflags='-w -s -extldflags "-static"' -o /go/bin/cli ./cli

FROM scratch
# # Copy our static executable.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/bin/cli /go/bin/cli

EXPOSE 9100:9109
ENTRYPOINT ["/go/bin/cli"]
