//go:build !custom || inputs || inputs.jolokia2_agent

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/jolokia2_agent" // register plugin
