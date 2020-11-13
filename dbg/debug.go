package dbg

import (
	"encoding/json"
	"log"
)

var Debug bool

func SetDebug(val bool) {
	Debug = val
}

func ConsoleLog(obj ...interface{}) {
	if Debug == true {
		a, _ := json.MarshalIndent(obj, "", "  ")
		log.Println(string(a))
	}
}

// MonitorFunc helps with debugging functions
func MonitorFunc(funcName string) func() {
	ConsoleLog("entered " + funcName)
	return func() {
		ConsoleLog("exited " + funcName)
	}
}

func FatalDebug(obj ...interface{}) {
	if Debug == true {
		a, _ := json.MarshalIndent(obj, "", "  ")
		log.Fatal(string(a))
	}
}
