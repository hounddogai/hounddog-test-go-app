module example.com/hounddog-test-go-app/client

go 1.21

require (
	google.golang.org/grpc v1.66.0
	gorm.io/driver/sqlite v1.5.7
	gorm.io/gorm v1.25.12
	example.com/hounddog-test-go-app/proto v0.0.0
	example.com/hounddog-test-go-app/utils v0.0.0
)

replace example.com/hounddog-test-go-app/proto => ../proto
replace example.com/hounddog-test-go-app/utils => ../utils
