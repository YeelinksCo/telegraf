//go:build !custom || (migrations && (inputs || inputs.mqtt_consumer))

package all

import _ "github.com/YeelinksCo/telegraf/migrations/inputs_mqtt_consumer" // register migration
