//go:build !custom || inputs || inputs.processes

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/processes" // register plugin
