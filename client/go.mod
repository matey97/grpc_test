module "github.com/matey97/grpc_test/client"

go 1.15

replace github.com/matey97/grpc_test/grpc_test => ../grpc_test

require (
	github.com/matey97/grpc_test/grpc_test v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.4.3 // indirect
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0
)