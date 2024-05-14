package api

import (
	"HCManger/setup"
	"encoding/json"
	"net/http"
)

func HandleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//log.Println("요청 :", r.RemoteAddr)

	switch r.URL.Path {
	case "/nodemap":
		json.NewEncoder(w).Encode(setup.NodeMap)
	case "/modelmap":
		json.NewEncoder(w).Encode(setup.ModelMap)
	default:
		http.NotFound(w, r)
	}
}
