package eureka

import (
	"strconv"
	"github.com/hudl/fargo"
	"github.com/satori/go.uuid"
	"time"
	"net/http"
	"log"
	"strings"
	"os"
	"os/signal"
	"syscall"
)


var _INSTANCEID string = uuid.NewV4().String()
var _HEARTBEAT_MAX_CONSECUTIVE_ERRORS int = 5
var _HEARTBEAT_SLEEPTIMEBETWEENHEARTBEATINSECONDS time.Duration = 10

var _SECUREPORT int = 8443
var _DATACENTER_NAME string = "MyOwn"

var _Configuration map[string]interface{}

type MyStruct struct {
	Name string `yml:"name"`
}


func Register(cfg map[string]interface{}) {
	_Configuration = cfg



	eurekaUrl := cleanEurekaUrlIfNeeded(_Configuration["eureka.client.serviceUrl.defaultZone"].(string))

	conn := fargo.NewConn(eurekaUrl)
	instance := new(fargo.Instance)
	instance.App = _Configuration["spring.application.name"].(string)
	instance.DataCenterInfo.Name = _DATACENTER_NAME
	instance.HealthCheckUrl = "http://" + _Configuration["eureka.instance.ip-address"].(string) + ":" + strconv.Itoa(_Configuration["server.port"].(int)) + "/health"
	instance.HomePageUrl = "http://" + _Configuration["eureka.instance.ip-address"].(string) + ":" + strconv.Itoa(_Configuration["server.port"].(int)) + "/"
	instance.StatusPageUrl = "http://" + _Configuration["eureka.instance.ip-address"].(string) + ":" + strconv.Itoa(_Configuration["server.port"].(int)) + "/info"
	instance.IPAddr = _Configuration["eureka.instance.ip-address"].(string)
	instance.HostName = _Configuration["hostname"].(string)
	instance.SecurePort = _SECUREPORT
	instance.SecureVipAddress = _Configuration["spring.application.name"].(string)
	instance.VipAddress = _Configuration["spring.application.name"].(string)
	instance.Status = fargo.StatusType("UP")
	instance.SetMetadataString("instanceId", _INSTANCEID)

	err := conn.RegisterInstance(instance)

	if err != nil {
		log.Println("cannot register in eureka")
	}

	startHeartbeat(eurekaUrl, _Configuration["spring.application.name"].(string), _Configuration["hostname"].(string), _INSTANCEID)

}

func cleanEurekaUrlIfNeeded(eurekaUrl string) string {
	newEurekaUrl := strings.Split(eurekaUrl, ",")[0]
	if newEurekaUrl[len(newEurekaUrl)-1:] == "/" {
		newEurekaUrl = newEurekaUrl[:len(newEurekaUrl)-1]
	}
	return newEurekaUrl
}

func Deregister() {

	eurekaUrl := cleanEurekaUrlIfNeeded(_Configuration["eureka.client.serviceUrl.defaultZone"].(string)) + "/apps/" + _Configuration["spring.application.name"].(string) + "/" + _Configuration["hostname"].(string) +  ":" + _INSTANCEID
	req, _ := http.NewRequest(http.MethodDelete, eurekaUrl, nil)
	res, _ := http.DefaultClient.Do(req)
	if res.StatusCode == http.StatusOK {
		log.Println("Deregistered correctly")
	} else {
		log.Println("Error while deregistering")
	}
}

func startHeartbeat(eurekaUrl string, appName string, hostname string, instance string) {
	consecutiveErrors := 0
	for {
		url := eurekaUrl + "/apps/" + appName + "/" + hostname + ":" + instance

		req, _ := http.NewRequest("PUT", url, nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != http.StatusOK {
			consecutiveErrors++
			if consecutiveErrors >= _HEARTBEAT_MAX_CONSECUTIVE_ERRORS {
				Deregister()
				Register(_Configuration)
			}
		}

		res.Body.Close()
		time.Sleep(_HEARTBEAT_SLEEPTIMEBETWEENHEARTBEATINSECONDS * time.Second)
	}
}

func CaptureInterruptSignal() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}