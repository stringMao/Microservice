module GateSvr

go 1.15

require (
	Common v0.0.0
	github.com/Unknwon/goconfig v0.0.0-20200908083735-df7de6a44db8
	github.com/gin-gonic/gin v1.7.2
	github.com/go-xorm/xorm v0.7.9
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo v1.8.5
	github.com/hashicorp/consul/api v1.9.1
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.4 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/protobuf v1.27.1
)

replace Common v0.0.0 => ../Common
