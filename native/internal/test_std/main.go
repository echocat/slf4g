package main

import (
	_ "github.com/echocat/slf4g/native"
	_ "github.com/echocat/slf4g/std/bridge/hook"
	stdlog "log"
)

func main() {
	stdlog.Print("abc")
}
