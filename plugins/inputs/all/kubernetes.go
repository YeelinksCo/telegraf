//go:build !custom || inputs || inputs.kubernetes

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/kubernetes" // register plugin
