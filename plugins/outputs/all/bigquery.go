//go:build !custom || outputs || outputs.bigquery

package all

import _ "github.com/YeelinksCo/telegraf/plugins/outputs/bigquery" // register plugin
