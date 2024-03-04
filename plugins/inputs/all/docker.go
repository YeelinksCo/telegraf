//go:build !custom || inputs || inputs.docker

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/docker" // register plugin
