package conf

import "reflect"

type Config struct {
	MapEnv bool
	Types  TypesTable
	Expect reflect.Kind
}

// Option for configuring config.
type Option func(c *Config)
