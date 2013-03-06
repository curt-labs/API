package analytics

const (
	ERROR_SEVERITY_SILENCE    = 0
	ERROR_SEVERITY_WARNINGS   = 1
	ERROR_SEVERITY_EXCEPTIONS = 2

	ENDPOINT_HOST = "www.google-analytics.com"
	ENDPOINT_PATH = "/__utm.gif"
)

type Config struct {
	EndpointHost string
	EndpointPath string
}

func NewConfig() Config {
	return Config{
		EndpointHost: ENDPOINT_HOST,
		EndpointPath: ENDPOINT_PATH,
	}
}

func (c *Config) SetHost(host string) {
	if host != "" {
		c.EndpointHost = host
	}
}

func (c *Config) SetPath(path string) {
	if path != "" {
		c.EndpointPath = path
	}
}
