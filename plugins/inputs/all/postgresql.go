//go:build !custom || inputs || inputs.postgresql

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/postgresql" // register plugin
