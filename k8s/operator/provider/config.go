package provider

type ProvidersConfig map[ID]Config

type Config struct {
	Meta map[string]any `json:"meta,omitempty" yaml:"meta,omitempty"`
}

func (c *Config) Get(key string) any {
	return c.Meta[key]
}
