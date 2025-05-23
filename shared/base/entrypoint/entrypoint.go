package entrypoint

type Target struct {
	OS   string
	Arch string
	Name string
}

var DefaultTargets = []Target{
	{OS: "windows", Arch: "amd64", Name: "win-amd64.exe"},
	{OS: "windows", Arch: "386", Name: "win-386.exe"},
	{OS: "linux", Arch: "amd64", Name: "linux-amd64"},
	{OS: "linux", Arch: "386", Name: "linux-386"},
	{OS: "darwin", Arch: "amd64", Name: "darwin-amd64"},
	{OS: "darwin", Arch: "arm64", Name: "darwin-arm64"},
}