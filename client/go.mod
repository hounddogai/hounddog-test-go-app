module example.com/hounddog-test-go-app/client

go 1.24.0

require (
	example.com/hounddog-test-go-app/proto v0.0.0
	example.com/hounddog-test-go-app/utils v0.0.0
	google.golang.org/grpc v1.79.3
	gorm.io/driver/sqlite v1.5.7
	gorm.io/gorm v1.25.12
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace example.com/hounddog-test-go-app/proto => ../proto

replace example.com/hounddog-test-go-app/utils => ../utils
