//go:build !custom || outputs || outputs.sensu

package all

import _ "github.com/YeelinksCo/telegraf/plugins/outputs/sensu" // register plugin
