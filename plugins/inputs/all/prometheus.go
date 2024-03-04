//go:build !custom || inputs || inputs.prometheus

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/prometheus" // register plugin
