//go:build !custom || inputs || inputs.socket_listener

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/socket_listener" // register plugin
