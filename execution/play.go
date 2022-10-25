package execution

type Play struct {
	Name string `json:"name"`
	ID string `json:"id"`
	StartTime float32 `json:"startTime"`
	EndTime float32 `json:"endTime"`
	Tasks []string `json:"tasks,omitempty"`
}