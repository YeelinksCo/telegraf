//go:build !custom || inputs || inputs.cloudwatch

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/cloudwatch" // register plugin
