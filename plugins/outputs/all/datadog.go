//go:build !custom || outputs || outputs.datadog

package all

import _ "github.com/YeelinksCo/telegraf/plugins/outputs/datadog" // register plugin
