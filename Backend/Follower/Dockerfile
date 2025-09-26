FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN go mod init temp-build-module

COPY go.work go.work.sum ./
COPY Follower/go.mod Follower/go.sum ./Follower/
COPY api_gateway/go.mod api_gateway/go.sum ./api_gateway/
COPY Stakeholders/go.mod Stakeholders/go.sum ./Stakeholders/
COPY Blog/go.mod Blog/go.sum ./Blog/
COPY Shopping/go.mod Shopping/go.sum ./Shopping/
COPY saga/go.mod ./saga/
COPY common/go.mod ./common/
RUN go get -d ./...
COPY common/ ./common/
COPY saga/ ./saga/ 
COPY Follower/ ./Follower/

WORKDIR /app/Follower

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/follower .

FROM alpine:latest
COPY --from=builder /app/follower /usr/bin/follower
EXPOSE 8084
ENTRYPOINT ["/usr/bin/follower"]