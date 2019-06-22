package conf

type Config struct {
	MapEnv bool
	Types  TypesTable
}

// Option for configuring config.
type Option func(c *Config)
