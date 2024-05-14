package api

import (
	"HCManger/setup"
	"encoding/json"
	"log"
	"net/http"
)

func IncreaseInfer(w http.ResponseWriter, r *http.Request) {
	log.Println("작업할당 메세지", r.RemoteAddr)
	var message requestData
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Message is invalid", http.StatusBadRequest)
		return
	}

	nodeID := setup.ConvertToNodeID(r.RemoteAddr, message.Port)
	modelkey := message.Id + "@" + message.Model
	version := message.Version

	if _, exists := setup.NodeMap[nodeID].TaskInfo[modelkey]; !exists {
		setup.NodeMap[nodeID].TaskInfo[modelkey] = make(map[string]setup.TaskInfo)
		setup.NodeMap[nodeID].TaskInfo[modelkey][version] = setup.TaskInfo{}
	}

	prevTaskinfo := setup.NodeMap[nodeID].TaskInfo[modelkey][version]

	setup.NodeMap[nodeID].TaskInfo[modelkey][version] = setup.TaskInfo{
		AverageComputation: prevTaskinfo.AverageComputation,
		Completion:         prevTaskinfo.Completion,
		AverageCompletion:  prevTaskinfo.AverageCompletion,
		LoadedAmount:       prevTaskinfo.LoadedAmount + 1,
	}

}
