package color

import (
	"os"
	"strings"
)

type SupportAssumptionDetection func() bool

var (
	SupportAssumptionDetections = []SupportAssumptionDetection{
		SupportAssumptionDetectionIntellij,
		SupportAssumptionDetectionGitlabCi,
	}

	SupportAssumptionDetectionIntellij = func() bool {
		// https://stackoverflow.com/questions/61920425/intellij-terminal-environment-variable-set-global
		return os.Getenv("IDEA_INITIAL_DIRECTORY") != "" ||
			strings.Contains(os.Getenv("TERMINAL_EMULATOR"), "JetBrains")
	}

	SupportAssumptionDetectionGitlabCi = func() bool {
		// https://docs.gitlab.com/ee/ci/variables/#list-all-environment-variables
		return os.Getenv("CI_JOB_ID") != "" && os.Getenv("CI_RUNNER_ID") != ""
	}
)
