package resolver

// PackageResolver finds urls and registries from a package file
type PackageResolver interface {
	ReadPackagesFromFile(string) error
	Registries() []string
	NumPackages() int
}
