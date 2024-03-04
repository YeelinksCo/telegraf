//go:build !custom || serializers || serializers.prometheusremotewrite

package all

import (
	_ "github.com/YeelinksCo/telegraf/plugins/serializers/prometheusremotewrite" // register plugin
)
