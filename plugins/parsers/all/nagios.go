//go:build !custom || parsers || parsers.nagios

package all

import _ "github.com/YeelinksCo/telegraf/plugins/parsers/nagios" // register plugin
