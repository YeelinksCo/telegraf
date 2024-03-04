//go:build !custom || outputs || outputs.postgresql

package all

import _ "github.com/YeelinksCo/telegraf/plugins/outputs/postgresql" // register plugin
