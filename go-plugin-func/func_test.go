package main

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
)

var fpath = "sample/linux_invoker.so"

func TestReadSharedAndInvoke(t *testing.T) {
	t.Run("test-invoke", func(t *testing.T) {
		ctx := context.Background()
		sharedLib, err := os.Open(fpath)
		if err != nil {
			t.Fatal(err.Error())
		}
		sharedLibPrefix := uuid.New().String()
		libPath, pl, err := ReadSharedFile(
			ctx, sharedLib, sharedLibPrefix)
		if err != nil {
			t.Fatal(err.Error())
		}
		defer os.Remove(libPath)

		res, err := Invoke(ctx, pl)
		if err != nil {
			t.Fatal(err.Error())
		}

		expected := "Hello Universe"
		if expected != string(res) {
			t.Fatalf("assertion error, expected: %v, actual: %v",
				expected, string(res))
		}
	})
}
