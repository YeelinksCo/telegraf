//go:build !custom || migrations

package all

import _ "github.com/YeelinksCo/telegraf/migrations/general_metricfilter" // register migration
