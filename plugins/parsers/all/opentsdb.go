//go:build !custom || parsers || parsers.opentsdb

package all

import _ "github.com/YeelinksCo/telegraf/plugins/parsers/opentsdb" // register plugin
