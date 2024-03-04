//go:build !custom || inputs || inputs.nats

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/nats" // register plugin
