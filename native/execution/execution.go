// Package execution contains tooling methods for simple executions.
package execution

// Execution is a simple interface for a simple execution.
// It is used to wrap executions that can throw an error. It is most likely
// used to simplify the error handling and make the methods better readable.
// To achieve that please use Execute(..).
type Execution func() error

// Execute will execute the given instances of Execution.
func Execute(executions ...Execution) error {
	for _, execution := range executions {
		if execution != nil {
			if err := execution(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Join will combine the given instances of Execution to one.
func Join(executions ...Execution) Execution {
	return func() error {
		return Execute(executions...)
	}
}
