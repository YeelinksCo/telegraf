//go:build !custom || (migrations && (inputs || inputs.disk))

package all

import _ "github.com/YeelinksCo/telegraf/migrations/inputs_disk" // register migration
