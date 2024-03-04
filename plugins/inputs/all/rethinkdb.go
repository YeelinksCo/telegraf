//go:build !custom || inputs || inputs.rethinkdb

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/rethinkdb" // register plugin
