//go:build !custom || outputs || outputs.opentsdb

package all

import _ "github.com/YeelinksCo/telegraf/plugins/outputs/opentsdb" // register plugin
