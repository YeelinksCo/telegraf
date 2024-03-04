//go:build !custom || inputs || inputs.beanstalkd

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/beanstalkd" // register plugin
