//go:build !custom || processors || processors.dedup

package all

import _ "github.com/YeelinksCo/telegraf/plugins/processors/dedup" // register plugin
