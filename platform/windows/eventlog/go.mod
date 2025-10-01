module github.com/echocat/slf4g/platform/windows/eventlog

go 1.23.0

toolchain go1.24.1

replace (
	github.com/echocat/slf4g => ../../../
	github.com/echocat/slf4g/native => ../../../native
)

require (
	github.com/echocat/slf4g v0.0.0
	github.com/echocat/slf4g/native v0.0.0
	golang.org/x/sys v0.35.0
)
