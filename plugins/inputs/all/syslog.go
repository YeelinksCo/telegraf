//go:build !custom || inputs || inputs.syslog

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/syslog" // register plugin
