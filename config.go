package kaiju

type Config struct {
	BindHost string `yaml:"bindHost" json:"bindHost"`
	Port     int    `yaml:"bindPort" json:"bindPort"`
}
