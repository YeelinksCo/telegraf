//go:build !custom || inputs || inputs.redis_sentinel

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/redis_sentinel" // register plugin
