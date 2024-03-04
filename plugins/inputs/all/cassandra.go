//go:build !custom || inputs || inputs.cassandra

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/cassandra" // register plugin
