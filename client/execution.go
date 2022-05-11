package client

import scheduler "github.com/starkandwayne/scheduler-for-ocf/core"

type byExecutionStart []*scheduler.Execution

func (s byExecutionStart) Len() int {
	return len(s)
}

func (s byExecutionStart) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byExecutionStart) Less(i int, j int) bool {
	return s[i].ExecutionStartTime.Before(s[j].ExecutionStartTime)
}
