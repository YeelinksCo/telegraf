//go:build !custom || inputs || inputs.tcp_listener

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/tcp_listener" // register plugin
