//go:build !custom || serializers || serializers.csv

package all

import (
	_ "github.com/YeelinksCo/telegraf/plugins/serializers/csv" // register plugin
)
