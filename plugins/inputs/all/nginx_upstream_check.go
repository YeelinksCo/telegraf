//go:build !custom || inputs || inputs.nginx_upstream_check

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/nginx_upstream_check" // register plugin
