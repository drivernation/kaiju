package kaiju

type Config struct {
	BindHost string `yaml:"bindHost" json:"bindHost" envconfig:"bind_host" default:"localhost"`
	Port     int    `yaml:"bindPort" json:"bindPort" envconfig:"bind_port" default:"8080"`
}
