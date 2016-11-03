package actuator

import (
	"encoding/json"
	"runtime"
)

func metrics() string {

	values := refillMetricsMap()

	b, err := json.Marshal(values)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func refillMetricsMap() runtime.MemStats {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return mem
}
