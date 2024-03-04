//go:build !custom || outputs || outputs.nsq

package all

import _ "github.com/YeelinksCo/telegraf/plugins/outputs/nsq" // register plugin
