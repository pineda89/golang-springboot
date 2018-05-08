# golang-starter

## Eureka integration is deprecated

golang-springboot is a project than implements with spring-cloud-config, eureka and actuator like a springboot implementation

The project parts are:
* actuator - exposes some methods like springboot actuator (/health, /metrics, /env, /info)
* config - connects to "spring cloud config server" and loads the configuration
* eureka - register, deregister and Heartbeat with eureka service

### How use it?

we have two options for starting:
* using env vars
* using command line params

#### env vars

we must export the nexts env vars:

| env vars | description | example |
| ---------- | ---------- | ---------- |
| spring_profiles_active | spring cloud config profile | test
| spring_cloud_config_uri | uri of spring cloud config server | http://127.0.0.1:30606
| spring_cloud_config_label | spring cloud config label | test
| server_port | exposed port | 8080
| eureka_instance_ip_address | IP of this application, accesible by eureka | 10.0.1.5
| spring_application_name | name of application | my-golang-application

#### command line params

we report the following parameters in the same order:
* spring_profiles_active
* spring_cloud_config_uri
* spring_cloud_config_label
* server_port
* eureka_instance_ip_address
* spring_application_name

example:

```golang
go run main.go pre http://localhost:8080 pre 8080 10.5.5.10 my-golang-app
```

### Import it in your project

Just download and compile the code. Also you must call the methods for start it

```golang
package main

import (
	"github.com/pineda89/golang-springboot/actuator"
	"github.com/pineda89/golang-springboot/config"
	"github.com/pineda89/golang-springboot/eureka"
)

func main() {
	config.LoadConfig()

	go StartWebServer(config.Configuration["server.port"].(int))
	go actuator.InitializeActuator()

	go eureka.Register(config.Configuration)
	CaptureInterruptSignal()
	eureka.Deregister()
}

func StartWebServer(port int) {
	http.HandleFunc("/", mainHandler)
	http.ListenAndServe(":" + strconv.Itoa(port), nil)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, service.MyApplicationMethod())
}

func CaptureInterruptSignal() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

```

This implementation exposes the methods "/", "/info", "/health", "/metrics" and "/env"
Also loads configuration from spring cloud config server, and register, deregister and healthcheck on eureka

