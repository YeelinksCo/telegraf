//go:build !custom || inputs || inputs.nomad

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/nomad" // register plugin
