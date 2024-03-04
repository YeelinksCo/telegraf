//go:build !custom || serializers || serializers.msgpack

package all

import (
	_ "github.com/YeelinksCo/telegraf/plugins/serializers/msgpack" // register plugin
)
