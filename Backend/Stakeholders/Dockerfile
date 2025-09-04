FROM golang:1.25-alpine as builder
WORKDIR /app

RUN go mod init temp-build-module

COPY go.work go.work.sum ./
COPY Stakeholders/go.mod Stakeholders/go.sum ./Stakeholders/
COPY Blog/go.mod Blog/go.sum ./Blog/

COPY api_gateway/go.mod api_gateway/go.sum ./api_gateway/

RUN go get -d ./...

COPY Stakeholders/ ./Stakeholders/

WORKDIR /app/Stakeholders

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/stakeholders .

FROM alpine:latest
COPY --from=builder /app/stakeholders /usr/bin/stakeholders
EXPOSE 8081
ENTRYPOINT ["/usr/bin/stakeholders"]