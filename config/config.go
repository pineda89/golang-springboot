package config

import (
	"os"
	"net/http"
	"github.com/Jeffail/gabs"
	"io/ioutil"
	"reflect"
	"strings"
	"strconv"
	"log"
)

var _DEFAULT_PORT int = 8080

var Configuration map[string]interface{} = make(map[string]interface{})


func LoadConfig() {
	log.Println("Loading config...")
	params := preloadConfigurationParams()
	newConfig := loadBasicsFromEnvironmentVars(params[0], params[1], params[2], params[3], params[4], params[5])
	getConfigFromSpringCloudConfigServer(newConfig["spring.cloud.config.uri"].(string), newConfig)
	Configuration = newConfig
	log.Println("Config loaded correctly")
}

func loadBasicsFromEnvironmentVars(spring_profiles_active, spring_cloud_config_uri, spring_cloud_config_label, server_port, eureka_instance_ip_address, spring_application_name string) map[string]interface{} {
	var newConfig map[string]interface{} = make(map[string]interface{})
	newConfig["spring.profiles.active"] = spring_profiles_active
	newConfig["spring.cloud.config.uri"] = spring_cloud_config_uri
	newConfig["spring.cloud.config.label"] = spring_cloud_config_label
	newConfig["server.port"] = server_port
	newConfig["eureka.instance.ip-address"] = eureka_instance_ip_address
	newConfig["spring.application.name"] = spring_application_name
	newConfig["hostname"], _ = os.Hostname()

	port, err := strconv.Atoi(newConfig["server.port"].(string))
	if err != nil {
		newConfig["server.port"] = _DEFAULT_PORT
	} else {
		newConfig["server.port"] = port
	}

	if newConfig["spring.profiles.active"] == "" || newConfig["spring.cloud.config.uri"] == "" || newConfig["spring.cloud.config.label"] == "" || newConfig["server.port"] == "" || newConfig["eureka.instance.ip-address"] == 0 || newConfig["spring.application.name"] == "" {
		panic("spring_profiles_active , spring_cloud_config_uri , spring_cloud_config_label , server_port , eureka_instance_ip_address, spring_application_name environment vars are mandatories")
	}

	return newConfig
}

func getConfigFromSpringCloudConfigServer(uriEndpoint string, newConfig map[string]interface{}) {
	finalEndpoint := uriEndpoint + "/" + newConfig["spring.application.name"].(string) + "/" + newConfig["spring.profiles.active"].(string) + "/"
	log.Println("Getting config from " + finalEndpoint)
	rs, err := getJsonFromSpringCloudConfigServer(finalEndpoint)
	if err != nil {
		panic("can't load configuration from " + finalEndpoint)
	}
	rewriteConfig(rs, newConfig)
}

func rewriteConfig(container *gabs.Container, newConfig map[string]interface{}) {
	newConfig["label"], _ = container.Path("label").Data().(string)
	newConfig["name"], _ = container.Path("name").Data().(string)
	source := container.Path("propertySources").Path("source")
	propertySources, _ := source.Children()

	iterateOverEachKeyAndReplaceVars(propertySources, newConfig)
	replaceVars(newConfig)

}

func replaceVars(newConfig map[string]interface{}) {
	for field, value := range newConfig {
		if isString(value) {
			if strings.Contains(value.(string), "${") {
				modifiedValue := value.(string)
				splitted := strings.Split(value.(string), "${")
				for i:=0;i<len(splitted);i++ {
					fieldToFind := strings.Split(splitted[i], "}")[0]
					if newConfig[fieldToFind] != nil {
						modifiedValue = strings.Replace(modifiedValue, "${"+fieldToFind+"}", newConfig[fieldToFind].(string), 10)
						newConfig[field] = modifiedValue
					}
				}
			}
		}
	}
}

func iterateOverEachKeyAndReplaceVars(containers []*gabs.Container, newConfig map[string]interface{}) {
	for _, child := range containers {
		keyvalueconfigurationmap, _ := child.ChildrenMap()
		for configurationField, configurationValue := range keyvalueconfigurationmap {
			modifiedConfigurationValue := configurationValue.Data()
			if isString(modifiedConfigurationValue) {
				if configurationValueThanMustBeReplaced(modifiedConfigurationValue) {
					modifiedConfigurationValue = replaceConfigurationValueAndReturnTheNewValue(modifiedConfigurationValue, newConfig)
				}
			}
			addNewKeyValueToConfigurationIfNotExists(configurationField, modifiedConfigurationValue, newConfig)
		}

	}
}

func replaceConfigurationValueAndReturnTheNewValue(modifiedConfigurationValue interface{}, newConfig map[string]interface{}) interface{} {
	splitted := strings.Split(modifiedConfigurationValue.(string), "${")
	for i:=1; i<len(splitted);i++  {
		fieldName := strings.Split(splitted[i], "}")[0]
		if newConfig[fieldName] != nil {
			fieldValue := newConfig[fieldName].(string)
			modifiedConfigurationValue = strings.Replace(modifiedConfigurationValue.(string), "${"+fieldName+"}", fieldValue, 1)
		}
	}
	return modifiedConfigurationValue
}

func configurationValueThanMustBeReplaced(data interface{}) bool {
	return strings.Contains(data.(string), "${")
}
func isString(data interface{}) bool {
	return reflect.TypeOf(data).String() == "string"
}

func addNewKeyValueToConfigurationIfNotExists(key string, value interface{}, newConfig map[string]interface{}) {
	if newConfig[key] == nil {
		newConfig[key] = value
	}
}

func getJsonFromSpringCloudConfigServer(url string) (*gabs.Container, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}

	jsonParsed, err := gabs.ParseJSON(b)

	if err != nil {
		return nil, err
	}

	return jsonParsed, nil
}

func AddKeyValueToConfig(key string, value interface{}) {
	Configuration[key] = value
}

func preloadConfigurationParams() []string {
	var params []string = make([]string, 6)

	params[0] = os.Getenv("spring_profiles_active")
	params[1] = os.Getenv("spring_cloud_config_uri")
	params[2] = os.Getenv("spring_cloud_config_label")
	params[3] = os.Getenv("server_port")
	params[4] = os.Getenv("eureka_instance_ip_address")
	params[5] = os.Getenv("spring_application_name")

	if len(os.Args)>5 {
		for i:=0;i<6;i++ {
			params[i] = os.Args[i+1]
		}
	}

	return params
}