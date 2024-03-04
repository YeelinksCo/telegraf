//go:build !custom || serializers || serializers.graphite

package all

import (
	_ "github.com/YeelinksCo/telegraf/plugins/serializers/graphite" // register plugin
)
