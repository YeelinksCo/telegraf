//go:build !custom || inputs || inputs.clickhouse

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/clickhouse" // register plugin
