FROM golang:1.18 as builder
MAINTAINER OpsTree Solutions

# Set working directory
WORKDIR /go/src/employee

# Copy source code
COPY ./ /go/src/employee

# FIX: Use 'go mod download' and build the specific path '.' 
# to avoid scanning external cache folders
RUN go mod download && \
    go build -o employee .

FROM alpine:latest
MAINTAINER OpsTree Solutions
WORKDIR /app

# Install necessary libraries for Alpine compatibility
RUN apk add --no-cache libc6-compat curl

# Copy binary from builder
COPY --from=builder /go/src/employee/employee /app/

# Setup config
RUN mkdir -p /root/config
COPY ./config.yaml /root/config/config.yaml
ENV CONFIG_FILE "/root/config/config.yaml"

EXPOSE 8083
ENTRYPOINT ["./employee"]
