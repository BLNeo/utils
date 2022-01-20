package template

type Http struct {
	RunMode string `toml:"runmode"`
	Address string `toml:"address"`
	Port    string `toml:"port"`
	CorsUrl string `toml:"corsurl"`
}
