package api

import (
	"HCManger/setup"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"
)

type TaskInfoFromAgent struct {
	LoadedAmount         int     `json:"loaded_amount"`
	AverageInferenceTime float32 `json:"average_inference_time"`
}

type AliveMsg struct {
	Port       string                                  `json:"port"`
	GpuName    string                                  `json:"gpuname"`
	Model_info map[string]map[string]TaskInfoFromAgent `json:"model_info"` //key1 : Provider@Name, Key2 : Version
}

func AliveHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	var message AliveMsg
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		log.Println("join request :", setup.ErasePort(r.RemoteAddr), "invalid message")
		http.Error(w, "Message is invalid", http.StatusBadRequest)
		return
	}

	nodeID := setup.ConvertToNodeID(r.RemoteAddr, message.Port)
	log.Println("Alive :", r.RemoteAddr)

	aliveTime[nodeID] = time.Now()
	//리스트에 이미 있는 노드라면
	if _, exists := setup.NodeMap[nodeID]; exists {
		updateModelInfo(nodeID, message.Model_info)
		/*
			for modelkey, value := range nodeInfo.TaskInfo {
				for version, prevTaskInfo := range value {
					newTaskInfo := setup.TaskInfo{
						AverageComputation: message.Model_info[modelkey][version].AverageInferenceTime,
						Completion:         prevTaskInfo.Completion,
						AverageCompletion:  prevTaskInfo.AverageCompletion,
						LoadedAmount:       prevTaskInfo.LoadedAmount,
					}
					setup.NodeMap[nodeID].TaskInfo[modelkey][version] = newTaskInfo
				}
			}
		*/
	} else { //처음 받는 노드에 대한 정보라면
		log.Println("Join :", nodeID)
		setup.NodeMap[nodeID] = setup.NodeInfo{
			ID:       nodeID,
			Address:  setup.ErasePort(r.RemoteAddr) + ":" + message.Port,
			TFLOPS:   setup.ConvertToFlops(message.GpuName),
			TaskInfo: make(map[string]map[string]setup.TaskInfo),
		}
		updateModelInfo(nodeID, message.Model_info)
		go healthCheck2(nodeID)
	}

}

var aliveTime map[string]time.Time = make(map[string]time.Time)

func healthCheck2(nodeID string) {
	defer delete(setup.NodeMap, nodeID)
	//일정 시간 안들어오면 리스트에서 삭제

	for {
		time.Sleep(time.Second * 60)
		currTime := time.Now()
		prevAliveTime := aliveTime[nodeID]
		duration := currTime.Sub(prevAliveTime)

		if duration >= 5*time.Second {
			log.Println("Exit :", nodeID)

			for modelkey, taskinfoPerModelkey := range setup.NodeMap[nodeID].TaskInfo {
				for version := range taskinfoPerModelkey {
					//modelkey, version이 이 Node가 가지고 있던 Model에 대한 정보
					//setup.ModelMap[modelkey][version]에서 이 nodeID 찾아 제거
					nodeList := setup.ModelMap[modelkey][version]

					for n, m := range nodeList {
						if m == nodeID {
							setup.ModelMap[modelkey][version] = append(setup.ModelMap[modelkey][version][:n], setup.ModelMap[modelkey][version][n+1:]...)
							break
						}
					}
				}
			}

			delete(setup.NodeMap, nodeID)

			return
		}
	}

	//log.Printf("TCP(%s) : %s\n", nodeID, string(buffer))
}

func healthCheck(nodeID string) {
	defer delete(setup.NodeMap, nodeID)
	//일정 시간 안들어오면 리스트에서 삭제
	tcpAddress := setup.ErasePort(setup.NodeMap[nodeID].Address) + ":6934"
	conn, err := net.Dial("tcp", tcpAddress)
	if err != nil {
		log.Println(tcpAddress+"tcp Dial 실패:", err)
		return
	}
	defer conn.Close()
	log.Println(tcpAddress + "tcp Dial 성공")

	log.Println(conn.RemoteAddr().String())

	buffer := make([]byte, 126)
	for {
		_, err = conn.Read(buffer)
		if err != nil {
			log.Println(nodeID+"끊김 : ", err)

			for modelkey, taskinfoPerModelkey := range setup.NodeMap[nodeID].TaskInfo {
				for version := range taskinfoPerModelkey {
					//modelkey, version이 이 Node가 가지고 있던 Model에 대한 정보
					//setup.ModelMap[modelkey][version]에서 이 nodeID 찾아 제거
					nodeList := setup.ModelMap[modelkey][version]

					for n, m := range nodeList {
						if m == nodeID {
							setup.ModelMap[modelkey][version] = append(setup.ModelMap[modelkey][version][:n], setup.ModelMap[modelkey][version][n+1:]...)
							break
						}
					}
				}
			}

			delete(setup.NodeMap, nodeID)

			return
		}

		//log.Printf("TCP(%s) : %s\n", nodeID, string(buffer))
	}
}

// node의 TaskInfo 최신 정보로 갱신
// 기존 TaskInfo를 순회하며, message의 TaskInfoFromAgent의 데이터를 그대로 삽입, TaskInfoFromAgent에서는 삭제
// 순회 중 TaskInfoFromAgent에서는 찾을 수 없는 모델의 정보라면, 모델이 삭제되었다는 뜻일 것 => 삭제
// TaskInfo 순회를 마친 후 남은 TaskInfoFromAgent의 데이터는 아직 갱신이 덜 된 새로운 모델일 것 => 삽입
func updateModelInfo(nodeID string, modelInfo map[string]map[string]TaskInfoFromAgent) {

	//기존 TaskInfo에서 modelInfo에 있는 정보를 덮어쓰기
	for modelkey, value := range setup.NodeMap[nodeID].TaskInfo {
		for version := range value {

			if newtaskPerModelkey, exists := modelInfo[modelkey]; exists {
				if newTask, exists := newtaskPerModelkey[version]; exists {
					//modelkey, version이 modelInfo에 있으면
					updateTaskInfo(nodeID, modelkey, version, newTask)
				} else {
					//modelkey 는 있고, version이 없으면
					delete(setup.NodeMap[nodeID].TaskInfo[modelkey], version)
					if len(setup.NodeMap[nodeID].TaskInfo[modelkey]) == 0 {
						delete(setup.NodeMap[nodeID].TaskInfo, modelkey)
					}

					//modelMap[modelkey][version]에서 nodeID 없애기
					deleteInModelMap1(modelkey, version, nodeID)
				}
				delete(modelInfo[modelkey], version)
				if len(modelInfo[modelkey]) == 0 {
					delete(modelInfo, modelkey)
				}
			} else {
				//modelkey이 modelInfo에 없으면
				delete(setup.NodeMap[nodeID].TaskInfo, modelkey)
				//modelMap에서도 없애주기
				deleteInModelMap2(modelkey, nodeID)
			}
		}
	}

	// TaskInfoFromAgent에 남은 것 삽입
	for modelkey, value := range modelInfo {
		for version, newTask := range value {
			updateTaskInfo(nodeID, modelkey, version, newTask)

			//modelMap[modelkey][version] 에 append
			if _, exists := setup.ModelMap[modelkey]; !exists {
				setup.ModelMap[modelkey] = make(map[string][]string)
			}
			setup.ModelMap[modelkey][version] = append(setup.ModelMap[modelkey][version], nodeID)
		}
	}
}

// nodeID의 modelkey, version에 해당하는 taskInfo 갱신
func updateTaskInfo(nodeID, modelkey, version string, taskFromAgent TaskInfoFromAgent) {
	prevTask := setup.NodeMap[nodeID].TaskInfo

	if taskPerModelkey, exists := prevTask[modelkey]; exists {
		if task, exists := taskPerModelkey[version]; exists {
			setup.NodeMap[nodeID].TaskInfo[modelkey][version] = setup.TaskInfo{
				AverageComputation: task.AverageComputation,
				Completion:         task.Completion,
				AverageCompletion:  task.AverageCompletion,
				LoadedAmount:       float32(taskFromAgent.LoadedAmount),
			}
		} else {
			setup.NodeMap[nodeID].TaskInfo[modelkey][version] = setup.TaskInfo{
				AverageComputation: taskFromAgent.AverageInferenceTime,
				Completion:         taskFromAgent.AverageInferenceTime,
				AverageCompletion:  taskFromAgent.AverageInferenceTime,
				LoadedAmount:       float32(taskFromAgent.LoadedAmount),
			}
		}
	} else {
		setup.NodeMap[nodeID].TaskInfo[modelkey] = make(map[string]setup.TaskInfo)
		setup.NodeMap[nodeID].TaskInfo[modelkey][version] = setup.TaskInfo{
			AverageComputation: taskFromAgent.AverageInferenceTime,
			Completion:         taskFromAgent.AverageInferenceTime,
			AverageCompletion:  taskFromAgent.AverageInferenceTime,
			LoadedAmount:       float32(taskFromAgent.LoadedAmount),
		}
	}
}

func deleteInModelMap1(modelkey, version, nodeID string) {
	setup.ModelMap[modelkey][version] = removeElement(setup.ModelMap[modelkey][version], nodeID)
	if len(setup.ModelMap[modelkey][version]) == 0 {
		delete(setup.ModelMap[modelkey], version)
	}
}

func deleteInModelMap2(modelkey, nodeID string) {
	for version := range setup.ModelMap[modelkey] {
		deleteInModelMap1(modelkey, version, nodeID)
	}
	if len(setup.ModelMap[modelkey]) == 0 {
		delete(setup.ModelMap, modelkey)
	}
}

func removeElement(slice []string, element string) []string {
	for i, v := range slice {
		if v == element {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
