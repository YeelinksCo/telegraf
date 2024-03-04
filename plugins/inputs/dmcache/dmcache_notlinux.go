//go:build !linux

package dmcache

import (
	"github.com/YeelinksCo/telegraf"
)

func (*DMCache) Gather(_ telegraf.Accumulator) error {
	return nil
}

func dmSetupStatus() ([]string, error) {
	return []string{}, nil
}
