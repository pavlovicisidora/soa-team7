FROM golang:alpine as build_container
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o blog

FROM alpine
COPY --from=build_container /app/blog /usr/bin
EXPOSE 8082
ENTRYPOINT ["blog"]
