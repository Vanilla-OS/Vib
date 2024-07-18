module github.com/vanilla-os/vib

go 1.22

require github.com/spf13/cobra v1.8.1

require (
	cuelang.org/go v0.9.2
	github.com/ebitengine/purego v0.7.1
	github.com/mitchellh/mapstructure v1.5.0
)

require (
	github.com/cockroachdb/apd/v3 v3.2.1 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.15.0 // indirect
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/vanilla-os/vib/api v0.0.0-20240618053016-44e9ee99064a
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/vanilla-os/vib/api => ./api
