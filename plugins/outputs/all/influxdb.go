//go:build !custom || outputs || outputs.influxdb

package all

import _ "github.com/YeelinksCo/telegraf/plugins/outputs/influxdb" // register plugin
