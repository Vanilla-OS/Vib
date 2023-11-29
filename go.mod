module github.com/vanilla-os/vib

go 1.21.4

require github.com/spf13/cobra v1.7.0

require github.com/mitchellh/mapstructure v1.5.0

require (
	github.com/google/uuid v1.3.0
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/vanilla-os/vib/api v0.0.0
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/vanilla-os/vib/api v0.0.0 => ./api/
