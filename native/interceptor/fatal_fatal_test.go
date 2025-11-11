//go:build !darwin

// This test should currently really not run on macOS, because it might
// lead to crashes.

package interceptor

import (
	"os"
	"os/exec"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_fatalExit(t *testing.T) {
	if os.Getenv("DO_IT_NOW_REALLY") == "1" {
		fatalExit(66)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
	cmd.Env = append(os.Environ(), "DO_IT_NOW_REALLY=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok {
		if e.ExitCode() != 66 {
			assert.Failf(t, "Expected to fail with exit code <1>; bot got error: <%+v>", e.ExitCode())
		}
	} else if err != nil {
		assert.Failf(t, "Expected to fail with exit code <1>; bot got error: <%+v>", err)
	} else {
		assert.Fail(t, "Expected to fail with exit code <1>; bot it exists with 0.")
	}
}
