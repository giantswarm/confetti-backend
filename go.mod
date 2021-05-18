module github.com/giantswarm/confetti-backend

go 1.15

require (
	github.com/atreugo/websocket v1.0.7
	github.com/giantswarm/microerror v0.3.0
	github.com/giantswarm/micrologger v0.5.0
	github.com/savsgio/atreugo/v11 v11.7.0
	github.com/spf13/cobra v1.1.3
)

replace (
	github.com/coreos/etcd => github.com/coreos/etcd v3.3.24+incompatible
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
)
