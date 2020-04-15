package main

import (
	"github.com/cloudfoundry/bundler-cnb/bundler"
	"github.com/cloudfoundry/packit"
)

func main() {
	buildpackYMLParser := bundler.NewBuildpackYMLParser()

	packit.Detect(bundler.Detect(buildpackYMLParser))
}
