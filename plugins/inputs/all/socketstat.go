//go:build !custom || inputs || inputs.socketstat

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/socketstat" // register plugin
