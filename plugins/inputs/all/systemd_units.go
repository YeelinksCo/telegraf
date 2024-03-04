//go:build !custom || inputs || inputs.systemd_units

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/systemd_units" // register plugin
