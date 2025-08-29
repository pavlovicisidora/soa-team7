FROM golang:1.25-alpine as builder
WORKDIR /app

RUN go mod init temp-build-module

COPY go.work go.work.sum ./
COPY api_gateway/go.mod api_gateway/go.sum ./api_gateway/
COPY Blog/go.mod Blog/go.sum ./Blog/
COPY Stakeholders/go.mod Stakeholders/go.sum ./Stakeholders/

RUN go get -d ./...

COPY api_gateway/ ./api_gateway/
COPY Blog/ ./Blog/ 

WORKDIR /app/api_gateway

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api-gateway .

FROM alpine:latest
WORKDIR /
COPY --from=builder /app/api-gateway /usr/bin/api-gateway
EXPOSE 8080
ENTRYPOINT ["/usr/bin/api-gateway"]