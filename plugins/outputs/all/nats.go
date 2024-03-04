//go:build !custom || outputs || outputs.nats

package all

import _ "github.com/YeelinksCo/telegraf/plugins/outputs/nats" // register plugin
