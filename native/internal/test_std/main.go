package main

import (
	_ "github.com/echocat/slf4g/bridge-std/hook"
	_ "github.com/echocat/slf4g/native"
	stdlog "log"
)

func main() {
	stdlog.Print("abc")
}
