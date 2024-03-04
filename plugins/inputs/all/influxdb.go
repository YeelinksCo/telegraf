//go:build !custom || inputs || inputs.influxdb

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/influxdb" // register plugin
