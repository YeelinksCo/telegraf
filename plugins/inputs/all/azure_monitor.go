//go:build !custom || inputs || inputs.azure_monitor

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/azure_monitor" // register plugin
