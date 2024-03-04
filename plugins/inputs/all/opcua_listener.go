//go:build !custom || inputs || inputs.opcua_listener

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/opcua_listener" // register plugin
