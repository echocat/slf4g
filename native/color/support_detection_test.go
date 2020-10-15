package color

import (
	"os"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_CanSupportBeAssumed(t *testing.T) {
	oldAssumptions := SupportAssumptionDetections
	defer func() {
		SupportAssumptionDetections = oldAssumptions
	}()

	var givenResult bool
	SupportAssumptionDetections = []SupportAssumptionDetection{func() bool {
		return givenResult
	}}

	givenResult = false
	assert.ToBeEqual(t, givenResult, CanSupportBeAssumed())

	givenResult = true
	assert.ToBeEqual(t, givenResult, CanSupportBeAssumed())
}

func Test_SupportAssumptionDetectionIntellij(t *testing.T) {
	defer prepareEnv("IDEA_INITIAL_DIRECTORY", "TERMINAL_EMULATOR")()

	assert.ToBeEqual(t, false, SupportAssumptionDetectionIntellij())

	setEnv("IDEA_INITIAL_DIRECTORY", "foo")
	assert.ToBeEqual(t, true, SupportAssumptionDetectionIntellij())

	setEnv("IDEA_INITIAL_DIRECTORY", "")
	setEnv("TERMINAL_EMULATOR", "foo")
	assert.ToBeEqual(t, false, SupportAssumptionDetectionIntellij())

	setEnv("TERMINAL_EMULATOR", "hello, JetBrains, world")
	assert.ToBeEqual(t, true, SupportAssumptionDetectionIntellij())
}

func Test_SupportAssumptionDetectionGitlabCi(t *testing.T) {
	defer prepareEnv("CI_JOB_ID", "CI_RUNNER_ID")()

	assert.ToBeEqual(t, false, SupportAssumptionDetectionGitlabCi())

	setEnv("CI_JOB_ID", "foo")
	setEnv("CI_RUNNER_ID", "")
	assert.ToBeEqual(t, false, SupportAssumptionDetectionGitlabCi())

	setEnv("CI_JOB_ID", "")
	setEnv("CI_RUNNER_ID", "foo")
	assert.ToBeEqual(t, false, SupportAssumptionDetectionGitlabCi())

	setEnv("CI_JOB_ID", "foo")
	setEnv("CI_RUNNER_ID", "bar")
	assert.ToBeEqual(t, true, SupportAssumptionDetectionGitlabCi())
}
func Test_SupportAssumptionDetectionGithubActions(t *testing.T) {
	defer prepareEnv("GITHUB_RUN_ID", "GITHUB_WORKFLOW")()

	assert.ToBeEqual(t, false, SupportAssumptionDetectionGithubActions())

	setEnv("GITHUB_RUN_ID", "foo")
	setEnv("GITHUB_WORKFLOW", "")
	assert.ToBeEqual(t, false, SupportAssumptionDetectionGithubActions())

	setEnv("GITHUB_RUN_ID", "")
	setEnv("GITHUB_WORKFLOW", "foo")
	assert.ToBeEqual(t, false, SupportAssumptionDetectionGithubActions())

	setEnv("GITHUB_RUN_ID", "foo")
	setEnv("GITHUB_WORKFLOW", "bar")
	assert.ToBeEqual(t, true, SupportAssumptionDetectionGithubActions())
}

func setEnv(key, value string) {
	if value != "" {
		_ = os.Setenv(key, value)
	} else {
		_ = os.Unsetenv(key)
	}
}

func prepareEnv(keys ...string) func() {
	oldValues := map[string]*string{}

	for _, key := range keys {
		if oldValue, hasOld := os.LookupEnv(key); hasOld {
			oldValues[key] = &oldValue
		} else {
			oldValues[key] = nil
		}

		_ = os.Unsetenv(key)
	}

	return func() {
		for key, value := range oldValues {
			if value != nil {
				_ = os.Setenv(key, *value)
			} else {
				_ = os.Unsetenv(key)
			}
		}
	}
}
