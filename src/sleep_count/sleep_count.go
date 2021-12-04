package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

type DarwinSleepCountPlugin struct {
	Prefix string
}

func (d DarwinSleepCountPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(d.MetricKeyPrefix())

	return map[string]mp.Graphs{
		"sleep": {
			Label: labelPrefix,
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "sleep", Label: "Sleep Count"},
				{Name: "dark_wake", Label: "Dark Wake Count"},
				{Name: "user_wake", Label: "User Wake Count"},
			},
		},
	}
}

func (d DarwinSleepCountPlugin) FetchMetrics() (map[string]float64, error) {
	// $ pmset -g stats
	// Sleep Count:862
	// Dark Wake Count:853
	// User Wake Count:15
	output, err := exec.Command("pmset", "-g", "stats").Output()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch pmset metrics: %s", err)
	}
	lines := strings.Split(string(output), "\n")
	var sleep_count, dark_wake_count, user_wake_count float64

	if strings.HasPrefix(lines[0], "Sleep Count") {
		s := strings.Split(lines[0], ":")
		sleep_count, err = strconv.ParseFloat(s[1], 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse Sleep Count: %s", err)
		}
	} else {
		return nil, fmt.Errorf("Failed to fetch Sleep Count")
	}

	if strings.HasPrefix(lines[1], "Dark Wake Count") {
		s := strings.Split(lines[1], ":")
		dark_wake_count, err = strconv.ParseFloat(s[1], 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse Dark Wake Count: %s", err)
		}
	} else {
		return nil, fmt.Errorf("Failed to fetch Dark Wake Count")
	}

	if strings.HasPrefix(lines[2], "User Wake Count") {
		s := strings.Split(lines[2], ":")
		user_wake_count, err = strconv.ParseFloat(s[1], 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse User Wake Count: %s", err)
		}
	} else {
		return nil, fmt.Errorf("Failed to fetch User Wake Count")
	}
	return map[string]float64{
		"sleep":     sleep_count,
		"dark_wake": dark_wake_count,
		"user_wake": user_wake_count}, nil
}

func (d DarwinSleepCountPlugin) MetricKeyPrefix() string {
	if d.Prefix == "" {
		d.Prefix = "battery"
	}
	return d.Prefix
}

func main() {
	optPrefix := flag.String("metric-key-prefix", "sleep", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	u := DarwinSleepCountPlugin{
		Prefix: *optPrefix,
	}
	plugin := mp.NewMackerelPlugin(u)
	plugin.Tempfile = *optTempfile
	plugin.Run()
}
