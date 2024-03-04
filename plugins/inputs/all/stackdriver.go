//go:build !custom || inputs || inputs.stackdriver

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/stackdriver" // register plugin
