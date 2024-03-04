//go:build !custom || outputs || outputs.cloudwatch_logs

package all

import _ "github.com/YeelinksCo/telegraf/plugins/outputs/cloudwatch_logs" // register plugin
