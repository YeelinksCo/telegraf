//go:build !custom || parsers || parsers.collectd

package all

import _ "github.com/YeelinksCo/telegraf/plugins/parsers/collectd" // register plugin
