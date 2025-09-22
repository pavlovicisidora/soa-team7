FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN go mod init temp-build-module

COPY go.work go.work.sum ./
COPY Shopping/go.mod Shopping/go.sum ./Shopping/
COPY Blog/go.mod Blog/go.sum ./Blog/
COPY api_gateway/go.mod api_gateway/go.sum ./api_gateway/
COPY Stakeholders/go.mod Stakeholders/go.sum ./Stakeholders/
COPY Follower/go.mod Follower/go.sum ./Follower/
COPY saga/go.mod ./saga/
COPY Tour/pom.xml ./Tour/ 
COPY common/go.mod ./common/

RUN go get -d ./...

COPY common/ ./common/
COPY Shopping/ ./Shopping/
COPY saga/ ./saga/ 
COPY api_gateway/ ./api_gateway/

WORKDIR /app/Shopping

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/shopping .

FROM alpine:latest
COPY --from=builder /app/shopping /usr/bin/shopping
EXPOSE 8085
ENTRYPOINT ["/usr/bin/shopping"]
