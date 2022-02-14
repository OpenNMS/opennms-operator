package values

/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"github.com/OpenNMS/opennms-operator/config"
	"github.com/OpenNMS/opennms-operator/internal/model/values"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

// GetDefaultValues - get the default Helm/Template values
func GetDefaultValues(operatorConfig config.OperatorConfig) values.TemplateValues {
	defaultValues := LoadValues(operatorConfig.DefaultOpenNMSValuesFile)
	defaultValues = SetServiceImages(operatorConfig, defaultValues)
	return defaultValues
}

// LoadValues - load Helm/Template values from the given file
func LoadValues(filename string) values.TemplateValues {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error reading default values from file: %v", err)
	}
	var defValues values.Values
	err = yaml.Unmarshal(yamlFile, &defValues)
	if err != nil {
		log.Fatalf("error unmarshalling default values from file: %v", err)
	}
	templateValues := values.TemplateValues{
		Values: defValues,
	}
	return templateValues
}

// SetServiceImages - set service images in the Helm/Template values based on ENV config
func SetServiceImages(config config.OperatorConfig, v values.TemplateValues) values.TemplateValues {
	v.Values.Grafana.Image = config.ServiceImageGrafana
	v.Values.Auth.Image = config.ServiceImageAuth

	v.Values.OpenNMS.InitContainerImage = config.ServiceImageInit
	return v
}