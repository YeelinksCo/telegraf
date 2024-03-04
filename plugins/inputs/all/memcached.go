//go:build !custom || inputs || inputs.memcached

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/memcached" // register plugin
