//go:build !custom || inputs || inputs.procstat

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/procstat" // register plugin
