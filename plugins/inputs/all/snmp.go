//go:build !custom || inputs || inputs.snmp

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/snmp" // register plugin
