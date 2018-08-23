package main

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

func getLevel(level string) int {
	switch level {
	case "debug":
		return logs.LevelDebug
	case "trace":
		return logs.LevelTrace
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "error":
		return logs.LevelError
	default:
		return logs.LevelDebug
	}
}

func initLogs(filename string, level string) (err error) {
	config := make(map[string]interface{})
	config["filename"] = filename
	config["level"] = getLevel(level)

	configStr, err := json.Marshal(config)
	if err != nil {
		return
	}

	logs.SetLogger(logs.AdapterFile, string(configStr))
	// logs.SetLogFuncCall(true) // print file name and row number
	return
}
