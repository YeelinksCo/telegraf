//go:build !custom || inputs || inputs.consul

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/consul" // register plugin
