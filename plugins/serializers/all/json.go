//go:build !custom || serializers || serializers.json

package all

import (
	_ "github.com/YeelinksCo/telegraf/plugins/serializers/json" // register plugin
)
