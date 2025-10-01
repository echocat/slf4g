package main

import (
	sdklog "log"

	"github.com/echocat/slf4g/native/location"

	_ "github.com/echocat/slf4g/hooks/sdklog"
	_ "github.com/echocat/slf4g/native"
)

func main() {
	location.DefaultDiscovery = location.NewCallerDiscovery()

	sdklog.Print("abc")
}
