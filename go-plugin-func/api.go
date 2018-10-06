package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"plugin"
)

type InvokeFromShared interface {
	Invoke(ctx context.Context, args ...[]interface{}) ([]byte, error)
}

func ReadSharedFile(_ context.Context, content io.Reader, identifier string) (string, error) {
	log.Println("in ReadSharedFile")
	dirPath := os.TempDir()
	sharedFile, err := ioutil.TempFile(
		dirPath, fmt.Sprintf("%v-*.so", identifier))
	if err != nil {
		return "", err
	}
	defer sharedFile.Close()

	fileName := sharedFile.Name()
	log.Println(fmt.Sprintf("temp file %s created", fileName))

	_, err = io.Copy(sharedFile, content)
	return fileName, err
}

func Invoke(ctx context.Context, pl *plugin.Plugin, args ...[]interface{}) ([]byte, error) {
	log.Println("in Invoke")
	typeDef, err := pl.Lookup(LookUpName)
	if err != nil {
		return nil, err
	}
	log.Println(fmt.Sprintf("Lookup for %s was successful", LookUpName))

	invoker, ok := typeDef.(InvokeFromShared)
	if !ok {
		return nil, errors.New("unable to cast exposed type to Invoker type")
	}
	log.Println("Type cast was successful")
	// add parameters somehow, maybe an HTTP query parameters?
	res, err := invoker.Invoke(ctx, args...)
	log.Println("Invoke was finished")
	return res, err
}
