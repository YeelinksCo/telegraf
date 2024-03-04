//go:build !custom || inputs || inputs.statsd

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/statsd" // register plugin
