package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/fnproject/fdk-go"
	"github.com/google/uuid"
)

const LookUpName string = "Invoker"

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

func myHandler(ctx context.Context, in io.Reader) ([]byte, error) {

	libPath, err := ReadSharedFile(ctx, in, uuid.New().String())
	if err != nil {
		return nil, err
	}

	defer os.Remove(libPath)
	_, err = os.Open(libPath)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("file exists and i can open it")
	}
	cmd := exec.CommandContext(ctx, "file", libPath)
	b, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	log.Println(string(b))
	//
	//pl, err := plugin.Open(libPath)
	//if err != nil {
	//	return nil, err
	//}
	//
	//log.Println("plugin loaded")
	//
	//return Invoke(ctx, pl)
	return nil, err
}

func main() {
	fdk.Handle(fdk.HandlerFunc(withError))
}
