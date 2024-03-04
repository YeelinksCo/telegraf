//go:build !custom || (migrations && (inputs || inputs.jolokia))

package all

import _ "github.com/YeelinksCo/telegraf/migrations/inputs_jolokia" // register migration
