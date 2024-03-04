//go:build !custom || inputs || inputs.conntrack

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/conntrack" // register plugin
