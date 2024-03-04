package mqtt_consumer

import (
	"github.com/YeelinksCo/telegraf"
)

type mqttLogger struct {
	telegraf.Logger
}

func (l mqttLogger) Printf(fmt string, args ...interface{}) {
	l.Logger.Debugf(fmt, args...)
}

func (l mqttLogger) Println(args ...interface{}) {
	l.Logger.Debug(args...)
}
