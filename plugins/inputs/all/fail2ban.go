//go:build !custom || inputs || inputs.fail2ban

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/fail2ban" // register plugin
