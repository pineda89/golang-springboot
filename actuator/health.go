package actuator

import "encoding/json"

func health() string {
	healthJson := new(healthJson)

	healthJson.Status = "UP"

	b, err := json.Marshal(healthJson)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

type healthJson struct {
	Status string `json:"status"`
}

