//go:build !custom || serializers || serializers.nowmetric

package all

import (
	_ "github.com/YeelinksCo/telegraf/plugins/serializers/nowmetric" // register plugin
)
