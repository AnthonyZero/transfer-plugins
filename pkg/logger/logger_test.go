package logger

import (
	"testing"
)

func TestJSONLogger(t *testing.T) {
	InitZapLogger("C:/develop/logs/transfer", ToLevel("warn"))
	Info("info message")
	Debug("debug message")
	Warnf("%s message", "warn")
	Error("error mssage")
	//Panic("panic message")
}
