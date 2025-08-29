FROM golang:1.25-alpine as builder
WORKDIR /app

RUN go mod init temp-build-module

COPY go.work go.work.sum ./
COPY Blog/go.mod Blog/go.sum ./Blog/
COPY api_gateway/go.mod api_gateway/go.sum ./api_gateway/
COPY Stakeholders/go.mod Stakeholders/go.sum ./Stakeholders/

RUN go get -d ./...

COPY Blog/ ./Blog/

WORKDIR /app/Blog

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/blog .

FROM alpine:latest
COPY --from=builder /app/blog /usr/bin/blog
EXPOSE 8082
ENTRYPOINT ["/usr/bin/blog"]