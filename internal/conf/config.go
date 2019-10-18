package conf

import (
	"fmt"
	"reflect"
)

type Config struct {
	MapEnv    bool
	Types     TypesTable
	Operators OperatorsTable
	Expect    reflect.Kind
	Optimize  bool
}

func New(i interface{}) *Config {
	var mapEnv bool
	if _, ok := i.(map[string]interface{}); ok {
		mapEnv = true
	}

	return &Config{
		MapEnv:   mapEnv,
		Types:    CreateTypesTable(i),
		Optimize: true,
	}
}

// Check validates the compiler configuration.
func (c *Config) Check() error {
	// check that all functions that define operator overloading
	// exist in environment and have correct signatures.
	for op, fns := range c.Operators {
		for _, fn := range fns {
			fnType, ok := c.Types[fn]
			if !ok || fnType.Type.Kind() != reflect.Func {
				return fmt.Errorf("function %s for %s operator does not exist in environment", fn, op)
			}
			requiredNumIn := 2
			if fnType.Method {
				requiredNumIn = 3 // As first argument of method is receiver.
			}
			if fnType.Type.NumIn() != requiredNumIn || fnType.Type.NumOut() != 1 {
				return fmt.Errorf("function %s for %s operator does not have a correct signature", fn, op)
			}
		}
	}
	return nil
}
