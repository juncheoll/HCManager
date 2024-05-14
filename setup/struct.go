package setup

type NodeInfo struct {
	ID       string                         `json:"id"`
	Address  string                         `json:"address"`
	TFLOPS   float32                        `json:"tflops"`
	TaskInfo map[string]map[string]TaskInfo `json:"taskinfo"`
}

type TaskInfo struct {
	AverageComputation float32 `json:"averagecomputation"`
	Completion         float32 `json:"completion"`
	AverageCompletion  float32 `json:"averagecompletion"`
	LoadedAmount       float32 `json:"loadedamount"`
}

type TaskInfoFromGateway struct {
}
