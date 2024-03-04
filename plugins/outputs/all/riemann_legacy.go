//go:build !custom || outputs || outputs.riemann_legacy

package all

import _ "github.com/YeelinksCo/telegraf/plugins/outputs/riemann_legacy" // register plugin
