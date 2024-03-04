package cisco_telemetry_mdt

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
	"testing"
	"time"

	dialout "github.com/cisco-ie/nx-telemetry-proto/mdt_dialout"
	telemetryBis "github.com/cisco-ie/nx-telemetry-proto/telemetry_bis"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	"github.com/YeelinksCo/telegraf"
	"github.com/YeelinksCo/telegraf/metric"
	"github.com/YeelinksCo/telegraf/testutil"
)

func TestHandleTelemetryTwoSimple(t *testing.T) {
	c := &CiscoTelemetryMDT{
		Log:       testutil.Logger{},
		Transport: "dummy",
		Aliases: map[string]string{
			"alias": "type:model/some/path",
		},
	}
	acc := &testutil.Accumulator{}
	err := c.Start(acc)
	// error is expected since we are passing in dummy transport
	require.Error(t, err)

	telemetry := &telemetryBis.Telemetry{
		MsgTimestamp: 1543236572000,
		EncodingPath: "type:model/some/path",
		NodeId:       &telemetryBis.Telemetry_NodeIdStr{NodeIdStr: "hostname"},
		Subscription: &telemetryBis.Telemetry_SubscriptionIdStr{SubscriptionIdStr: "subscription"},
		DataGpbkv: []*telemetryBis.TelemetryField{
			{
				Fields: []*telemetryBis.TelemetryField{
					{
						Name: "keys",
						Fields: []*telemetryBis.TelemetryField{
							{
								Name:        "name",
								ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "str"},
							},
							{
								Name:        "uint64",
								ValueByType: &telemetryBis.TelemetryField_Uint64Value{Uint64Value: 1234},
							},
						},
					},
					{
						Name: "content",
						Fields: []*telemetryBis.TelemetryField{
							{
								Name:        "bool",
								ValueByType: &telemetryBis.TelemetryField_BoolValue{BoolValue: true},
							},
						},
					},
				},
			},
			{
				Fields: []*telemetryBis.TelemetryField{
					{
						Name: "keys",
						Fields: []*telemetryBis.TelemetryField{
							{
								Name:        "name",
								ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "str2"},
							},
						},
					},
					{
						Name: "content",
						Fields: []*telemetryBis.TelemetryField{
							{
								Name:        "bool",
								ValueByType: &telemetryBis.TelemetryField_BoolValue{BoolValue: false},
							},
						},
					},
				},
			},
		},
	}
	data, err := proto.Marshal(telemetry)
	require.NoError(t, err)

	c.handleTelemetry(data)
	require.Empty(t, acc.Errors)

	tags := map[string]string{
		"path":         "type:model/some/path",
		"name":         "str",
		"uint64":       "1234",
		"source":       "hostname",
		"subscription": "subscription",
	}
	fields := map[string]interface{}{"bool": true}
	acc.AssertContainsTaggedFields(t, "alias", fields, tags)

	tags = map[string]string{
		"path":         "type:model/some/path",
		"name":         "str2",
		"source":       "hostname",
		"subscription": "subscription",
	}
	fields = map[string]interface{}{"bool": false}
	acc.AssertContainsTaggedFields(t, "alias", fields, tags)
}

func TestIncludeDeleteField(t *testing.T) {
	type TelemetryEntry struct {
		name        string
		fieldName   string
		uint32Value uint32
		uint64Value uint64
		stringValue string
	}
	encodingPath := TelemetryEntry{
		name:        "path",
		stringValue: "openconfig-interfaces:interfaces/interface/subinterfaces/subinterface/openconfig-if-ip:ipv6/addresses/address",
	}
	name := TelemetryEntry{name: "name", stringValue: "Loopback10"}
	index := TelemetryEntry{name: "index", stringValue: "0"}
	ip := TelemetryEntry{name: "ip", fieldName: "state/ip", stringValue: "10::10"}
	prefixLength := TelemetryEntry{name: "prefix-length", fieldName: "state/prefix_length", uint32Value: uint32(128), uint64Value: 128}
	origin := TelemetryEntry{name: "origin", fieldName: "state/origin", stringValue: "STATIC"}
	status := TelemetryEntry{name: "status", fieldName: "state/status", stringValue: "PREFERRED"}
	source := TelemetryEntry{name: "source", stringValue: "hostname"}
	subscription := TelemetryEntry{name: "subscription", stringValue: "subscription"}
	deleteKey := "delete"
	stateKey := "state"

	testCases := []struct {
		telemetry *telemetryBis.Telemetry
		expected  []telegraf.Metric
	}{{
		telemetry: &telemetryBis.Telemetry{
			MsgTimestamp: 1543236572000,
			EncodingPath: encodingPath.stringValue,
			NodeId:       &telemetryBis.Telemetry_NodeIdStr{NodeIdStr: source.stringValue},
			Subscription: &telemetryBis.Telemetry_SubscriptionIdStr{SubscriptionIdStr: subscription.stringValue},
			DataGpbkv: []*telemetryBis.TelemetryField{
				{
					Fields: []*telemetryBis.TelemetryField{
						{
							Name: "keys",
							Fields: []*telemetryBis.TelemetryField{
								{
									Name:        name.name,
									ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: name.stringValue},
								},
								{
									Name:        index.name,
									ValueByType: &telemetryBis.TelemetryField_Uint32Value{Uint32Value: index.uint32Value},
								},
								{
									Name:        ip.name,
									ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: ip.stringValue},
								},
							},
						},
						{
							Name: "content",
							Fields: []*telemetryBis.TelemetryField{
								{
									Name: stateKey,
									Fields: []*telemetryBis.TelemetryField{
										{
											Name:        ip.name,
											ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: ip.stringValue},
										},
										{
											Name:        prefixLength.name,
											ValueByType: &telemetryBis.TelemetryField_Uint32Value{Uint32Value: prefixLength.uint32Value},
										},
										{
											Name:        origin.name,
											ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: origin.stringValue},
										},
										{
											Name:        status.name,
											ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: status.stringValue},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		expected: []telegraf.Metric{
			metric.New(
				"deleted",
				map[string]string{
					encodingPath.name: encodingPath.stringValue,
					name.name:         name.stringValue,
					index.name:        index.stringValue,
					ip.name:           ip.stringValue,
					source.name:       source.stringValue,
					subscription.name: subscription.stringValue,
				},
				map[string]interface{}{
					deleteKey:              false,
					ip.fieldName:           ip.stringValue,
					prefixLength.fieldName: prefixLength.uint64Value,
					origin.fieldName:       origin.stringValue,
					status.fieldName:       status.stringValue,
				},
				time.Now(),
			)},
	},
		{
			telemetry: &telemetryBis.Telemetry{
				MsgTimestamp: 1543236572000,
				EncodingPath: encodingPath.stringValue,
				NodeId:       &telemetryBis.Telemetry_NodeIdStr{NodeIdStr: source.stringValue},
				Subscription: &telemetryBis.Telemetry_SubscriptionIdStr{SubscriptionIdStr: subscription.stringValue},
				DataGpbkv: []*telemetryBis.TelemetryField{
					{
						Delete: true,
						Fields: []*telemetryBis.TelemetryField{
							{
								Name: "keys",
								Fields: []*telemetryBis.TelemetryField{
									{
										Name:        name.name,
										ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: name.stringValue},
									},
									{
										Name:        index.name,
										ValueByType: &telemetryBis.TelemetryField_Uint32Value{Uint32Value: index.uint32Value},
									},
									{
										Name:        ip.name,
										ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: ip.stringValue},
									},
								},
							},
						},
					},
				},
			},
			expected: []telegraf.Metric{
				metric.New(
					"deleted",
					map[string]string{
						encodingPath.name: encodingPath.stringValue,
						name.name:         name.stringValue,
						index.name:        index.stringValue,
						ip.name:           ip.stringValue,
						source.name:       source.stringValue,
						subscription.name: subscription.stringValue,
					},
					map[string]interface{}{deleteKey: true},
					time.Now(),
				)},
		},
	}
	for _, test := range testCases {
		c := &CiscoTelemetryMDT{
			Log:                testutil.Logger{},
			Transport:          "dummy",
			Aliases:            map[string]string{"deleted": encodingPath.stringValue},
			IncludeDeleteField: true}
		acc := &testutil.Accumulator{}
		// error is expected since we are passing in dummy transport
		require.ErrorContains(t, c.Start(acc), "dummy")
		data, err := proto.Marshal(test.telemetry)
		require.NoError(t, err)

		c.handleTelemetry(data)
		actual := acc.GetTelegrafMetrics()
		testutil.RequireMetricsEqual(t, test.expected, actual, testutil.IgnoreTime())
	}
}

func TestHandleTelemetrySingleNested(t *testing.T) {
	c := &CiscoTelemetryMDT{
		Log:       testutil.Logger{},
		Transport: "dummy",
		Aliases: map[string]string{
			"nested": "type:model/nested/path",
		},
	}
	acc := &testutil.Accumulator{}
	err := c.Start(acc)
	// error is expected since we are passing in dummy transport
	require.Error(t, err)

	telemetry := &telemetryBis.Telemetry{
		MsgTimestamp: 1543236572000,
		EncodingPath: "type:model/nested/path",
		NodeId:       &telemetryBis.Telemetry_NodeIdStr{NodeIdStr: "hostname"},
		Subscription: &telemetryBis.Telemetry_SubscriptionIdStr{SubscriptionIdStr: "subscription"},
		DataGpbkv: []*telemetryBis.TelemetryField{
			{
				Fields: []*telemetryBis.TelemetryField{
					{
						Name: "keys",
						Fields: []*telemetryBis.TelemetryField{
							{
								Name: "nested",
								Fields: []*telemetryBis.TelemetryField{
									{
										Name: "key",
										Fields: []*telemetryBis.TelemetryField{
											{
												Name:        "level",
												ValueByType: &telemetryBis.TelemetryField_DoubleValue{DoubleValue: 3},
											},
										},
									},
								},
							},
						},
					},
					{
						Name: "content",
						Fields: []*telemetryBis.TelemetryField{
							{
								Name: "nested",
								Fields: []*telemetryBis.TelemetryField{
									{
										Name: "value",
										Fields: []*telemetryBis.TelemetryField{
											{
												Name:        "foo",
												ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "bar"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	data, err := proto.Marshal(telemetry)
	require.NoError(t, err)

	c.handleTelemetry(data)
	require.Empty(t, acc.Errors)

	tags := map[string]string{
		"path":         "type:model/nested/path",
		"level":        "3",
		"source":       "hostname",
		"subscription": "subscription",
	}
	fields := map[string]interface{}{"nested/value/foo": "bar"}
	acc.AssertContainsTaggedFields(t, "nested", fields, tags)
}

func TestHandleEmbeddedTags(t *testing.T) {
	c := &CiscoTelemetryMDT{
		Transport:    "dummy",
		Aliases:      map[string]string{"extra": "type:model/extra"},
		EmbeddedTags: []string{"type:model/extra/list/name"},
	}
	acc := &testutil.Accumulator{}
	err := c.Start(acc)
	// error is expected since we are passing in dummy transport
	require.Error(t, err)

	telemetry := &telemetryBis.Telemetry{
		MsgTimestamp: 1543236572000,
		EncodingPath: "type:model/extra",
		NodeId:       &telemetryBis.Telemetry_NodeIdStr{NodeIdStr: "hostname"},
		Subscription: &telemetryBis.Telemetry_SubscriptionIdStr{SubscriptionIdStr: "subscription"},
		DataGpbkv: []*telemetryBis.TelemetryField{
			{
				Fields: []*telemetryBis.TelemetryField{
					{
						Name: "keys",
						Fields: []*telemetryBis.TelemetryField{
							{
								Name:        "foo",
								ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "bar"},
							},
						},
					},
					{
						Name: "content",
						Fields: []*telemetryBis.TelemetryField{
							{
								Name: "list",
								Fields: []*telemetryBis.TelemetryField{
									{
										Name:        "name",
										ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "entry1"},
									},
									{
										Name:        "test",
										ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "foo"},
									},
								},
							},
							{
								Name: "list",
								Fields: []*telemetryBis.TelemetryField{
									{
										Name:        "name",
										ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "entry2"},
									},
									{
										Name:        "test",
										ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "bar"},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	data, err := proto.Marshal(telemetry)
	require.NoError(t, err)

	c.handleTelemetry(data)
	require.Empty(t, acc.Errors)

	tags1 := map[string]string{
		"path":         "type:model/extra",
		"foo":          "bar",
		"source":       "hostname",
		"subscription": "subscription",
		"list/name":    "entry1",
	}
	fields1 := map[string]interface{}{"list/test": "foo"}
	tags2 := map[string]string{
		"path":         "type:model/extra",
		"foo":          "bar",
		"source":       "hostname",
		"subscription": "subscription",
		"list/name":    "entry2",
	}
	fields2 := map[string]interface{}{"list/test": "bar"}
	acc.AssertContainsTaggedFields(t, "extra", fields1, tags1)
	acc.AssertContainsTaggedFields(t, "extra", fields2, tags2)
}

func TestHandleNXAPI(t *testing.T) {
	c := &CiscoTelemetryMDT{
		Transport: "dummy",
		Aliases:   map[string]string{"nxapi": "show nxapi"},
	}
	acc := &testutil.Accumulator{}
	err := c.Start(acc)
	// error is expected since we are passing in dummy transport
	require.Error(t, err)

	telemetry := &telemetryBis.Telemetry{
		MsgTimestamp: 1543236572000,
		EncodingPath: "show nxapi",
		NodeId:       &telemetryBis.Telemetry_NodeIdStr{NodeIdStr: "hostname"},
		Subscription: &telemetryBis.Telemetry_SubscriptionIdStr{SubscriptionIdStr: "subscription"},
		DataGpbkv: []*telemetryBis.TelemetryField{
			{
				Fields: []*telemetryBis.TelemetryField{
					{
						Name: "keys",
						Fields: []*telemetryBis.TelemetryField{
							{
								Name:        "foo",
								ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "bar"},
							},
						},
					},
					{
						Name: "content",
						Fields: []*telemetryBis.TelemetryField{
							{
								Fields: []*telemetryBis.TelemetryField{
									{
										Name: "TABLE_nxapi",
										Fields: []*telemetryBis.TelemetryField{
											{
												Fields: []*telemetryBis.TelemetryField{
													{
														Name: "ROW_nxapi",
														Fields: []*telemetryBis.TelemetryField{
															{
																Fields: []*telemetryBis.TelemetryField{
																	{
																		Name:        "index",
																		ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "i1"},
																	},
																	{
																		Name:        "value",
																		ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "foo"},
																	},
																},
															},
															{
																Fields: []*telemetryBis.TelemetryField{
																	{
																		Name:        "index",
																		ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "i2"},
																	},
																	{
																		Name:        "value",
																		ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "bar"},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	data, err := proto.Marshal(telemetry)
	require.NoError(t, err)

	c.handleTelemetry(data)
	require.Empty(t, acc.Errors)

	tags1 := map[string]string{
		"path":         "show nxapi",
		"foo":          "bar",
		"TABLE_nxapi":  "i1",
		"row_number":   "0",
		"source":       "hostname",
		"subscription": "subscription",
	}
	fields1 := map[string]interface{}{"value": "foo"}
	tags2 := map[string]string{
		"path":         "show nxapi",
		"foo":          "bar",
		"TABLE_nxapi":  "i2",
		"row_number":   "0",
		"source":       "hostname",
		"subscription": "subscription",
	}
	fields2 := map[string]interface{}{"value": "bar"}
	acc.AssertContainsTaggedFields(t, "nxapi", fields1, tags1)
	acc.AssertContainsTaggedFields(t, "nxapi", fields2, tags2)
}

func TestHandleNXAPIXformNXAPI(t *testing.T) {
	c := &CiscoTelemetryMDT{
		Log:       testutil.Logger{},
		Transport: "dummy",
		Aliases:   map[string]string{"nxapi": "show nxapi"},
	}
	acc := &testutil.Accumulator{}
	err := c.Start(acc)
	// error is expected since we are passing in dummy transport
	require.Error(t, err)

	telemetry := &telemetryBis.Telemetry{
		MsgTimestamp: 1543236572000,
		EncodingPath: "show processes cpu",
		NodeId:       &telemetryBis.Telemetry_NodeIdStr{NodeIdStr: "hostname"},
		Subscription: &telemetryBis.Telemetry_SubscriptionIdStr{SubscriptionIdStr: "subscription"},
		DataGpbkv: []*telemetryBis.TelemetryField{
			{
				Fields: []*telemetryBis.TelemetryField{
					{
						Name: "keys",
						Fields: []*telemetryBis.TelemetryField{
							{
								Name:        "foo",
								ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "bar"},
							},
						},
					},
					{
						Name: "content",
						Fields: []*telemetryBis.TelemetryField{
							{
								Fields: []*telemetryBis.TelemetryField{
									{
										Name: "TABLE_process_cpu",
										Fields: []*telemetryBis.TelemetryField{
											{
												Fields: []*telemetryBis.TelemetryField{
													{
														Name: "ROW_process_cpu",
														Fields: []*telemetryBis.TelemetryField{
															{
																Fields: []*telemetryBis.TelemetryField{
																	{
																		Name:        "index",
																		ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "i1"},
																	},
																	{
																		Name:        "value",
																		ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "foo"},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	data, err := proto.Marshal(telemetry)
	require.NoError(t, err)

	c.handleTelemetry(data)
	require.Empty(t, acc.Errors)

	tags1 := map[string]string{
		"path":              "show processes cpu",
		"foo":               "bar",
		"TABLE_process_cpu": "i1",
		"row_number":        "0",
		"source":            "hostname",
		"subscription":      "subscription",
	}
	fields1 := map[string]interface{}{"value": "foo"}
	acc.AssertContainsTaggedFields(t, "show processes cpu", fields1, tags1)
}

func TestHandleNXXformMulti(t *testing.T) {
	c := &CiscoTelemetryMDT{
		Transport: "dummy",
		Aliases:   map[string]string{"dme": "sys/lldp"},
	}
	acc := &testutil.Accumulator{}
	err := c.Start(acc)
	// error is expected since we are passing in dummy transport
	require.Error(t, err)

	telemetry := &telemetryBis.Telemetry{
		MsgTimestamp: 1543236572000,
		EncodingPath: "sys/lldp",
		NodeId:       &telemetryBis.Telemetry_NodeIdStr{NodeIdStr: "hostname"},
		Subscription: &telemetryBis.Telemetry_SubscriptionIdStr{SubscriptionIdStr: "subscription"},
		DataGpbkv: []*telemetryBis.TelemetryField{
			{
				Fields: []*telemetryBis.TelemetryField{
					{
						Name: "keys",
						Fields: []*telemetryBis.TelemetryField{
							{
								Name:        "foo",
								ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "bar"},
							},
						},
					},
					{
						Name: "content",
						Fields: []*telemetryBis.TelemetryField{
							{
								Fields: []*telemetryBis.TelemetryField{
									{
										Name: "fooEntity",
										Fields: []*telemetryBis.TelemetryField{
											{
												Fields: []*telemetryBis.TelemetryField{
													{
														Name: "attributes",
														Fields: []*telemetryBis.TelemetryField{
															{
																Fields: []*telemetryBis.TelemetryField{
																	{
																		Name:        "rn",
																		ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "some-rn"},
																	},
																	{
																		Name:        "portIdV",
																		ValueByType: &telemetryBis.TelemetryField_Uint32Value{Uint32Value: 12},
																	},
																	{
																		Name:        "portDesc",
																		ValueByType: &telemetryBis.TelemetryField_Uint64Value{Uint64Value: 100},
																	},
																	{
																		Name:        "test",
																		ValueByType: &telemetryBis.TelemetryField_Uint64Value{Uint64Value: 281474976710655},
																	},
																	{
																		Name:        "subscriptionId",
																		ValueByType: &telemetryBis.TelemetryField_Uint64Value{Uint64Value: 2814749767106551},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	data, err := proto.Marshal(telemetry)
	require.NoError(t, err)

	c.handleTelemetry(data)
	require.Empty(t, acc.Errors)
	//validate various transformation scenaarios newly added in the code.
	fields := map[string]interface{}{
		"portIdV":        "12",
		"portDesc":       "100",
		"test":           int64(281474976710655),
		"subscriptionId": "2814749767106551",
	}
	acc.AssertContainsFields(t, "dme", fields)
}

func TestHandleNXDME(t *testing.T) {
	c := &CiscoTelemetryMDT{
		Transport: "dummy",
		Aliases:   map[string]string{"dme": "sys/dme"},
	}
	acc := &testutil.Accumulator{}
	err := c.Start(acc)
	// error is expected since we are passing in dummy transport
	require.Error(t, err)

	telemetry := &telemetryBis.Telemetry{
		MsgTimestamp: 1543236572000,
		EncodingPath: "sys/dme",
		NodeId:       &telemetryBis.Telemetry_NodeIdStr{NodeIdStr: "hostname"},
		Subscription: &telemetryBis.Telemetry_SubscriptionIdStr{SubscriptionIdStr: "subscription"},
		DataGpbkv: []*telemetryBis.TelemetryField{
			{
				Fields: []*telemetryBis.TelemetryField{
					{
						Name: "keys",
						Fields: []*telemetryBis.TelemetryField{
							{
								Name:        "foo",
								ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "bar"},
							},
						},
					},
					{
						Name: "content",
						Fields: []*telemetryBis.TelemetryField{
							{
								Fields: []*telemetryBis.TelemetryField{
									{
										Name: "fooEntity",
										Fields: []*telemetryBis.TelemetryField{
											{
												Fields: []*telemetryBis.TelemetryField{
													{
														Name: "attributes",
														Fields: []*telemetryBis.TelemetryField{
															{
																Fields: []*telemetryBis.TelemetryField{
																	{
																		Name:        "rn",
																		ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "some-rn"},
																	},
																	{
																		Name:        "value",
																		ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "foo"},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	data, err := proto.Marshal(telemetry)
	require.NoError(t, err)

	c.handleTelemetry(data)
	require.Empty(t, acc.Errors)

	tags1 := map[string]string{
		"path":         "sys/dme",
		"foo":          "bar",
		"fooEntity":    "some-rn",
		"source":       "hostname",
		"subscription": "subscription",
	}
	fields1 := map[string]interface{}{"value": "foo"}
	acc.AssertContainsTaggedFields(t, "dme", fields1, tags1)
}

func TestTCPDialoutOverflow(t *testing.T) {
	c := &CiscoTelemetryMDT{
		Log:            testutil.Logger{},
		Transport:      "tcp",
		ServiceAddress: "127.0.0.1:0",
	}
	acc := &testutil.Accumulator{}
	err := c.Start(acc)
	require.NoError(t, err)

	hdr := struct {
		MsgType       uint16
		MsgEncap      uint16
		MsgHdrVersion uint16
		MsgFlags      uint16
		MsgLen        uint32
	}{MsgLen: uint32(1000000000)}

	addr := c.Address()
	conn, err := net.Dial(addr.Network(), addr.String())
	require.NoError(t, err)
	require.NoError(t, binary.Write(conn, binary.BigEndian, hdr))
	_, err = conn.Read([]byte{0})
	require.True(t, err == nil || errors.Is(err, io.EOF))
	require.NoError(t, conn.Close())

	c.Stop()

	require.Contains(t, acc.Errors, errors.New("dialout packet too long: 1000000000"))
}

func mockTelemetryMicroburstMessage() *telemetryBis.Telemetry {
	data, err := os.ReadFile("./testdata/microburst")
	if err != nil {
		panic(err)
	}

	newMessage := &telemetryBis.Telemetry{}
	err = proto.Unmarshal(data, newMessage)
	if err != nil {
		panic(err)
	}
	return newMessage
}

func mockTelemetryMessage() *telemetryBis.Telemetry {
	return &telemetryBis.Telemetry{
		MsgTimestamp: 1543236572000,
		EncodingPath: "type:model/some/path",
		NodeId:       &telemetryBis.Telemetry_NodeIdStr{NodeIdStr: "hostname"},
		Subscription: &telemetryBis.Telemetry_SubscriptionIdStr{SubscriptionIdStr: "subscription"},
		DataGpbkv: []*telemetryBis.TelemetryField{
			{
				Fields: []*telemetryBis.TelemetryField{
					{
						Name: "keys",
						Fields: []*telemetryBis.TelemetryField{
							{
								Name:        "name",
								ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "str"},
							},
						},
					},
					{
						Name: "content",
						Fields: []*telemetryBis.TelemetryField{
							{
								Name:        "value",
								ValueByType: &telemetryBis.TelemetryField_Sint64Value{Sint64Value: -1},
							},
						},
					},
				},
			},
		},
	}
}

func TestGRPCDialoutMicroburst(t *testing.T) {
	c := &CiscoTelemetryMDT{
		Log:            testutil.Logger{},
		Transport:      "grpc",
		ServiceAddress: "127.0.0.1:0",
		Aliases: map[string]string{
			"some":     "microburst",
			"parallel": "type:model/parallel/path",
			"other":    "type:model/other/path",
		},
	}
	acc := &testutil.Accumulator{}
	err := c.Start(acc)
	require.NoError(t, err)

	telemetry := mockTelemetryMicroburstMessage()
	data, err := proto.Marshal(telemetry)
	require.NoError(t, err)

	c.handleTelemetry(data)
	require.Empty(t, acc.Errors)
	tags := map[string]string{
		"microburst":   "microburst",
		"path":         "microburst",
		"source":       "n9k-eor-tm4",
		"subscription": "1",
	}
	fields := map[string]interface{}{
		"duration":      uint64(1200),
		"endDepth":      int64(0),
		"interfaceName": "Eth0/0/0",
		"peak":          int64(500),
		"queue":         "queue-255",
		"queueType":     "unicast",
		"threshold":     int64(0),
		"ts":            "2023-08-03T20:12:59.655308Z",
	}
	acc.AssertContainsTaggedFields(t, "microburst", fields, tags)
}

func TestTCPDialoutMultiple(t *testing.T) {
	c := &CiscoTelemetryMDT{
		Log:            testutil.Logger{},
		Transport:      "tcp",
		ServiceAddress: "127.0.0.1:0",
		Aliases: map[string]string{
			"some":     "type:model/some/path",
			"parallel": "type:model/parallel/path",
			"other":    "type:model/other/path",
		},
	}
	acc := &testutil.Accumulator{}
	err := c.Start(acc)
	require.NoError(t, err)

	telemetry := mockTelemetryMessage()

	hdr := struct {
		MsgType       uint16
		MsgEncap      uint16
		MsgHdrVersion uint16
		MsgFlags      uint16
		MsgLen        uint32
	}{}

	addr := c.Address()
	conn, err := net.Dial(addr.Network(), addr.String())
	require.NoError(t, err)

	data, err := proto.Marshal(telemetry)
	require.NoError(t, err)
	hdr.MsgLen = uint32(len(data))
	require.NoError(t, binary.Write(conn, binary.BigEndian, hdr))
	_, err = conn.Write(data)
	require.NoError(t, err)

	conn2, err := net.Dial(addr.Network(), addr.String())
	require.NoError(t, err)

	telemetry.EncodingPath = "type:model/parallel/path"
	data, err = proto.Marshal(telemetry)
	require.NoError(t, err)
	hdr.MsgLen = uint32(len(data))
	require.NoError(t, binary.Write(conn2, binary.BigEndian, hdr))
	_, err = conn2.Write(data)
	require.NoError(t, err)
	_, err = conn2.Write([]byte{0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0})
	require.NoError(t, err)
	_, err = conn2.Read([]byte{0})
	require.True(t, err == nil || errors.Is(err, io.EOF))
	require.NoError(t, conn2.Close())

	telemetry.EncodingPath = "type:model/other/path"
	data, err = proto.Marshal(telemetry)
	require.NoError(t, err)
	hdr.MsgLen = uint32(len(data))
	require.NoError(t, binary.Write(conn, binary.BigEndian, hdr))
	_, err = conn.Write(data)
	require.NoError(t, err)
	_, err = conn.Write([]byte{0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0})
	require.NoError(t, err)
	_, err = conn.Read([]byte{0})
	require.True(t, err == nil || errors.Is(err, io.EOF))
	c.Stop()
	require.NoError(t, conn.Close())

	// We use the invalid dialout flags to let the server close the connection
	require.Equal(t, []error{errors.New("invalid dialout flags: 257"), errors.New("invalid dialout flags: 257")}, acc.Errors)

	tags := map[string]string{
		"path":         "type:model/some/path",
		"name":         "str",
		"source":       "hostname",
		"subscription": "subscription",
	}
	fields := map[string]interface{}{"value": int64(-1)}
	acc.AssertContainsTaggedFields(t, "some", fields, tags)

	tags = map[string]string{
		"path":         "type:model/parallel/path",
		"name":         "str",
		"source":       "hostname",
		"subscription": "subscription",
	}
	fields = map[string]interface{}{"value": int64(-1)}
	acc.AssertContainsTaggedFields(t, "parallel", fields, tags)

	tags = map[string]string{
		"path":         "type:model/other/path",
		"name":         "str",
		"source":       "hostname",
		"subscription": "subscription",
	}
	fields = map[string]interface{}{"value": int64(-1)}
	acc.AssertContainsTaggedFields(t, "other", fields, tags)
}

func TestGRPCDialoutError(t *testing.T) {
	c := &CiscoTelemetryMDT{
		Log:            testutil.Logger{},
		Transport:      "grpc",
		ServiceAddress: "127.0.0.1:0",
	}
	acc := &testutil.Accumulator{}
	err := c.Start(acc)
	require.NoError(t, err)

	addr := c.Address()
	conn, err := grpc.Dial(addr.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	client := dialout.NewGRPCMdtDialoutClient(conn)
	stream, err := client.MdtDialout(context.Background())
	require.NoError(t, err)

	args := &dialout.MdtDialoutArgs{Errors: "foobar"}
	require.NoError(t, stream.Send(args))

	// Wait for the server to close
	_, err = stream.Recv()
	require.True(t, err == nil || errors.Is(err, io.EOF))
	c.Stop()

	require.Equal(t, []error{errors.New("GRPC dialout error: foobar")}, acc.Errors)
}

func TestGRPCDialoutMultiple(t *testing.T) {
	c := &CiscoTelemetryMDT{
		Log:            testutil.Logger{},
		Transport:      "grpc",
		ServiceAddress: "127.0.0.1:0",
		Aliases: map[string]string{
			"some":     "type:model/some/path",
			"parallel": "type:model/parallel/path",
			"other":    "type:model/other/path",
		},
	}
	acc := &testutil.Accumulator{}
	err := c.Start(acc)
	require.NoError(t, err)
	telemetry := mockTelemetryMessage()

	addr := c.Address()
	conn, err := grpc.Dial(addr.String(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	require.NoError(t, err)
	client := dialout.NewGRPCMdtDialoutClient(conn)
	stream, err := client.MdtDialout(context.TODO())
	require.NoError(t, err)

	data, err := proto.Marshal(telemetry)
	require.NoError(t, err)
	args := &dialout.MdtDialoutArgs{Data: data, ReqId: 456}
	require.NoError(t, stream.Send(args))

	conn2, err := grpc.Dial(addr.String(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	require.NoError(t, err)
	client2 := dialout.NewGRPCMdtDialoutClient(conn2)
	stream2, err := client2.MdtDialout(context.TODO())
	require.NoError(t, err)

	telemetry.EncodingPath = "type:model/parallel/path"
	data, err = proto.Marshal(telemetry)
	require.NoError(t, err)
	args = &dialout.MdtDialoutArgs{Data: data}
	require.NoError(t, stream2.Send(args))
	require.NoError(t, stream2.Send(&dialout.MdtDialoutArgs{Errors: "testclose"}))
	_, err = stream2.Recv()
	require.True(t, err == nil || errors.Is(err, io.EOF))
	require.NoError(t, conn2.Close())

	telemetry.EncodingPath = "type:model/other/path"
	data, err = proto.Marshal(telemetry)
	require.NoError(t, err)
	args = &dialout.MdtDialoutArgs{Data: data}
	require.NoError(t, stream.Send(args))
	require.NoError(t, stream.Send(&dialout.MdtDialoutArgs{Errors: "testclose"}))
	_, err = stream.Recv()
	require.True(t, err == nil || errors.Is(err, io.EOF))

	c.Stop()
	require.NoError(t, conn.Close())

	require.Equal(t, []error{errors.New("GRPC dialout error: testclose"), errors.New("GRPC dialout error: testclose")}, acc.Errors)

	tags := map[string]string{
		"path":         "type:model/some/path",
		"name":         "str",
		"source":       "hostname",
		"subscription": "subscription",
	}
	fields := map[string]interface{}{"value": int64(-1)}
	acc.AssertContainsTaggedFields(t, "some", fields, tags)

	tags = map[string]string{
		"path":         "type:model/parallel/path",
		"name":         "str",
		"source":       "hostname",
		"subscription": "subscription",
	}
	fields = map[string]interface{}{"value": int64(-1)}
	acc.AssertContainsTaggedFields(t, "parallel", fields, tags)

	tags = map[string]string{
		"path":         "type:model/other/path",
		"name":         "str",
		"source":       "hostname",
		"subscription": "subscription",
	}
	fields = map[string]interface{}{"value": int64(-1)}
	acc.AssertContainsTaggedFields(t, "other", fields, tags)
}

func TestGRPCDialoutKeepalive(t *testing.T) {
	c := &CiscoTelemetryMDT{
		Log:            testutil.Logger{},
		Transport:      "grpc",
		ServiceAddress: "127.0.0.1:0",
		EnforcementPolicy: GRPCEnforcementPolicy{
			PermitKeepaliveWithoutCalls: true,
			KeepaliveMinTime:            0,
		},
	}
	acc := &testutil.Accumulator{}
	err := c.Start(acc)
	require.NoError(t, err)

	addr := c.Address()
	conn, err := grpc.Dial(addr.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	client := dialout.NewGRPCMdtDialoutClient(conn)
	stream, err := client.MdtDialout(context.Background())
	require.NoError(t, err)

	telemetry := mockTelemetryMessage()
	data, err := proto.Marshal(telemetry)
	require.NoError(t, err)
	args := &dialout.MdtDialoutArgs{Data: data, ReqId: 456}
	require.NoError(t, stream.Send(args))

	c.Stop()
	require.NoError(t, conn.Close())
}

func TestSourceFieldRewrite(t *testing.T) {
	c := &CiscoTelemetryMDT{
		Log:       testutil.Logger{},
		Transport: "dummy",
		Aliases:   map[string]string{"alias": "type:model/some/path"},
	}
	c.SourceFieldName = "mdt_source"
	acc := &testutil.Accumulator{}
	err := c.Start(acc)
	// error is expected since we are passing in dummy transport
	require.Error(t, err)

	telemetry := &telemetryBis.Telemetry{
		MsgTimestamp: 1543236572000,
		EncodingPath: "type:model/some/path",
		NodeId:       &telemetryBis.Telemetry_NodeIdStr{NodeIdStr: "hostname"},
		Subscription: &telemetryBis.Telemetry_SubscriptionIdStr{SubscriptionIdStr: "subscription"},
		DataGpbkv: []*telemetryBis.TelemetryField{
			{
				Fields: []*telemetryBis.TelemetryField{
					{
						Name: "keys",
						Fields: []*telemetryBis.TelemetryField{
							{
								Name:        "source",
								ValueByType: &telemetryBis.TelemetryField_StringValue{StringValue: "str"},
							},
						},
					},
					{
						Name: "content",
						Fields: []*telemetryBis.TelemetryField{
							{
								Name:        "bool",
								ValueByType: &telemetryBis.TelemetryField_BoolValue{BoolValue: false},
							},
						},
					},
				},
			},
		},
	}
	data, err := proto.Marshal(telemetry)
	require.NoError(t, err)

	c.handleTelemetry(data)
	require.Empty(t, acc.Errors)

	tags := map[string]string{
		"path":         "type:model/some/path",
		"mdt_source":   "str",
		"source":       "hostname",
		"subscription": "subscription",
	}
	fields := map[string]interface{}{"bool": false}
	acc.AssertContainsTaggedFields(t, "alias", fields, tags)
}
