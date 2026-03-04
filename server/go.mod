module example.com/hounddog-test-go-app/server

go 1.21

require (
	example.com/hounddog-test-go-app/proto v0.0.0
	example.com/hounddog-test-go-app/utils v0.0.0
	google.golang.org/grpc v1.66.0
)

require (
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240604185151-ef581f913117 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace example.com/hounddog-test-go-app/proto => ../proto

replace example.com/hounddog-test-go-app/utils => ../utils
