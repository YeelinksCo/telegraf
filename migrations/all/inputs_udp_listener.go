//go:build !custom || (migrations && (inputs || inputs.udp_listener))

package all

import _ "github.com/YeelinksCo/telegraf/migrations/inputs_udp_listener" // register migration
