package main

import (
	"context"
	"strings"

	"github.com/cryptellation/version/dagger/internal/dagger"
)

// Git provides access to a git repository.
type Git struct {
	container *dagger.Container
}

// NewGit creates a new Git container with the given source directory and token.
func NewGit(ctx context.Context, srcDir *dagger.Directory, token *dagger.Secret) (Git, error) {
	var err error

	// Create container
	container := dag.Container().
		From("alpine/git").
		WithMountedDirectory("/git", srcDir).
		WithWorkdir("/git").
		WithoutEntrypoint()

	// Set authentication based on the token
	tokenString, err := token.Plaintext(ctx)
	if err != nil {
		return Git{}, err
	}

	// Change the url to use the token
	container, err = container.WithExec([]string{
		"git", "remote", "set-url", "origin", "https://lerenn:" + tokenString + "@github.com/cryptellation/version.git",
	}).Sync(ctx)
	if err != nil {
		return Git{}, err
	}

	// Set Git author
	container, err = setGitAuthor(ctx, container)
	if err != nil {
		return Git{}, err
	}

	return Git{
		container: container,
	}, nil
}

// GetLastCommitTitle returns the title of the last commit.
func (g *Git) GetLastCommitTitle(ctx context.Context) (string, error) {
	res, err := g.container.
		WithExec([]string{"git", "log", "-1", "--pretty=%B"}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}

	// Remove potential new line
	res = strings.TrimSuffix(res, "\n")

	return res, nil
}

// PublishTagFromReleaseTitle creates a new tag based on the last commit title.
// The title should follow the angular convention.
func (g *Git) PublishTagFromReleaseTitle(ctx context.Context) error {
	// Get new semver
	title, err := g.GetLastCommitTitle(ctx)
	if err != nil {
		return err
	}

	// Get last tag
	lastTag, err := g.GetLastTag(ctx)
	if err != nil {
		return err
	}

	// Process newSemVer change
	change, newSemVer, err := ProcessSemVerChange(lastTag, title)
	if err != nil {
		return err
	}
	if change == SemVerChangeNone {
		return nil
	}

	// Tag commit
	g.container, err = g.container.
		WithExec([]string{"git", "tag", "v" + newSemVer}).
		Sync(ctx)
	if err != nil {
		return err
	}

	// Push new tag
	g.container, err = g.container.
		WithExec([]string{"git", "push", "--tags"}).
		Sync(ctx)

	return err
}

// GetLastTag returns the last tag of the repository.
func (g *Git) GetLastTag(ctx context.Context) (string, error) {
	res, err := g.container.
		WithExec([]string{"git", "describe", "--tags", "--abbrev=0"}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}

	// Remove potential new line
	res = strings.TrimSuffix(res, "\n")

	return res, nil
}

func setGitAuthor(
	ctx context.Context,
	container *dagger.Container,
) (*dagger.Container, error) {
	// Add infos on author
	container, err := container.
		WithExec([]string{"git", "config", "--global", "user.email", "louis.fradin+cryptellation-ci@gmail.com"}).
		Sync(ctx)
	if err != nil {
		return nil, err
	}
	container, err = container.
		WithExec([]string{"git", "config", "--global", "user.name", "Cryptellation CI"}).
		Sync(ctx)
	if err != nil {
		return nil, err
	}

	return container, nil
}
