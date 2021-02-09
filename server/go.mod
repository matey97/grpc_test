module github.com/matey97/grpc_test/server

go 1.15

replace github.com/matey97/grpc_test/grpc_test => ../grpc_test

require (
	github.com/matey97/grpc_test/grpc_test v0.0.0-00010101000000-000000000000
	cloud.google.com/go/bigquery v1.8.0
	github.com/golang/protobuf v1.4.3 // indirect
	google.golang.org/api v0.24.0
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0
)
