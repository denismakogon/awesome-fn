package main

import (
	"context"
)

type invoker string

func (i *invoker) Invoke(ctx context.Context, args ...[]interface{}) ([]byte, error) {
	return []byte("Hello Universe"), nil
}

// exported as symbol named "Greeter"
var Invoker invoker
