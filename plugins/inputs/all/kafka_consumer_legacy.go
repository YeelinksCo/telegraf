//go:build !custom || inputs || inputs.kafka_consumer_legacy

package all

import _ "github.com/YeelinksCo/telegraf/plugins/inputs/kafka_consumer_legacy" // register plugin
