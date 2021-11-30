package grafana

import (
	"context"
	"testing"
)

func TestCurrentValues(t *testing.T) {
	var (
		queriesTestDataMap = map[string]string{
			"update_global_index_gas_used{}":   "gas used",
			"update_global_index_gas_wanted{}": "gas wanted",
			"update_global_index_uusd_fee{}":   "uusd fee",
			"config_crc32":                     "{{label}}",
		}

		queriesTestData = []expr{
			{
				Query:        "update_global_index_gas_used{}",
				LegendFormat: "gas used",
			},
			{
				Query:        "update_global_index_gas_wanted{}",
				LegendFormat: "gas wanted",
			},
			{
				Query:        "update_global_index_uusd_fee{}",
				LegendFormat: "uusd fee",
			},
			{
				Query:        "config_crc32",
				LegendFormat: "{{label}}",
			},
		}
	)

	client := newClient(httpPrefix+addr, token, timeout)

	datasources, err := client.currentValues(context.Background(), queriesTestData)
	if err != nil {
		t.Fatal(err)
	}

	if len(datasources) == 0 {
		t.Fatal("current values is empty!")
	}

	for _, d := range datasources {
		switch d.Query {
		case "update_global_index_gas_used{}":
			checkSingleValues(t, d.Values, queriesTestDataMap[d.Query])
		case "update_global_index_gas_wanted{}":
			checkSingleValues(t, d.Values, queriesTestDataMap[d.Query])
		case "update_global_index_uusd_fee{}":
			checkSingleValues(t, d.Values, queriesTestDataMap[d.Query])
		case "config_crc32":
			checkMultipleValues(t, d.Values, queriesTestDataMap[d.Query])
		}
	}
}

func checkSingleValues(t *testing.T, values []LabelValue, label string) {
	if len(values) > 1 {
		t.Fatal("lenght of values more than 1, but should be 1!")
	}

	if values[0].Label != label {
		t.Fatal("wrong label!")
	}

	if values[0].Value == "" {
		t.Fatal("empty value!")
	}
}

func checkMultipleValues(t *testing.T, values []LabelValue, wrongLabel string) {
	if len(values) == 1 {
		t.Fatal("lenght of values should be more than 1")
	}

	for _, v := range values {
		if v.Label == wrongLabel || v.Label == "" {
			t.Fatal("wrong label!")
		}

		if v.Value == "" {
			t.Fatal("empty value!")
		}
	}
}
