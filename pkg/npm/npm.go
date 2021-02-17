package npm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/hermo/registry-check/pkg/resolver"
	"github.com/pkg/errors"
)

// packageLockJSON represents a package-lock.json for determining it's version
type packageLockJSON struct {
	LockFileVersion int `json:"lockFileVersion"`
}

// lockFileV1 represents a v1 package-lock.json
type lockFileV1 struct {
	Dependencies map[string]lockFileV1Dependency `json:"dependencies,omitempty"`
}

// lockFileV2 represents a v2 package-lock.json
type lockFileV2 struct {
	Dependencies map[string]lockFileV1Dependency `json:"dependencies,omitempty"`
	Packages     map[string]lockFileV2Package    `json:"packages,omitempty"`
}

// lockFileV1Dependency represents the recursive structure of v1 package-lock.json dependencies
type lockFileV1Dependency struct {
	Resolved     string                          `json:"resolved,omitempty"`
	Dependencies map[string]lockFileV1Dependency `json:"dependencies,omitempty"`
}

// lockFileV2Package represents a package in the flat format of v2 package-lock.json files
type lockFileV2Package struct {
	Resolved string `json:"resolved,omitempty"`
}

// ProcessedLockFile represents a package-lock.json file after processing
type ProcessedLockFile struct {
	urls        []string
	numPackages int
}

// NewPackageLockFile creates a new PackageResolver for package-lock.json files
func NewPackageLockFile() resolver.PackageResolver {
	return &ProcessedLockFile{urls: []string{}, numPackages: 0}
}

// ReadPackagesFromFile reads a package-lock.json file and extracts dependency urls
func (n *ProcessedLockFile) ReadPackagesFromFile(filename string) error {
	rawfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	data := packageLockJSON{}
	if err := json.Unmarshal([]byte(rawfile), &data); err != nil {
		return errors.Wrap(err, "Could not parse package-lock.json file")
	}

	switch data.LockFileVersion {
	case 1:
		var parsed lockFileV1
		if err := json.Unmarshal([]byte(rawfile), &parsed); err != nil {
			return errors.Wrap(err, "Could not parse v1 package-lock.json file")
		}
		traverseDeps(&n.urls, &parsed.Dependencies)
	case 2:
		var parsed lockFileV2
		if err := json.Unmarshal([]byte(rawfile), &parsed); err != nil {
			return errors.Wrap(err, "Could not parse v2 package-lock.json file")
		}
		traverseDeps(&n.urls, &parsed.Dependencies)
		traversePackages(&n.urls, &parsed.Packages)

	default:
		return fmt.Errorf("Unsupported lockfile version %d", data.LockFileVersion)
	}
	return nil
}

// traverseDeps recursively iterates through the "dependencies" key of a package-lock.json file
func traverseDeps(urls *[]string, deps *map[string]lockFileV1Dependency) {
	for _, dep := range *deps {
		*urls = append(*urls, dep.Resolved)
		traverseDeps(urls, &dep.Dependencies)
	}
}

// traverseDeps iterates through the "packages" key of a v2 package-lock.json file
func traversePackages(urls *[]string, deps *map[string]lockFileV2Package) {
	for _, dep := range *deps {
		*urls = append(*urls, dep.Resolved)
	}
}

// Registries returns a list of all registries used by a package
func (n *ProcessedLockFile) Registries() []string {
	registries := make(map[string]bool)
	var out []string
	for _, s := range n.urls {
		if s == "" {
			continue
		}
		n.numPackages++
		url, err := url.Parse(s)
		if err != nil {
			panic(err)
		}

		key := fmt.Sprintf("%s://%s", url.Scheme, url.Host)
		if !registries[key] {
			registries[key] = true
			out = append(out, key)
		}
	}
	return out
}

// NumPackages reveals how many (not necessarily unique) packages a lockfile contains
func (n *ProcessedLockFile) NumPackages() int {
	return n.numPackages
}
