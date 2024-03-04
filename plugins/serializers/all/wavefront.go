//go:build !custom || serializers || serializers.wavefront

package all

import (
	_ "github.com/YeelinksCo/telegraf/plugins/serializers/wavefront" // register plugin
)
