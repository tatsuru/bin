package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

type DarwinBatteryPlugin struct {
	Prefix string
}

func (d DarwinBatteryPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(d.MetricKeyPrefix())

	return map[string]mp.Graphs{
		"percentage": {
			Label: labelPrefix,
			Unit:  mp.UnitPercentage,
			Metrics: []mp.Metrics{
				{Name: "percentage", Label: "Battery1 Percentage"},
			},
		},
	}
}

func (d DarwinBatteryPlugin) FetchMetrics() (map[string]float64, error) {
	// $ pmset -g batt
	// Now drawing from 'Battery Power'
	// -InternalBattery-0 (id=19202147)       99%; discharging; 10:21 remaining present: true
	output, err := exec.Command("pmset", "-g", "batt").Output()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch pmset metrics: %s", err)
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) != 3 {
		return nil, fmt.Errorf("more than 1 battery?: %s", err)
	}
	xs := strings.Fields(lines[1])
	pct_s := strings.ReplaceAll(xs[2], "%;", "")
	pct, err := strconv.ParseFloat(pct_s, 64)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch pmset metrics: %s", err)
	}

	return map[string]float64{"percentage": pct}, nil
}

func (d DarwinBatteryPlugin) MetricKeyPrefix() string {
	if d.Prefix == "" {
		d.Prefix = "battery"
	}
	return d.Prefix
}

func main() {
	optPrefix := flag.String("metric-key-prefix", "battery", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	u := DarwinBatteryPlugin{
		Prefix: *optPrefix,
	}
	plugin := mp.NewMackerelPlugin(u)
	plugin.Tempfile = *optTempfile
	plugin.Run()
}
