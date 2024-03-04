package parallel

import "github.com/YeelinksCo/telegraf"

type Parallel interface {
	Enqueue(telegraf.Metric)
	Stop()
}
