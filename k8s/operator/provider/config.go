package provider

type ProvidersConfig map[ID]Config

type Config struct {
	Meta map[string]any `mapstructure:"meta,omitempty"`
}

func (c *Config) Get(key string) any {
	return c.Meta[key]
}
