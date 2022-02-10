module grpc_client

go 1.16

require (
	github.com/asim/go-micro/plugins/client/grpc/v3 v3.7.0
	github.com/asim/go-micro/plugins/registry/consul/v3 v3.7.0
	github.com/asim/go-micro/v3 v3.7.0
	github.com/sirupsen/logrus v1.7.0
	message v0.0.0
)

replace message => ../message
