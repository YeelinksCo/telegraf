//go:build !custom || outputs || outputs.cloudwatch

package all

import _ "github.com/YeelinksCo/telegraf/plugins/outputs/cloudwatch" // register plugin
