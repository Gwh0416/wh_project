package kk

import (
	"encoding/json"
	"time"

	"gwh.com/project-common/tms"
)

type FieldMap map[string]any

type KafkaLog struct {
	Type     string
	Action   string
	Time     string
	Msg      string
	Field    FieldMap
	FuncName string
}

func Error(err error, funcName string, fieldMap FieldMap) []byte {
	kl := KafkaLog{
		Type:     "error",
		Action:   "click",
		Time:     tms.Format(time.Now()),
		Msg:      err.Error(),
		Field:    fieldMap,
		FuncName: funcName,
	}
	bytes, _ := json.Marshal(kl)
	return bytes
}

func Info(msg string, funcName string, fieldMap FieldMap) []byte {
	kl := KafkaLog{
		Type:     "info",
		Action:   "click",
		Time:     tms.Format(time.Now()),
		Msg:      msg,
		Field:    fieldMap,
		FuncName: funcName,
	}
	bytes, _ := json.Marshal(kl)
	return bytes
}
