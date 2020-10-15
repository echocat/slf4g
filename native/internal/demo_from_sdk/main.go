package main

import (
	sdklog "log"

	_ "github.com/echocat/slf4g/native"
	_ "github.com/echocat/slf4g/sdk/bridge/hook"
)

func main() {
	sdklog.Print("abc")
}
