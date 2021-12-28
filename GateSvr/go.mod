module GateSvr

go 1.15

require (
	Common v0.0.0
	github.com/garyburd/redigo v1.6.2
	github.com/gin-gonic/gin v1.7.2
	github.com/go-xorm/xorm v0.7.9
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo v1.8.5
	github.com/hashicorp/consul/api v1.9.1
	github.com/smartystreets/goconvey v1.7.2 // indirect
)

replace Common v0.0.0 => ../Common
