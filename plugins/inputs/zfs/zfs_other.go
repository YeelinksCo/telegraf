//go:build !linux && !freebsd

package zfs

import (
	"github.com/YeelinksCo/telegraf"
	"github.com/YeelinksCo/telegraf/plugins/inputs"
)

func (*Zfs) Gather(_ telegraf.Accumulator) error {
	return nil
}

func init() {
	inputs.Add("zfs", func() telegraf.Input {
		return &Zfs{}
	})
}
