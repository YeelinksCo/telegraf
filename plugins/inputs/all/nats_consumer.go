//go:build !custom || inputs || inputs.nats_consumer

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/nats_consumer" // register plugin
