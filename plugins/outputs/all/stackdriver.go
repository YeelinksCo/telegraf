//go:build !custom || outputs || outputs.stackdriver

package all

import _ "github.com/YeelinksCo/telegraf/plugins/outputs/stackdriver" // register plugin
