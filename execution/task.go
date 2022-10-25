package execution

type Task struct {
	Name string `json:"name"`
	ID string `json:"id"`
	StartTime float32 `json:"startTime"`
	EndTime float32 `json:"endTime"`
	Targets []TargetResults `json:"targets"`
}

type TargetResults struct {
	Machine	string `json:"machine,omitempty"`
	Action           string                 `json:"action,omitempty"`
	Changed          bool                   `json:"changed,omitempty"`
	Stdout           []string               `json:"stdout,omitempty"`
	Stderr           []string               `json:"stderr,omitempty"`
	Failed           bool                   `json:"failed,omitempty"`
	FailedWhenResult bool                   `json:"failedWhenResult,omitempty"`
	Skipped          bool                   `json:"skipped,omitempty"`
	SkipReason       string                 `json:"skipReason,omitempty"`
	Unreachable      bool                   `json:"unreachable,omitempty"`
}