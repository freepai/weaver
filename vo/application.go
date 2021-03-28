package vo

type Application struct {
	Name string `yaml:"name"`
	Group string `yaml:"group"`
	Resources []string `yaml:"resources"`
}
