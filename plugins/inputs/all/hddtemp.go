//go:build !custom || inputs || inputs.hddtemp

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/hddtemp" // register plugin
