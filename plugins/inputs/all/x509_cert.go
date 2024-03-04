//go:build !custom || inputs || inputs.x509_cert

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/x509_cert" // register plugin
