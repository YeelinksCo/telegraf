//go:build !custom || outputs || outputs.loki

package all

import _ "github.com/YeelinksCo/telegraf/plugins/outputs/loki" // register plugin
