//go:build !custom || outputs || outputs.opentelemetry

package all

import _ "github.com/YeelinksCo/telegraf/plugins/outputs/opentelemetry" // register plugin
