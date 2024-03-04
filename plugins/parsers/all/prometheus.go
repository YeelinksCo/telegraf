//go:build !custom || parsers || parsers.prometheus

package all

import _ "github.com/YeelinksCo/telegraf/plugins/parsers/prometheus" // register plugin
