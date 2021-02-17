package presentation

// Output contains all the things the app might want to output
type Output struct {
	ParsedSuccessfully bool
	Error              error
	Filename           string
	LockfileType       string
	NumPackages        int
	Registries         []string
}

// Presenter presents output
type Presenter interface {
	Present(output *Output)
}
