package setup

// key:ip, value: NodeInfo
var NodeMap map[string]NodeInfo = make(map[string]NodeInfo)

// key1:ID@NAME, key2:Version, value:nodelist
var ModelMap map[string]map[string][]string = make(map[string]map[string][]string)

var HttpAddress string = "0.0.0.0:80"
