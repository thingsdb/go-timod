# timod

Library for creating ThingsDB modules using the Go language

## Installation

Simple install the package to your [$GOPATH](https://github.com/golang/go/wiki/GOPATH) with the [go tool](https://golang.org/cmd/go/) from shell:

```shell
$ go get github.com/thingsdb/go-timod
```

Make sure [Git](https://git-scm.com/downloads) is installed on your machine and in your system's PATH.

## Usage

Modules for ThingsDB are simple binary files which should read from `stdin` and write a response back to `stdout`. Work by the module must be non-blocking. When a request is received from ThingsDB, the data contains a package id (`Pid`) which should be used for the response. This is used by ThingsDB to map the response to the correct request since responses do not have to be written back to ThingsDB in order.

If the module requires configuration data, for example a connection string, then this configuration will be send immediately after start-up but may be received again when the module configuration is changed in ThingsDB.

Do not use functions like `Println` and `Printf` since these function will write to `stdout` and this is reserved for ThingsDB. Instead, use `log.Print..` to write to `stderr` instead.

The following code may be used as a template:

```
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	timod "requests/timod"

	"github.com/vmihailenco/msgpack"
)

type reqData struct {
	URL     string      `msgpack:"url"`
	Method  string      `msgpack:"metohd"`
	Body    []byte      `msgpack:"body"`
	Headers [][2]string `msgpack:"headers"`
	Params  [][2]string `msgpack:"params"`
}

func handleReqData(pkg *timod.Pkg, data *reqData) {

	params := url.Values{}

	reqURL, err := url.Parse(data.URL)
	if err != nil {
		timod.WriteEx(pkg.Pid, timod.ExBadData, fmt.Sprintf("failed to parse URL (%s)", err))
		return
	}

	for i := 0; i < len(data.Params); i++ {
		param := data.Params[i]
		key, value := param[0], param[1]
		params.Set(key, value)
	}

	reqURL.RawQuery = params.Encode()

	req, err := http.NewRequest(data.Method, reqURL.String(), nil)
	if err != nil {
		timod.WriteEx(pkg.Pid, timod.ExBadData, fmt.Sprintf("failed to create HTTP request (%s)", err))
		return
	}

	for i := 0; i < len(data.Headers); i++ {
		header := data.Headers[i]
		key, value := header[0], header[1]
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		timod.WriteEx(pkg.Pid, timod.ExBadData, fmt.Sprintf("failed to do the HTTP request (%s)", err))
		return
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		timod.WriteEx(pkg.Pid, timod.ExBadData, fmt.Sprintf("failed to read bytes from HTTP response (%s)", err))
		return
	}

	// timod.WriteResponse(pkg.Pid, &resBytes)
}

func onModuleReq(pkg *timod.Pkg) {
	var data reqData

	err := msgpack.Unmarshal(pkg.Data, &data)
	if err == nil {
		handleReqData(pkg, &data)
	} else {
		timod.WriteEx(pkg.Pid, timod.ExBadData, fmt.Sprintf("failed to unpack request (%s)", err))
	}
}

func handler(buf *timod.Buffer, quit chan bool) {
	for {
		select {
		case pkg := <-buf.PkgCh:
			switch timod.Proto(pkg.Tp) {
			case timod.ProtoModuleInit:
				log.Println("No init required for this module")
			case timod.ProtoModuleReq:
				onModuleReq(pkg)
			default:
				log.Printf("Error: Unexpected package type: %d", pkg.Tp)
			}
		case err := <-buf.ErrCh:
			log.Printf("Error: %s", err)
			quit <- true
			return
		}
	}
}

func main() {
    // Starts the module
	timod.StartModule("demo", handler)
}

```