package formatter

type checkedExecution func() error

func executeChecked(executions ...checkedExecution) error {
	for _, execution := range executions {
		if execution != nil {
			if err := execution(); err != nil {
				return err
			}
		}
	}
	return nil
}

func joinCheckedExecutions(executions ...checkedExecution) checkedExecution {
	return func() error {
		return executeChecked(executions...)
	}
}
