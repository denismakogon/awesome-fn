package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"plugin"

	"github.com/fnproject/fdk-go"
	"github.com/google/uuid"
)

const LookUpName string = "Invoker"

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
var client = &http.Client{Transport: tr}

func withError(ctx context.Context, in io.Reader, out io.Writer) {
	res, err := myHandler(ctx, in)
	if err != nil {
		io.WriteString(out, err.Error())
		log.Println(err.Error())
		fdk.WriteStatus(out, http.StatusInternalServerError)
		return
	}
	out.Write(res)
	fdk.WriteStatus(out, http.StatusOK)
}

type SharedObjectURL struct {
	URL string `json:"url"`
}

func myHandler(ctx context.Context, in io.Reader) ([]byte, error) {
	var so SharedObjectURL
	err := json.NewDecoder(in).Decode(&so)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(so.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	libPath, err := ReadSharedFile(ctx, resp.Body, uuid.New().String())
	if err != nil {
		return nil, err
	}

	defer os.Remove(libPath)

	pl, err := plugin.Open(libPath)
	if err != nil {
		return nil, err
	}

	log.Println("plugin loaded")

	return Invoke(ctx, pl)
}

func main() {
	fdk.Handle(fdk.HandlerFunc(withError))
}
