//go:build !custom || inputs || inputs.snmp_trap

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/snmp_trap" // register plugin
