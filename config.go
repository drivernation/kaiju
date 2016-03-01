package kaiju

type Config struct {
	BindHost string `yaml:"bindHost" json:"bindHost" envconfig:"bind_host" default:"localhost"`
	Port     int    `yaml:"bindPort" json:"bindPort" enconfig:"bind_port" default:"8080"`
}
