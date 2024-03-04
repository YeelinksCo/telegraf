//go:build !custom || inputs || inputs.udp_listener

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/udp_listener" // register plugin
