//go:build !custom || inputs || inputs.ipmi_sensor

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/ipmi_sensor" // register plugin
