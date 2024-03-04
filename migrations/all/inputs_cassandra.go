//go:build !custom || (migrations && (inputs || inputs.cassandra))

package all

import _ "github.com/YeelinksCo/telegraf/migrations/inputs_cassandra" // register migration
