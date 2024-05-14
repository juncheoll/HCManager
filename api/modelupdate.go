package api

import (
	"HCManger/setup"
	"encoding/json"
	"log"
	"net/http"
)

func ModelUpdate(w http.ResponseWriter, r *http.Request) {
	var message requestData
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Message is invalid", http.StatusBadRequest)
		return
	}

	log.Println("노드에게 modelupdate 수신 :", message)
	nodeID := setup.ConvertToNodeID(r.RemoteAddr, message.Port)
	modelkey := message.Id + "@" + message.Model

	if _, exists := setup.NodeMap[nodeID].TaskInfo[modelkey]; !exists {
		setup.NodeMap[nodeID].TaskInfo[modelkey] = make(map[string]setup.TaskInfo)
	}
	setup.NodeMap[nodeID].TaskInfo[modelkey][message.Version] = setup.TaskInfo{}

	if _, exxists := setup.ModelMap[modelkey]; !exxists {
		setup.ModelMap[modelkey] = make(map[string][]string)
	}
	setup.ModelMap[modelkey][message.Version] = append(setup.ModelMap[modelkey][message.Version], nodeID)
}
