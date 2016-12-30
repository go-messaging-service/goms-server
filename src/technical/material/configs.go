package material

type Config struct {
	Path string
}

type TopicConfig struct {
	Config
	Topics []string
}
