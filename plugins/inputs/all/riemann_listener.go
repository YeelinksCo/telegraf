//go:build !custom || inputs || inputs.riemann_listener

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/riemann_listener" // register plugin
