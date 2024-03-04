//go:build !custom || inputs || inputs.consul_agent

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/consul_agent" // register plugin
