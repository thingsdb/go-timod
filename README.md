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

The following code may be used as a template: (see: https://github.com/thingsdb/ThingsDB/tree/master/modules/go/demo)

```go
package main

import (
    "fmt"
    "log"

    timod "github.com/thingsdb/go-timod"

    "github.com/vmihailenco/msgpack"
)

func handler(buf *timod.Buffer, quit chan bool) {
    for {
        select {
        case pkg := <-buf.PkgCh:
            switch timod.Proto(pkg.Tp) {
            case timod.ProtoModuleConf:
                // Configuration data for this module is received from ThingsDB.
                //
                // The module should respond with:
                //
                // - timod.WriteConfOk(): if successful
                // - timod.WriteConfErr(): in case the configuration has failed
                log.Println("No configuration data is required for this module")
                timod.WriteConfOk() // Just write OK

            case timod.ProtoModuleReq:
                // A request from ThingsDB may be unpacked to a struct or to
                // an map[string]interface{}.
                //
                // The module should respond with:
                //
                // - timod.WriteResponse(pid, value): if successful
                // - timod.WriteEx(pid, err_code, err_msg): in case of an error
                //
                // (alternative to WriteResponse, WriteResponseRaw is available
                //  for a pre-packed, raw response)
                type Demo struct {
                    Message string `msgpack:"message"`
                }
                var demo Demo

                // pkg.Data contains Message Packed data, most likely you want
                // to unpack the data into a struct.
                err := msgpack.Unmarshal(pkg.Data, &demo)
                if err == nil {
                    // In this demo a `message` property will be unpacked and
                    // used as a return value.
                    timod.WriteResponse(pkg.Pid, &demo.Message)
                } else {
                    // In case of an error, make sure to call `WriteEx(..)` so
                    // ThingsDB can finish the future request with an
                    // appropriate error. (see error codes below)
                    timod.WriteEx(
                        pkg.Pid,
                        timod.ExBadData,
                        fmt.Sprintf("failed to unpack request (%s)", err))
				}

            default:
                log.Printf("Error: Unexpected package type: %d", pkg.Tp)
			}
        case err := <-buf.ErrCh:
            // In case of an error you probably want to quit the module.
            // ThingsDB will try to restart the module a few times if this
            // happens.
            log.Printf("Error: %s", err)
            quit <- true
        }
    }
}

func main() {
    // Starts the module
    timod.StartModule("demo", handler)

    // It is possible to add some cleanup code here
}
```

## Exceptions

These are the possible exceptions which may be used. Do not use any other
exception code, as ThingsDB will only accept values withing the given range.

- `ExCancelled` - operation is cancelled before completion
- `ExOperation` - operation is not valid in the current context
- `ExNumArguments` - wrong number of arguments
- `ExTypeError` - object of inappropriate type
- `ExValueError` - object has the right type but an inappropriate value
- `ExOverflow` - integer overflow
- `ExZeroDiv` - division or module by zero
- `ExMaxQuota` - max quota is reached
- `ExAuthError` - authentication error
- `ExForbidden` - forbidden (access denied)
- `ExLookupError` - requested resource not found
- `ExBadData` - unable to handle request due to invalid data
- `ExSyntaxError` - syntax error in query
- `ExNodeError` - node is temporary unable to handle the request
- `ExAssertError` - assertion statement has failed
- `ExCustom100`..`ExCustom127` - can be used as a custom errors
