//go:build !custom || inputs || inputs.intel_pmu

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/intel_pmu" // register plugin
