package setup

import "strings"

// TFLOPS
const (
	RTX4090      = 90
	RTX4080      = 48.74
	RTX4070Ti    = 40.09
	RTX3090Ti    = 40
	RTX3080Ti    = 34.1
	RTX3090      = 35.7
	RTX3080      = 29.7
	RTX3070Ti    = 21.7
	RTX4060Ti    = 22.06
	RTX3070      = 20.3
	RTX2080Ti    = 13.45
	RTX2080      = 11.14
	RTX3060Ti    = 16.2
	RTX3060      = 12.7
	RTX2080SUPER = 11.15
	GTX1080Ti    = 11.3
)

var gpuTFLOPSMap = map[string]float32{
	"4090":       RTX4090,
	"4080":       RTX4080,
	"4070 Ti":    RTX4070Ti,
	"4060 Ti":    RTX4060Ti,
	"3090 Ti":    RTX3090Ti,
	"3090":       RTX3090,
	"3080 Ti":    RTX3080Ti,
	"3080":       RTX3080,
	"3070 Ti":    RTX3070Ti,
	"3070":       RTX3070,
	"3060 Ti":    RTX3060Ti,
	"3060":       RTX3060,
	"2080 Ti":    RTX2080Ti,
	"2080 SUPER": RTX2080SUPER,
	"2080":       RTX2080,
	"1080 Ti":    GTX1080Ti,
	"N/A":        10,
}

func ConvertToFlops(gpuName string) float32 {
	for key, value := range gpuTFLOPSMap {
		if strings.Contains(gpuName, key) {
			return value
		}
	}
	return gpuTFLOPSMap["N/A"]
}
