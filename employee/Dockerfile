FROM golang:1.18 as builder
MAINTAINER OpsTree Solutions

# Set the working directory first
WORKDIR /go/src/employee

# Copy the source code
COPY ./ /go/src/employee

# FIX: Remove the recursive ./... and build the specific package
# This prevents Go from scanning your Jenkins 'go-cache' folder
RUN go mod download && \
    go build -o employee .

FROM alpine:latest
MAINTAINER OpsTree Solutions
WORKDIR /app

# Install necessary libraries for Go binaries to run on Alpine
RUN apk add --no-cache libc6-compat curl

# Copy the binary from the builder stage
COPY --from=builder /go/src/employee/employee /app/

# Copy config to the location the ENV expects
COPY ./config.yaml /root/config/config.yaml
ENV CONFIG_FILE "/root/config/config.yaml"

EXPOSE 8083
ENTRYPOINT ["./employee"]
