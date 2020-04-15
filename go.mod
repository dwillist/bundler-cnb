module github.com/cloudfoundry/bundler-cnb

require (
	github.com/buildpack/libbuildpack v1.25.11 // indirect
	github.com/cloudfoundry/dagger v0.0.0-20200213200846-c2a9723f08c4
	github.com/cloudfoundry/libcfbuildpack v1.91.23 // indirect
	github.com/cloudfoundry/node-engine-cnb v0.0.169 // indirect
	github.com/cloudfoundry/occam v0.0.0-20200227222407-1c4318bd26e2
	github.com/cloudfoundry/packit v0.0.0-20200316192019-270fa885ab52
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/golang/protobuf v1.3.5 // indirect
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/onsi/gomega v1.9.0
	github.com/sclevine/spec v1.4.0
	golang.org/x/net v0.0.0-20200320220750-118fecf932d8 // indirect
	golang.org/x/sys v0.0.0-20200321134203-328b4cd54aae // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.2.8
)

go 1.13

replace github.com/cloudfoundry/packit => /home/arch/workspace/Pivotal/packit
