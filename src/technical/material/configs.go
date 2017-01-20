package technicalMaterial

type Config struct {
	Path string
}

type TopicConfig struct {
	Config
	Topics []string `json:"topics"`
}

type ServerConfig struct {
	TopicLocation string      `json:"topic-config"`
	Connectors    []Connector `json:"connectors"`
}

type Connector struct {
	protocol string
	ip       string
	port     string
}
