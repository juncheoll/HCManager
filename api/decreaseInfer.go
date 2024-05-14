package api

import (
	"HCManger/setup"
	"encoding/json"
	"log"
	"net/http"
)

type requestData struct {
	Port      string  `json:"port"`
	Id        string  `json:"id"`
	Model     string  `json:"model"`
	Version   string  `json:"version"`
	BurstTime float64 `json:"burstTime"`
}

func DecreaseInfer(w http.ResponseWriter, r *http.Request) {
	log.Println("작업종료 메세지", r.RemoteAddr)
	var message requestData
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Message is invalid", http.StatusBadRequest)
		return
	}

	nodeID := setup.ConvertToNodeID(r.RemoteAddr, message.Port)
	modelkey := message.Id + "@" + message.Model
	Version := message.Version
	prevTaskinfo := setup.NodeMap[nodeID].TaskInfo[modelkey][Version]

	setup.NodeMap[nodeID].TaskInfo[modelkey][Version] = setup.TaskInfo{
		AverageComputation: (prevTaskinfo.AverageComputation + float32(message.BurstTime)) / 2,
		Completion:         prevTaskinfo.Completion,
		AverageCompletion:  prevTaskinfo.AverageCompletion,
		LoadedAmount:       prevTaskinfo.LoadedAmount - 1,
	}

}
