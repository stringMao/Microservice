module LoginSvr

go 1.15

require (
	Common v0.0.0
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.1181
	github.com/garyburd/redigo v1.6.2
	github.com/gin-gonic/gin v1.7.2
	github.com/go-xorm/xorm v0.7.9
	github.com/gomodule/redigo v1.8.5
)

replace Common v0.0.0 => ../Common
