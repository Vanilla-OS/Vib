module github.com/vanilla-os/vib

go 1.22.0

toolchain go1.23rc2

require github.com/spf13/cobra v1.8.1

require (
	github.com/containers/storage v1.56.0
	github.com/mitchellh/mapstructure v1.5.0
)

require (
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/Microsoft/hcsshim v0.12.9 // indirect
	github.com/containerd/cgroups/v3 v3.0.3 // indirect
	github.com/containerd/errdefs v0.3.0 // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/containerd/typeurl/v2 v2.2.0 // indirect
	github.com/cyphar/filepath-securejoin v0.3.4 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/go-intervals v0.0.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/klauspost/pgzip v1.2.6 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/mistifyio/go-zfs/v3 v3.0.1 // indirect
	github.com/moby/sys/capability v0.3.0 // indirect
	github.com/moby/sys/mountinfo v0.7.2 // indirect
	github.com/moby/sys/user v0.3.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/runtime-spec v1.2.0 // indirect
	github.com/opencontainers/selinux v1.11.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/tchap/go-patricia/v2 v2.3.1 // indirect
	github.com/ulikunitz/xz v0.5.12 // indirect
	github.com/vbatts/tar-split v0.11.6 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/grpc v1.67.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/vanilla-os/vib/api v0.0.0-20240812130736-2cc767ade5f9
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/vanilla-os/vib/api => ./api
