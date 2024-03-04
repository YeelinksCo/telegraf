//go:build !custom || outputs || outputs.socket_writer

package all

import _ "github.com/YeelinksCo/telegraf/plugins/outputs/socket_writer" // register plugin
