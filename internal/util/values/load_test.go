// +build unit

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

package values

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadValues(t *testing.T) {
	writeTestFiles(t)

	values := LoadValues(TestFilename)
	assert.NotNil(t, values)
	assert.Equal(t, "testNamespace", values.Values.Namespace, "should load the yaml file correctly")
}

var TestFilename = "./testtmp/test1.yaml"

var TestFileContents = `Namespace: testNamespace`

func writeTestFiles(t *testing.T) {
	file1 := []byte(TestFileContents)

	newpath := filepath.Join(".", "testtmp")
	err := os.MkdirAll(newpath, os.ModePerm)
	assert.Nil(t, err)

	err = os.WriteFile(TestFilename, file1, 0644)
	assert.Nil(t, err)
}