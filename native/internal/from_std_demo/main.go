package main

import (
	stdlog "log"

	_ "github.com/echocat/slf4g/native"
	_ "github.com/echocat/slf4g/sdk/bridge/hook"
)

func main() {
	stdlog.Print("abc")
}
