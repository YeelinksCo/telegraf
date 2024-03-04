//go:generate ../../../tools/readme_config_includer/generator
//go:build !linux

package sensors

import (
	_ "embed"

	"github.com/YeelinksCo/telegraf"
	"github.com/YeelinksCo/telegraf/plugins/inputs"
)

//go:embed sample.conf
var sampleConfig string

type Sensors struct {
	Log telegraf.Logger `toml:"-"`
}

func (s *Sensors) Init() error {
	s.Log.Warn("current platform is not supported")
	return nil
}
func (*Sensors) SampleConfig() string                { return sampleConfig }
func (*Sensors) Gather(_ telegraf.Accumulator) error { return nil }

func init() {
	inputs.Add("sensors", func() telegraf.Input {
		return &Sensors{}
	})
}
