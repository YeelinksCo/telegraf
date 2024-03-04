//go:build !custom || inputs || inputs.iptables

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/iptables" // register plugin
