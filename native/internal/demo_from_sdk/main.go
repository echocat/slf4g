package main

import (
	sdklog "log"

	"github.com/echocat/slf4g/native/location"

	_ "github.com/echocat/slf4g/native"
	_ "github.com/echocat/slf4g/sdk/bridge/hook"
)

func main() {
	location.DefaultDiscovery = location.NewCallerDiscovery()

	sdklog.Print("abc")
}
