//go:build !linux

package cgroup

import (
	"github.com/YeelinksCo/telegraf"
)

func (*CGroup) Gather(_ telegraf.Accumulator) error {
	return nil
}
