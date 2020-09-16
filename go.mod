module github.com/beatitudes/shippy-service-consignment

go 1.14

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/beatitudes/shippy-service-vessel v0.0.0-20200916024440-9b8d78eabf98 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/micro/go-micro/v2 v2.9.1
	go.mongodb.org/mongo-driver v1.4.1 // indirect
	google.golang.org/protobuf v1.25.0
)
