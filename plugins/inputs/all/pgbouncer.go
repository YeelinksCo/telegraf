//go:build !custom || inputs || inputs.pgbouncer

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/pgbouncer" // register plugin
