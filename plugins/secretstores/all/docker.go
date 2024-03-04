//go:build !custom || secretstores || secretstores.docker

package all

import _ "github.com/YeelinksCo/telegraf/plugins/secretstores/docker" // register plugin
