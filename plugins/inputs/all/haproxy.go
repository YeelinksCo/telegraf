//go:build !custom || inputs || inputs.haproxy

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/haproxy" // register plugin
