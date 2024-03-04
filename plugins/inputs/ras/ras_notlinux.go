//go:generate ../../../tools/readme_config_includer/generator
//go:build !linux || (linux && !386 && !amd64 && !arm && !arm64)

package ras

import (
	_ "embed"

	"github.com/YeelinksCo/telegraf"
	"github.com/YeelinksCo/telegraf/plugins/inputs"
)

//go:embed sample.conf
var sampleConfig string

type Ras struct {
	Log telegraf.Logger `toml:"-"`
}

func (r *Ras) Init() error {
	r.Log.Warn("current platform is not supported")
	return nil
}
func (*Ras) SampleConfig() string                { return sampleConfig }
func (*Ras) Gather(_ telegraf.Accumulator) error { return nil }
func (*Ras) Start(_ telegraf.Accumulator) error  { return nil }
func (*Ras) Stop()                               {}

func init() {
	inputs.Add("ras", func() telegraf.Input {
		return &Ras{}
	})
}
