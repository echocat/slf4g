package color

import (
	"os"
	"strings"
)

// SupportAssumptionDetection is a function that detects for the current
// environment if color support can be assumed. This can for example be done in
// if this application runs in the context of a GitLabCi run, inside an IDE, ...
type SupportAssumptionDetection func() bool

// CanSupportBeAssumed returns true if the support of color can be assumed.
func CanSupportBeAssumed() bool {
	for _, d := range SupportAssumptionDetections {
		if d() {
			return true
		}
	}
	return false
}

// SupportAssumptionDetections holds all global registered
// SupportAssumptionDetection variants that should be used to discover color the
// support.
var SupportAssumptionDetections = []SupportAssumptionDetection{
	SupportAssumptionDetectionIntellij,
	SupportAssumptionDetectionGitlabCi,
	SupportAssumptionDetectionGithubActions,
}

// SupportAssumptionDetectionIntellij returns true if this application is
// executed in the context of the IntelliJ IDEA framework
// (https://jetbrains.com/idea) or derivatives.
//
// See: https://stackoverflow.com/questions/61920425/intellij-terminal-environment-variable-set-global
func SupportAssumptionDetectionIntellij() bool {
	return os.Getenv("IDEA_INITIAL_DIRECTORY") != "" ||
		strings.Contains(os.Getenv("TERMINAL_EMULATOR"), "JetBrains")
}

// SupportAssumptionDetectionGitlabCi returns true if this application is
// executed in the context of a GitLabCI run (https://docs.gitlab.com/ee/ci/).
//
// See: https://docs.gitlab.com/ee/ci/variables/#list-all-environment-variables
func SupportAssumptionDetectionGitlabCi() bool {
	return os.Getenv("CI_JOB_ID") != "" && os.Getenv("CI_RUNNER_ID") != ""
}

// SupportAssumptionDetectionGithubActions returns true if this application is
// executed in the context of a GitHub Actions run
// (https://docs.github.com/en/free-pro-team@latest/actions).
//
// See: https://docs.gitlab.com/ee/ci/variables/#list-all-environment-variables
func SupportAssumptionDetectionGithubActions() bool {
	return os.Getenv("GITHUB_RUN_ID") != "" && os.Getenv("GITHUB_WORKFLOW") != ""
}
