package logger

import (
	"github.com/ipfs/go-log/v2"
)

var Logger *log.ZapEventLogger

func Init() {
	Logger = log.Logger("aa")

	log.SetAllLoggers(log.LevelWarn)
	log.SetLogLevel("aa", "info")
}
