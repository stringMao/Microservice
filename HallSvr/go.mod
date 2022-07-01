module HallSvr

go 1.15

require (
	Common v0.0.0
	github.com/gin-gonic/gin v1.7.2
	github.com/golang/protobuf v1.5.2
	github.com/smartystreets/goconvey v1.7.2 // indirect
)

replace Common v0.0.0 => ../Common
