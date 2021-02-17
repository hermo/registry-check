package analyzer

import (
	"fmt"
	"path"

	"github.com/hermo/registry-check/pkg/composer"
	"github.com/hermo/registry-check/pkg/npm"
	"github.com/hermo/registry-check/pkg/resolver"
)

// Result represents the findings of a AnalyzeLockfile call
type Result struct {
	LockfileType string
	NumPackages  int
	Registries   []string
}

// DetermineLockfileType determines what type of lock file filename represents
func DetermineLockfileType(filename string) (string, error) {
	switch fn := path.Base(filename); fn {
	case "package-lock.json":
		return "npm", nil
	case "composer.lock":
		return "composer", nil
	default:
		return "", fmt.Errorf("Cannot guess lockfile format from filename")
	}
}

// AnalyzeLockfile finds out what registries are used by a lockfile among other things
func AnalyzeLockfile(filename string, lockfileType string) (result Result, err error) {
	var resolver resolver.PackageResolver

	switch lockfileType {
	case "npm":
		resolver = npm.NewPackageLockFile()
	case "composer":
		resolver = composer.NewComposerLockFile()
	default:
		err = fmt.Errorf("Unsupported package type: %s", lockfileType)
		return
	}

	if err = resolver.ReadPackagesFromFile(filename); err != nil {
		return // nil, fmt.Errorf("ERROR: Can't read file %w", err)
	}
	result.Registries = resolver.Registries()
	result.NumPackages = resolver.NumPackages()
	result.LockfileType = lockfileType
	return
}
