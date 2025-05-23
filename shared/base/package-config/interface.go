package packageconfig

type PackageConfig struct {
	URL         string   `json:"url"`
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Maintainers []string `json:"maintainers"`
	Icon        string   `json:"icon"`
	NetBridge   struct {
		MaxPorts int `json:"maxports"`
	} `json:"netbridge"`
	BuildNumber string `json:"buildNumber,omitempty"`
}
