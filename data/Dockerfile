FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
# RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/mkuchenbecker/brewery3/data
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -a -tags netgo -ldflags='-w -s -extldflags "-static"' -o /go/bin/datastore
############################
# STEP 2 build a small image
############################
FROM scratch
# # Copy our static executable.
COPY --from=builder /go/bin/datastore /go/bin/datastore
# Run the hello binary.
ENTRYPOINT ["/go/bin/datastore"]
