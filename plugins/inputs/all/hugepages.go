//go:build !custom || inputs || inputs.hugepages

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/hugepages" // register plugin
