//go:build !custom || outputs || outputs.kinesis

package all

import _ "github.com/YeelinksCo/telegraf/plugins/outputs/kinesis" // register plugin
