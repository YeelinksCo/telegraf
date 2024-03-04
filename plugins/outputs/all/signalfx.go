//go:build !custom || outputs || outputs.signalfx

package all

import _ "github.com/YeelinksCo/telegraf/plugins/outputs/signalfx" // register plugin
