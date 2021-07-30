module TestClient

go 1.15

replace Common v0.0.0 => ../Common

require (
	Common v0.0.0
	github.com/golang/protobuf v1.5.0
)
