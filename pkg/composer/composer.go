package composer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/hermo/registry-check/pkg/resolver"
	"github.com/pkg/errors"
)

// packageDep represents a package in a composer.lock file
type packageDep struct {
	Dist struct {
		URL string `json:"url"`
	} `json:"dist"`
}

//lockFile represents relevant fields in a composer.lock file
type lockFile struct {
	Packages    []packageDep `json:"packages,omitempty"`
	DevPackages []packageDep `json:"packages-dev,omitempty"`
}

//ProcessedLockFile represents a composer.lock file after processing
type ProcessedLockFile struct {
	urls        []string
	numPackages int
}

// NewComposerLockFile creates a new PackageResolver for composer.lock files
func NewComposerLockFile() resolver.PackageResolver {
	return &ProcessedLockFile{urls: []string{}, numPackages: 0}
}

//ReadPackagesFromFile reads a composer.lock file and extracts dependency urls
func (n *ProcessedLockFile) ReadPackagesFromFile(filename string) error {
	rawfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	data := lockFile{}
	err = json.Unmarshal([]byte(rawfile), &data)
	if err != nil {
		return errors.Wrap(err, "Could not parse composer.lock file")
	}

	traverseDeps(&n.urls, &data.Packages)
	traverseDeps(&n.urls, &data.DevPackages)
	return nil
}

// traverseDeps iterates a dependency list and records stored URLs
func traverseDeps(urls *[]string, deps *[]packageDep) {
	for _, dep := range *deps {
		*urls = append(*urls, dep.Dist.URL)
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

// NumPackages reveals how many packages a lockfile contains
func (n *ProcessedLockFile) NumPackages() int {
	return n.numPackages
}
