//go:build !custom || outputs || outputs.newrelic

package all

import _ "github.com/YeelinksCo/telegraf/plugins/outputs/newrelic" // register plugin
