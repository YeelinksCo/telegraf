//go:build !custom || serializers || serializers.template

package all

import (
	_ "github.com/YeelinksCo/telegraf/plugins/serializers/template" // register plugin
)
