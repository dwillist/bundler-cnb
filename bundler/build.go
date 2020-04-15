package bundler

import (
	"path/filepath"
	"time"

	"github.com/cloudfoundry/packit"
	"github.com/cloudfoundry/packit/postal"
)

//go:generate faux --interface EntryResolver --output fakes/entry_resolver.go
type EntryResolver interface {
	Resolve([]packit.BuildpackPlanEntry) packit.BuildpackPlanEntry
}

//go:generate faux --interface DependencyManager --output fakes/dependency_manager.go
type DependencyManager interface {
	Resolve(path, id, version, stack string) (postal.Dependency, error)
	Install(dependency postal.Dependency, cnbPath, layerPath string) error
}

//go:generate faux --interface BuildPlanRefinery --output fakes/build_plan_refinery.go
type BuildPlanRefinery interface {
	BillOfMaterial(dependency postal.Dependency) packit.BuildpackPlan
}

func Build(entries EntryResolver, dependencies DependencyManager, planRefinery BuildPlanRefinery, logger LogEmitter, clock Clock) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)
		logger.Process("Resolving Bundler version")

		entry := entries.Resolve(context.Plan.Entries)

		dependency, err := dependencies.Resolve(filepath.Join(context.CNBPath, "buildpack.toml"), entry.Name, entry.Version, context.Stack)
		if err != nil {
			return packit.BuildResult{}, err
		}

		logger.SelectedDependency(entry, dependency, clock.Now())

		bundlerLayer, err := context.Layers.Get(Bundler, packit.LaunchLayer)
		if err != nil {
			return packit.BuildResult{}, err
		}

		bundlerLayer.Build = entry.Metadata["build"] == true
		bundlerLayer.Cache = entry.Metadata["build"] == true

		bom := planRefinery.BillOfMaterial(postal.Dependency{
			ID:      dependency.ID,
			Name:    dependency.Name,
			SHA256:  dependency.SHA256,
			Stacks:  dependency.Stacks,
			URI:     dependency.URI,
			Version: dependency.Version,
		})

		cachedSHA, ok := bundlerLayer.Metadata[DepKey].(string)
		if ok && cachedSHA == dependency.SHA256 {
			logger.Process("Reusing cached layer %s", bundlerLayer.Path)
			logger.Break()

			return packit.BuildResult{
				Plan:   bom,
				Layers: []packit.Layer{bundlerLayer},
			}, nil
		}

		logger.Process("Executing build process")

		err = bundlerLayer.Reset()
		if err != nil {
			return packit.BuildResult{}, err
		}

		bundlerLayer.Metadata = map[string]interface{}{
			DepKey:     dependency.SHA256,
			"built_at": clock.Now().Format(time.RFC3339Nano),
		}

		logger.Subprocess("Installing Bundler %s", dependency.Version)
		then := clock.Now()
		err = dependencies.Install(dependency, context.CNBPath, bundlerLayer.Path)
		if err != nil {
			return packit.BuildResult{}, err
		}
		logger.Action("Completed in %s", time.Since(then).Round(time.Millisecond))
		logger.Break()

		return packit.BuildResult{
			Plan:   bom,
			Layers: []packit.Layer{bundlerLayer},
		}, nil
	}
}
