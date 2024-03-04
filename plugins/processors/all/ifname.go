//go:build !custom || processors || processors.ifname

package all

import _ "github.com/YeelinksCo/telegraf/plugins/processors/ifname" // register plugin
