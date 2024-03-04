//go:build !custom || serializers || serializers.splunkmetric

package all

import (
	_ "github.com/YeelinksCo/telegraf/plugins/serializers/splunkmetric" // register plugin
)
