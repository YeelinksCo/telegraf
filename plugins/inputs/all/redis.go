//go:build !custom || inputs || inputs.redis

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/redis" // register plugin
