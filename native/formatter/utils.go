package formatter

type checkedExecution func() error

func executeChecked(executions ...checkedExecution) error {
	for _, execution := range executions {
		if err := execution(); err != nil {
			return err
		}
	}
	return nil
}