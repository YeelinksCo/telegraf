//go:build !custom || inputs || inputs.kafka_consumer

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/kafka_consumer" // register plugin
