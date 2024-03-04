//go:build !custom || secretstores || secretstores.os

package all

import _ "github.com/YeelinksCo/telegraf/plugins/secretstores/os" // register plugin
