//go:build !custom || inputs || inputs.openntpd

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/openntpd" // register plugin
