package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fnproject/fdk-go"
)

var debug, _ = strconv.Atoi(os.Getenv("FDK_DEBUG"))
var bufPool = &sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}

func logMessage(args ...interface{}) {
	if debug == 1 {
		log.Println(args...)
	}
}

func azureHandler(ctx context.Context, httpClient *http.Client,
	proxyRequest *http.Request, in io.Reader, out io.Writer) error {

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)

	logMessage("in handler")
	fctx := fdk.Context(ctx)
	logMessage("context provisioned")
	contentLen, err := io.Copy(buf, in)
	if err != nil {
		return err
	}

	logMessage("content len: ", contentLen)
	logMessage("^ here should be an income request body and its length")
	logMessage("HTTP method: ", fctx.Method)

	if err != nil {
		return err
	}

	proxyRequest.Body = ioutil.NopCloser(buf)
	for k, v := range fctx.Header {
		proxyRequest.Header[k] = v
	}
	proxyRequest.ContentLength = contentLen

	logMessage("headers: ", fctx.Header)
	logMessage("proxy request provisioned")

	resp, err := httpClient.Do(proxyRequest)
	logMessage("proxy request submitted")
	if err != nil {
		return err
	}

	logMessage("request was successful")
	defer resp.Body.Close()

	io.Copy(out, resp.Body)
	logMessage("setting headers back")

	for k, v := range resp.Header {
		fdk.SetHeader(out, k, strings.Join(v, ";"))
	}

	return nil
}

func SelectAzureFuncHandler(azureFunc string, httpClient *http.Client) fdk.HandlerFunc {

	proxyReq, _ := http.NewRequest(
		"POST", "http://localhost:80/api/"+azureFunc, nil)

	return func(ctx context.Context, in io.Reader, out io.Writer) {
		err := azureHandler(ctx, httpClient, proxyReq, in, out)
		if err != nil {
			io.WriteString(out, err.Error())
			fdk.WriteStatus(out, http.StatusInternalServerError)
			return
		}
		fdk.WriteStatus(out, http.StatusOK)
	}
}

func main() {

	cmd := exec.Command("./run_host.sh")
	cmd.Stdout = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	azureFunc := os.Getenv("AZURE_FUNCTION")
	if azureFunc == "" {
		panic("unable to find an Azure function to run")
	}

	httpClient := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 5 * time.Second,
			DisableKeepAlives:   false,
		},
	}

	for i := 0; i < 100; i++ {
		res, err := httpClient.Post(
			"http://localhost:80", "text/plain", nil)
		if err != nil {
			logMessage(err.Error())
		} else {

		}
		if res != nil {
			if res.StatusCode <= 202 {
				logMessage("Azure serverless app is running fine.")
				fdk.Handle(SelectAzureFuncHandler(azureFunc, httpClient))
			}
		}
		time.Sleep(1 * time.Second)
	}
}
