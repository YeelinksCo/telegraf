//go:build !custom || inputs || inputs.knx_listener

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/knx_listener" // register plugin
