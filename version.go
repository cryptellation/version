package version

import (
	"fmt"
)

var (
	// Version of the application.
	globalVersion string

	// Revision of the application.
	globalCommitHash string
)

// Version returns a string representing the current version.
func Version() string {
	return globalVersion
}

// SetVersion sets the version to tha value provided as ver unless ver is empty.
func SetVersion(ver string) {
	if ver != "" {
		globalVersion = ver
	}
}

// SetCommitHash sets the commit hash of the application to the value provided as hash.
// Empty values are accepted.
func SetCommitHash(hash string) {
	globalCommitHash = hash
}

// CommitHash returns a string representing the current commitHash.
func CommitHash() string {
	return globalCommitHash
}

// FullVersion returns a string representing the version and commit hash concatenated separated by a '-'.
//
// Returns only the version if the commit hash is not defined.
func FullVersion() string {
	switch {
	case globalVersion == "" && globalCommitHash == "":
		return "unknown"
	case globalVersion == "":
		return globalCommitHash
	case globalCommitHash == "":
		return globalVersion
	default:
		return fullVersion(globalVersion, globalCommitHash)
	}
}

func fullVersion(version, commitHash string) string {
	return fmt.Sprintf("%s-%s", version, commitHash)
}
