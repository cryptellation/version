// A generated module for Version functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"runtime"

	"github.com/cryptellation/version/dagger/internal/dagger"
)

type Version struct{}

// Publish a new release
func (ci *Version) PublishTag(
	ctx context.Context,
	sourceDir *dagger.Directory,
	token *dagger.Secret,
) error {
	// Create Git repo access
	repo, err := NewGit(ctx, sourceDir, token)
	if err != nil {
		return err
	}

	// Publish new tag
	return repo.PublishTagFromReleaseTitle(ctx)
}

// Lint runs golangci-lint on the source code in the given directory.
func (mod *Version) Lint(sourceDir *dagger.Directory) *dagger.Container {
	c := dag.Container().
		From("golangci/golangci-lint:v1.62.0").
		WithMountedCache("/root/.cache/golangci-lint", dag.CacheVolume("golangci-lint"))

	c = mod.withGoCodeAndCacheAsWorkDirectory(c, sourceDir)

	return c.WithExec([]string{"golangci-lint", "run", "--timeout", "10m"})
}

// UnitTests returns a container that runs the unit tests.
func (mod *Version) UnitTests(sourceDir *dagger.Directory) *dagger.Container {
	c := dag.Container().From("golang:" + goVersion() + "-alpine")
	return mod.withGoCodeAndCacheAsWorkDirectory(c, sourceDir).
		WithExec([]string{"sh", "-c",
			"go test ./...",
		})
}

func goVersion() string {
	return runtime.Version()[2:]
}

func (mod *Version) withGoCodeAndCacheAsWorkDirectory(
	c *dagger.Container,
	sourceDir *dagger.Directory,
) *dagger.Container {
	containerPath := "/go/src/github.com/cryptellation/version"
	return c.
		// Add Go caches
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("gobuild")).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("gocache")).

		// Add source code
		WithMountedDirectory(containerPath, sourceDir).

		// Add workdir
		WithWorkdir(containerPath)
}
