//go:build !custom || (migrations && (inputs || inputs.nats_consumer))

package all

import _ "github.com/YeelinksCo/telegraf/migrations/inputs_nats_consumer" // register migration
