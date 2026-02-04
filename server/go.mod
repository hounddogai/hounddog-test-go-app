module example.com/hounddog-test-go-app/server

go 1.21

require (
	google.golang.org/grpc v1.66.0
	example.com/hounddog-test-go-app/proto v0.0.0
	example.com/hounddog-test-go-app/utils v0.0.0
)

replace example.com/hounddog-test-go-app/proto => ../proto
replace example.com/hounddog-test-go-app/utils => ../utils
