module github.com/pavlovicisidora/soa-team7/Backend/Shopping

go 1.25.0

require (
	github.com/gorilla/mux v1.8.1
	github.com/pavlovicisidora/soa-team7/Backend/APIGateway v0.0.0-00010101000000-000000000000
	go.mongodb.org/mongo-driver v1.15.0
	google.golang.org/grpc v1.75.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/pavlovicisidora/soa-team7/Backend/APIGateway => ../api_gateway
