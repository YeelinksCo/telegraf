//go:build !custom || serializers || serializers.cloudevents

package all

import (
	_ "github.com/YeelinksCo/telegraf/plugins/serializers/cloudevents" // register plugin
)
