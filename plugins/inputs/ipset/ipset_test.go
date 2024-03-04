package ipset

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/YeelinksCo/telegraf/config"
	"github.com/YeelinksCo/telegraf/testutil"
)

func TestIpset(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		tags   []map[string]string
		fields [][]map[string]interface{}
		err    error
	}{
		{
			name:  "0 sets, no results",
			value: "",
		},
		{
			name: "Empty sets, no values",
			value: `create myset hash:net family inet hashsize 1024 maxelem 65536
				create myset2 hash:net,port family inet hashsize 16384 maxelem 524288 counters comment
				`,
		},
		{
			name: "Non-empty sets, but no counters, no results",
			value: `create myset hash:net family inet hashsize 1024 maxelem 65536
				add myset 1.2.3.4
				`,
		},
		{
			name: "Line with data but not enough fields",
			value: `create hash:net family inet hashsize 1024 maxelem 65536 counters
				add myset 4.5.6.7 packets 123 bytes
				`,
			err: errors.New("error parsing line (expected at least 7 fields): \t\t\t\tadd myset 4.5.6.7 packets 123 bytes"),
		},
		{
			name: "Non-empty sets, counters, no comment",
			value: `create myset hash:net family inet hashsize 1024 maxelem 65536 counters
				add myset 1.2.3.4 packets 1328 bytes 79680
				add myset 2.3.4.5 packets 0 bytes 0
				add myset 3.4.5.6 packets 3 bytes 222
				`,
			tags: []map[string]string{
				{"set": "myset", "rule": "1.2.3.4"},
				{"set": "myset", "rule": "3.4.5.6"},
			},
			fields: [][]map[string]interface{}{
				{map[string]interface{}{"packets_total": uint64(1328), "bytes_total": uint64(79680)}},
				{map[string]interface{}{"packets_total": uint64(3), "bytes_total": uint64(222)}},
			},
		},
		{
			name: "Sets with counters and comment",
			value: `create myset hash:net family inet hashsize 1024 maxelem 65536 counters comment
				add myset 1.2.3.4 packets 1328 bytes 79680 comment "first IP"
				add myset 2.3.4.5 packets 0 bytes 0 comment "2nd IP"
				add myset 3.4.5.6 packets 3 bytes 222 "3rd IP"
				`,
			tags: []map[string]string{
				{"set": "myset", "rule": "1.2.3.4"},
				{"set": "myset", "rule": "3.4.5.6"},
			},
			fields: [][]map[string]interface{}{
				{map[string]interface{}{"packets_total": uint64(1328), "bytes_total": uint64(79680)}},
				{map[string]interface{}{"packets_total": uint64(3), "bytes_total": uint64(222)}},
			},
		},
		{
			name: "Sets with and without timeouts",
			value: `create counter-test hash:ip family inet hashsize 1024 maxelem 65536 timeout 1800 counters
				add counter-test 192.168.1.1 timeout 1792 packets 8 bytes 672
				create counter-test2 hash:ip family inet hashsize 1024 maxelem 65536 counters
				add counter-test2 192.168.1.1 packets 18 bytes 673
				`,
			tags: []map[string]string{
				{"set": "counter-test", "rule": "192.168.1.1"},
				{"set": "counter-test2", "rule": "192.168.1.1"},
			},
			fields: [][]map[string]interface{}{
				{map[string]interface{}{"packets_total": uint64(8), "bytes_total": uint64(672), "timeout": uint64(1792)}},
				{map[string]interface{}{"packets_total": uint64(18), "bytes_total": uint64(673)}},
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i++
			ips := &Ipset{
				lister: func(config.Duration, bool) (*bytes.Buffer, error) {
					return bytes.NewBufferString(tt.value), nil
				},
			}
			acc := new(testutil.Accumulator)
			err := acc.GatherError(ips.Gather)
			if !reflect.DeepEqual(tt.err, err) {
				t.Errorf("%d: expected error '%#v' got '%#v'", i, tt.err, err)
			}
			if len(tt.tags) == 0 {
				n := acc.NFields()
				if n != 0 {
					t.Errorf("%d: expected 0 values got %d", i, n)
				}
				return
			}
			n := 0
			for j, tags := range tt.tags {
				for k, fields := range tt.fields[j] {
					if len(acc.Metrics) < n+1 {
						t.Errorf("%d: expected at least %d values got %d", i, n+1, len(acc.Metrics))
						break
					}
					m := acc.Metrics[n]
					if !reflect.DeepEqual(m.Measurement, measurement) {
						t.Errorf("%d %d %d: expected measurement '%#v' got '%#v'\n", i, j, k, measurement, m.Measurement)
					}
					if !reflect.DeepEqual(m.Tags, tags) {
						t.Errorf("%d %d %d: expected tags\n%#v got\n%#v\n", i, j, k, tags, m.Tags)
					}
					if !reflect.DeepEqual(m.Fields, fields) {
						t.Errorf("%d %d %d: expected fields\n%#v got\n%#v\n", i, j, k, fields, m.Fields)
					}
					n++
				}
			}
		})
	}
}

func TestIpset_Gather_listerError(t *testing.T) {
	errFoo := errors.New("error foobar")
	ips := &Ipset{
		lister: func(config.Duration, bool) (*bytes.Buffer, error) {
			return new(bytes.Buffer), errFoo
		},
	}
	acc := new(testutil.Accumulator)
	err := acc.GatherError(ips.Gather)
	if !reflect.DeepEqual(err, errFoo) {
		t.Errorf("Expected error %#v got\n%#v\n", errFoo, err)
	}
}
