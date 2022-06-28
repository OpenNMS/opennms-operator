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
	"fmt"
	"github.com/OpenNMS/opennms-operator/config"
	"github.com/OpenNMS/opennms-operator/internal/model/values"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

// GetDefaultValues - get the default Helm/Template values
func GetDefaultValues(operatorConfig config.OperatorConfig) values.TemplateValues {
	v := LoadValues(operatorConfig.DefaultOpenNMSValuesFile, operatorConfig.DefaultOperatorValuesFile)
	return values.TemplateValues{
		Values: v,
	}
}

// LoadValues - load Helm/Template values from the given files
func LoadValues(filename1, filename2 string) values.Values {
	yamlFile1, err := ioutil.ReadFile(filename1)
	if err != nil {
		log.Fatalf("error reading default values from file: %v", err)
	}
	yamlFile2, err := ioutil.ReadFile(filename2)
	if err != nil {
		log.Fatalf("error reading default values from file: %v", err)
	}
	combinedYaml := fmt.Sprintf("%s\n%s", yamlFile1, yamlFile2)
	var defValues values.Values
	err = yaml.Unmarshal([]byte(combinedYaml), &defValues)
	if err != nil {
		log.Fatalf("error unmarshalling default values from file: %v", err)
	}
	return defValues
}
