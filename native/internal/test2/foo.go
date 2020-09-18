package main

import (
	log "github.com/echocat/slf4g"
	_ "github.com/echocat/slf4g/native"
	stdlog "log"
)

func main() {
	log.ConfigureStd()
	stdlog.Print("abc")
}
