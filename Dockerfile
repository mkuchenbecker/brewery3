# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang v1.11 base image
FROM golang:1.11

# Add Maintainer Info
LABEL maintainer="Mike Kuchenbecker <mkuchenbecker@github.com>"

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/mkuchenbecker/brewery3

# Copy everything from the current directory to the PWD
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 8080

RUN make build

# Run the executable
CMD ["./server.bin"]