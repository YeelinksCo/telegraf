//go:build !custom || inputs || inputs.sysstat

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/sysstat" // register plugin
