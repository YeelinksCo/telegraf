//go:build !custom || inputs || inputs.kinesis_consumer

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/kinesis_consumer" // register plugin
