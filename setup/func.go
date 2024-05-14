package setup

import "strings"

func ErasePort(ip string) string {
	port := strings.Split(ip, ":")
	return port[0]
}

//임시 ID 생성자
func ConvertToNodeID(address string, port string) string {
	return ErasePort(address) + ":" + port
}
