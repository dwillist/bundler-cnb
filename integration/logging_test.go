package integration

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testLogging(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
		pack   occam.Pack
		docker occam.Docker
	)

	it.Before(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()
	})

	context("when the buildpack is run with pack build", func() {
		var (
			image occam.Image
			name  string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
		})

		it("logs useful information for the user", func() {
			var err error
			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithNoPull().
				WithBuildpacks(bundlerBuildpack).
				Execute(name, filepath.Join("testdata", "simple_app"))
			Expect(err).ToNot(HaveOccurred(), logs.String)

			buildpackVersion, err := GetGitVersion()
			Expect(err).ToNot(HaveOccurred())

			sequence := []interface{}{
				fmt.Sprintf("Bundler Buildpack %s", buildpackVersion),
				"  Resolving Bundler version",
				"    Candidate version sources (in priority order):",
				"      buildpack.yml -> \"2.1.x\"",
				"",
				MatchRegexp(`    Selected Bundler version \(using buildpack\.yml\): 2\.\d+\.\d+`),
				"",
				"  Executing build process",
				MatchRegexp(`    Installing Bundler 2\.\d+\.\d+`),
				MatchRegexp(`      Completed in \d+\.?\d*`),
			}

			Expect(GetBuildLogs(logs.String())).To(ContainSequence(sequence), logs.String())
		})
	})
}
