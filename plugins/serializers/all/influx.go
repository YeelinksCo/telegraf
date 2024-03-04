//go:build !custom || serializers || serializers.influx

package all

import (
	_ "github.com/YeelinksCo/telegraf/plugins/serializers/influx" // register plugin
)
