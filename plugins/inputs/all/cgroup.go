//go:build !custom || inputs || inputs.cgroup

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/cgroup" // register plugin
