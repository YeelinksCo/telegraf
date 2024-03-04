//go:build !custom || inputs || inputs.dns_query

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/dns_query" // register plugin
