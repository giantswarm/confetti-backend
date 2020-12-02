module github.com/giantswarm/confetti-backend

go 1.15

require (
	github.com/atreugo/websocket v1.0.6
	github.com/giantswarm/microerror v0.2.1
	github.com/giantswarm/micrologger v0.4.0
	github.com/savsgio/atreugo/v11 v11.5.4
	github.com/spf13/cobra v1.1.1
)

replace github.com/coreos/etcd => github.com/coreos/etcd v3.3.24+incompatible
