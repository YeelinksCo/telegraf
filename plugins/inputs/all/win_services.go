//go:build !custom || inputs || inputs.win_services

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/win_services" // register plugin
