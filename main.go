package main

import (
	"HCManger/api"
	"HCManger/logger"
	"HCManger/setup"
	"math"
	"math/rand"
	"time"

	"fmt"
	"log"
	"net/http"
)

func main() {
	logger.Init()

	//test()

	http.HandleFunc("/alive", api.AliveHandler)
	http.HandleFunc("/nodemap", api.HandleGet)
	http.HandleFunc("/modelmap", api.HandleGet)

	http.HandleFunc("/inference/start", api.IncreaseInfer)
	http.HandleFunc("/inference/end", api.DecreaseInfer)

	http.HandleFunc("/model/update", api.ModelUpdate)
	http.HandleFunc("/model/delete", api.Modeldelete)
	/*
		http.HandleFunc("/join", api.GPUJoinHandler)

		http.HandleFunc("/taskinfo/agent", api.TaskInfoHandler)
		http.HandleFunc("/taskinfo/gateway", api.TaskInfoHandler)

		http.HandleFunc("/modelinfo/upload", api.ModelInfoUploadHandler)
		http.HandleFunc("/modelinfo/delete", api.ModelInfoDeleteHandler)
	*/

	//go monitoring.Monitoring()

	log.Printf("Starting server at %s\n", setup.HttpAddress)

	server := &http.Server{
		Addr:        setup.HttpAddress,
		IdleTimeout: 120 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func test() {
	http.HandleFunc("/test", testHandler)

	//ip, port, gpuname, modelinfo(meta@Llama)
	createTestData()

}

func createTestData() {
	count := 1000
	cnt := 0

	modelkey := "meta@Llama-2-7B-Chat"
	version := "1"

	setup.ModelMap[modelkey] = make(map[string][]string)

	for p1 := 1; p1 <= 255; p1++ {
		for p2 := 1; p2 <= 255; p2++ {
			for p3 := 1; p3 <= 255; p3++ {
				for p4 := 1; p4 <= 255; p4++ {
					if count == cnt {
						return
					}
					ip := fmt.Sprintf("%d.%d.%d.%d", p1, p2, p3, p4)
					address := ip + ":8080"
					TFLOPS := float32(40)
					computation := rand.Float32() * 10
					taskInfo := setup.TaskInfo{
						LoadedAmount:       float32(rand.Intn(100)),
						AverageComputation: roundTo4DecimalPlaces(computation), // = inference Time
						Completion:         roundTo4DecimalPlaces(computation + rand.Float32()*2),
						AverageCompletion:  roundTo4DecimalPlaces(computation + rand.Float32()*2),
					}

					setup.NodeMap[ip] = setup.NodeInfo{
						ID:       ip,
						Address:  address,
						TFLOPS:   TFLOPS,
						TaskInfo: make(map[string]map[string]setup.TaskInfo),
					}

					setup.NodeMap[ip].TaskInfo[modelkey] = make(map[string]setup.TaskInfo)
					setup.NodeMap[ip].TaskInfo[modelkey][version] = taskInfo

					setup.ModelMap[modelkey][version] = append(setup.ModelMap[modelkey][version], ip)

					cnt++
				}
			}
		}
	}
}

func roundTo4DecimalPlaces(f float32) float32 {
	return float32(math.Round(float64(f*10000)) / 10000)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
}
