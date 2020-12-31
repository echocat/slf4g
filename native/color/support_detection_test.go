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
	SupportAssumptionDetections = []SupportAssumptionDetection{func() (bool, error) {
		return givenResult, nil
	}}

	givenResult = false
	actual1, err1 := CanSupportBeAssumed()
	assert.ToBeNil(t, err1)
	assert.ToBeEqual(t, givenResult, actual1)

	givenResult = true
	actual2, err2 := CanSupportBeAssumed()
	assert.ToBeNil(t, err2)
	assert.ToBeEqual(t, givenResult, actual2)
}

func Test_SupportAssumptionDetectionIntellij(t *testing.T) {
	defer prepareEnv("IDEA_INITIAL_DIRECTORY", "TERMINAL_EMULATOR")()

	actual1, err1 := SupportAssumptionDetectionIntellij()
	assert.ToBeNil(t, err1)
	assert.ToBeEqual(t, false, actual1)

	setEnv("IDEA_INITIAL_DIRECTORY", "foo")
	actual2, err2 := SupportAssumptionDetectionIntellij()
	assert.ToBeNil(t, err2)
	assert.ToBeEqual(t, true, actual2)

	setEnv("IDEA_INITIAL_DIRECTORY", "")
	setEnv("TERMINAL_EMULATOR", "foo")
	actual3, err3 := SupportAssumptionDetectionIntellij()
	assert.ToBeNil(t, err3)
	assert.ToBeEqual(t, false, actual3)

	setEnv("TERMINAL_EMULATOR", "hello, JetBrains, world")
	actual4, err4 := SupportAssumptionDetectionIntellij()
	assert.ToBeNil(t, err4)
	assert.ToBeEqual(t, true, actual4)
}

func Test_SupportAssumptionDetectionGitlabCi(t *testing.T) {
	defer prepareEnv("CI_JOB_ID", "CI_RUNNER_ID")()

	actual1, err1 := SupportAssumptionDetectionGitlabCi()
	assert.ToBeNil(t, err1)
	assert.ToBeEqual(t, false, actual1)

	setEnv("CI_JOB_ID", "foo")
	setEnv("CI_RUNNER_ID", "")
	actual2, err2 := SupportAssumptionDetectionGitlabCi()
	assert.ToBeNil(t, err2)
	assert.ToBeEqual(t, false, actual2)

	setEnv("CI_JOB_ID", "")
	setEnv("CI_RUNNER_ID", "foo")
	actual3, err3 := SupportAssumptionDetectionGitlabCi()
	assert.ToBeNil(t, err3)
	assert.ToBeEqual(t, false, actual3)

	setEnv("CI_JOB_ID", "foo")
	setEnv("CI_RUNNER_ID", "bar")
	actual4, err4 := SupportAssumptionDetectionGitlabCi()
	assert.ToBeNil(t, err4)
	assert.ToBeEqual(t, true, actual4)
}
func Test_SupportAssumptionDetectionGithubActions(t *testing.T) {
	defer prepareEnv("GITHUB_RUN_ID", "GITHUB_WORKFLOW")()

	actual1, err1 := SupportAssumptionDetectionGithubActions()
	assert.ToBeNil(t, err1)
	assert.ToBeEqual(t, false, actual1)

	setEnv("GITHUB_RUN_ID", "foo")
	setEnv("GITHUB_WORKFLOW", "")
	actual2, err2 := SupportAssumptionDetectionGithubActions()
	assert.ToBeNil(t, err2)
	assert.ToBeEqual(t, false, actual2)

	setEnv("GITHUB_RUN_ID", "")
	setEnv("GITHUB_WORKFLOW", "foo")
	actual3, err3 := SupportAssumptionDetectionGithubActions()
	assert.ToBeNil(t, err3)
	assert.ToBeEqual(t, false, actual3)

	setEnv("GITHUB_RUN_ID", "foo")
	setEnv("GITHUB_WORKFLOW", "bar")
	actual4, err4 := SupportAssumptionDetectionGithubActions()
	assert.ToBeNil(t, err4)
	assert.ToBeEqual(t, true, actual4)
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
