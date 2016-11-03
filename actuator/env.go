package actuator

import (
	"encoding/json"
	"os"
	"strings"
)

func env() string {
	env := new(envObject)
	env.SystemEnvironment = make(map[string]string)

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		env.SystemEnvironment[pair[0]] = pair[1]
	}


	b, err := json.Marshal(env)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

type envObject struct {
	Configuration map[string]interface{} `json:"configuration"`
	SystemEnvironment map[string]string `json:"systemEnvironment"`
}