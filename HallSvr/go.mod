module HallSvr

go 1.15

require (
	Common v0.0.0
	github.com/gin-gonic/gin v1.7.2
	github.com/golang/protobuf v1.5.0
	github.com/smartystreets/goconvey v1.7.2 // indirect
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55
)

replace Common v0.0.0 => ../Common
