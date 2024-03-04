//go:build !custom || serializers || serializers.prometheus

package all

import (
	_ "github.com/YeelinksCo/telegraf/plugins/serializers/prometheus" // register plugin
)
