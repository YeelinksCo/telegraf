//go:build !custom || secretstores || secretstores.systemd

package all

import _ "github.com/YeelinksCo/telegraf/plugins/secretstores/systemd" // register plugin
